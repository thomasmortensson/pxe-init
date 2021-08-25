package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thomasmortensson/pxe-init/internal/drivers/cli/common"
)

const (
	ArgGRPCPort         string = "grpc-port"
	ArgHTTPPort         string = "http-port"
	ArgForwardServer    string = "pxe-forward-server"
	ArgDatabaseHost     string = "db-host"
	ArgDatabasePort     string = "db-port"
	ArgDatabaseUser     string = "db-user"
	ArgDatabasePassword string = "db-password"
	ArgDatabaseName     string = "db-name"
	ArgDatabaseSSLMode  string = "db-ssl-mode"

	DefaultGRPCPort         int    = 5000
	DefaultHTTPPort         int    = 8080 // Unauthenticated as in pxe therefore no SSL :(
	DefaultForwardServer    string = "http://zbox.mortcloud.com"
	DefaultDatabaseHost     string = "127.0.0.1"
	DefaultDatabasePort     int    = 5432
	DefaultDatabaseUser     string = "pxe_init"
	DefaultDatabasePassword string = ""
	DefaultDatabaseName     string = "pxe_init"
	DefaultDatabaseSSLMode  string = "disable"

	EnvGRPCPort         string = "PXE_SRV_GRPC_PORT"
	EnvHTTPPort         string = "PXE_SRV_HTTP_PORT"
	EnvForwardServer    string = "PXE_SRV_FORWARD_SERVER"
	EnvDatabaseHost     string = "PXE_SRV_DB_HOST"
	EnvDatabasePort     string = "PXE_SRV_DB_PORT"
	EnvDatabaseUser     string = "PXE_SRV_DB_USER"
	EnvDatabasePassword string = "PXE_SRV_DB_PASSWORD" // nolint:gosec // This is not the password you're looking for
	EnvDatabaseName     string = "PXE_SRV_DB_NAME"
	EnvDatabaseSSLMode  string = "PXE_SRV_DB_SSL_MODE"
)

const (
	DatabaseSSLModeDisable    string = "disable"
	DatabaseSSLModeAllow      string = "allow"
	DatabaseSSLModePrefer     string = "prefer"
	DatabaseSSLModeRequire    string = "require"
	DatabaseSSLModeVerifyCA   string = "verify-ca"
	DatabaseSSLModeVerifyFull string = "verify-full"
)

var DatabaseSSLModes = [...]string{
	DatabaseSSLModeDisable,
	DatabaseSSLModeAllow,
	DatabaseSSLModePrefer,
	DatabaseSSLModeRequire,
	DatabaseSSLModeVerifyCA,
	DatabaseSSLModeVerifyFull,
}

var (
	// To be filled by compiler
	version string

	// RootCmd is the base pxe-init-server command.
	RootCmd = &cobra.Command{
		Use:     "pxe-init-server",
		Short:   "The pxe-init server process",
		Version: version,
	}

	globalFlags = struct {
		grpcPort      int
		httpPort      int
		forwardServer string

		databaseHost     string
		databasePort     int
		databaseUser     string
		databasePassword string
		databaseName     string
		databaseSSLMode  string
	}{}
)

func addRootFlags() {
	RootCmd.PersistentFlags().IntVar(
		&globalFlags.grpcPort,
		ArgGRPCPort,
		common.GetEnvOrDefaultInt(EnvGRPCPort, DefaultGRPCPort),
		"gRPC port to serve pxe-init-server on",
	)

	RootCmd.PersistentFlags().IntVar(
		&globalFlags.httpPort,
		ArgHTTPPort,
		common.GetEnvOrDefaultInt(EnvHTTPPort, DefaultHTTPPort),
		"HTTP port to serve pxe-init-server on",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.databaseHost,
		ArgDatabaseHost,
		common.GetEnvOrDefaultString(EnvDatabaseHost, DefaultDatabaseHost),
		"Database host endpoint to use for communication",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.forwardServer,
		ArgForwardServer,
		common.GetEnvOrDefaultString(EnvForwardServer, DefaultForwardServer),
		"Forward server for PXE assets",
	)

	RootCmd.PersistentFlags().IntVar(
		&globalFlags.databasePort,
		ArgDatabasePort,
		common.GetEnvOrDefaultInt(EnvDatabasePort, DefaultDatabasePort),
		"Database port to use for communication",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.databaseUser,
		ArgDatabaseUser,
		common.GetEnvOrDefaultString(EnvDatabaseUser, DefaultDatabaseUser),
		"Database user",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.databasePassword,
		ArgDatabasePassword,
		common.GetEnvOrDefaultString(EnvDatabasePassword, DefaultDatabasePassword),
		"Database password",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.databaseName,
		ArgDatabaseName,
		common.GetEnvOrDefaultString(EnvDatabaseName, DefaultDatabaseName),
		"Database name",
	)

	RootCmd.PersistentFlags().StringVar(
		&globalFlags.databaseSSLMode,
		ArgDatabaseSSLMode,
		common.GetEnvOrDefaultString(EnvDatabaseSSLMode, DefaultDatabaseSSLMode),
		fmt.Sprintf("Database SSL mode choice %v", DatabaseSSLModes),
	)
}

// Execute is the main entrypoint to the pxe-init-server program
func Execute() {
	addRootFlags()
	addServeCommand()

	if err := RootCmd.Execute(); err != nil {
		common.ErrorExit("Failed to run pxe-init-server", err)
	}
}

func verifySSLMode(sslmode string) error {
	for _, mode := range DatabaseSSLModes {
		if mode == sslmode {
			return nil
		}
	}
	return fmt.Errorf("unknown sslmode: %v", sslmode)
}
