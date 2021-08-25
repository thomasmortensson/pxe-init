package client

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/thomasmortensson/pxe-init/internal/drivers/cli/common"
)

const (
	ArgGrpcEndpoint string = "grpc-endpoint"

	DefaultGrpcEndpoint string = "grpc://localhost:5000"

	grpcTimeout = 10 * time.Second
)

var (
	// To be filled by compiler
	version string

	// RootCmd is the base pxe-init-client command.
	RootCmd = &cobra.Command{
		Use:     "pxe-init-client",
		Short:   "A command line client for the pxe-init service.",
		Version: version,
	}

	globalFlags = struct {
		endpoint string
		// Extend to use mTLS auth and stop setting client option: grpc.WithInsecure() (Use cert-manager :) )
		// I just didn't have time to set this up but I'd want to setup TLS on comms between server/client and the postgres database
		caFile   string
		certFile string
		keyFile  string
	}{}
)

func addRootFlags() {
	RootCmd.PersistentFlags().StringVar(
		&globalFlags.endpoint,
		ArgGrpcEndpoint,
		DefaultGrpcEndpoint,
		"gRPC endpoint to use for communication with pxe-init-server",
	)
}

// Execute is the main entrypoint to the pxe-init-client program
func Execute() {
	addRootFlags()
	addListImagesCommand()
	addRegisterImageMachineCommand()

	if err := RootCmd.Execute(); err != nil {
		common.ErrorPrint("Failed to run pxe-init-client", err)
	}
}

// endpointFromCmd returns the grpc-endpoint argument
func endpointFromCmd(cmd *cobra.Command) string {
	endpoint, err := cmd.Flags().GetString(ArgGrpcEndpoint)
	if err != nil {
		common.ErrorExit(fmt.Sprintf("Unable to get: %v", ArgGrpcEndpoint), err)
	}

	endpointURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		common.ErrorExit(fmt.Sprintf("Unable to parse provided endpoint: %v", endpoint), err)
	}
	return endpointURL.Host
}
