package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	rpc "github.com/thomasmortensson/pxe-init/internal/adapters/grpc"
	pb "github.com/thomasmortensson/pxe-init/internal/adapters/grpc/v1"
	adapterHttp "github.com/thomasmortensson/pxe-init/internal/adapters/http"
	"github.com/thomasmortensson/pxe-init/internal/adapters/postgres"
	"github.com/thomasmortensson/pxe-init/internal/drivers/cli/common"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the pxe-init-server",
		RunE:  serve,
	}
)

func serve(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		common.ErrorPrint("Failed to setup logger", err)
		return err
	}
	defer logger.Sync() // nolint:errcheck // In defer function

	// Get args
	grpcPort := globalFlags.grpcPort
	httpPort := globalFlags.httpPort

	if err = verifySSLMode(globalFlags.databaseSSLMode); err != nil {
		logger.Error(
			"Unknown SSL mode",
			zap.String("sslmode", globalFlags.databaseSSLMode),
			zap.Error(err),
		)
		return err
	}

	connProvider := postgres.NewConnProvider(
		globalFlags.databaseHost,
		globalFlags.databasePort,
		globalFlags.databaseUser,
		globalFlags.databasePassword,
		globalFlags.databaseName,
		globalFlags.databaseSSLMode,
	)

	db, err := postgres.NewPgDatastore(logger, connProvider)
	if err != nil {
		logger.Error(
			"Unable to setup DB connection",
			zap.Error(err),
		)
		return err
	}

	forwardServer, err := url.ParseRequestURI(globalFlags.forwardServer)
	if err != nil {
		logger.Error(
			"Unable to parse forward server address",
			zap.Error(err),
		)
		return err
	}

	// Setup HTTP server
	serverCtx := adapterHttp.NewServerContext(logger, db, forwardServer)
	router := adapterHttp.NewRouter(serverCtx)
	router.AddRoutes()

	httpServer := &http.Server{
		// Can tighten to not bind to all interfaces. For now this is fine
		Addr:    fmt.Sprintf(":%v", httpPort),
		Handler: router.Engine,
	}

	// Setup grpc server
	grpcEndpoint := fmt.Sprintf(":%d", grpcPort)
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logger.Error(
			"Failed to listen on port",
			zap.String("endpoint", grpcEndpoint),
			zap.Error(err),
		)
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pxeServer := rpc.NewPxeInitServer(logger, db, forwardServer)
	pb.RegisterPxeInitServer(
		grpcServer,
		pxeServer,
	)

	// Run servers
	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, os.Interrupt)

	go func() {
		// service http connections
		logger.Info(
			"HTTP server startup",
			zap.Int("http-port", httpPort),
		)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(
				"Failed to listen on address",
				zap.String("address", httpServer.Addr),
				zap.Error(err),
			)
		}
	}()

	go func() {
		// service grpc connections
		logger.Info(
			"gRPC server startup",
			zap.Int("grpc-port", grpcPort),
		)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error(
				"Failed to listen on endpoint",
				zap.String("endpoint", grpcEndpoint),
			)
		}
	}()

	<-endChan

	// Shutdown
	logger.Info("exiting pxe-init-server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failure",
			zap.Error(err),
		)
		return err
	}

	grpcServer.GracefulStop()

	return nil
}

func addServeCommand() {
	RootCmd.AddCommand(serveCmd)
}
