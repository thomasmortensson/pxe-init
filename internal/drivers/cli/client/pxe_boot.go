package client

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/thomasmortensson/pxe-init/internal/adapters/grpc/v1"
	"github.com/thomasmortensson/pxe-init/internal/drivers/cli/common"
)

const (
	ArgMacAddress = "mac-address"
	ArgImage      = "image-name"
)

var (
	registerImageMachineCmd = &cobra.Command{
		Use:   "register-image-machine",
		Short: "Register an image with machine with MAC",
		RunE:  registerImageMachine,
	}

	registerImageMachineFlags = struct {
		Mac   string
		Image string
	}{}
)

func registerImageMachine(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		common.ErrorPrint("Failed to setup logger", err)
		return err
	}
	defer logger.Sync() // nolint:errcheck // In defer function

	// Get args from cmd
	endpoint := endpointFromCmd(cmd)

	logger.Debug("In pxe-boot command",
		zap.String("endpoint", endpoint),
	)

	// socket connect to server
	conn, err := grpc.Dial(endpoint, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		logger.Error(
			"Failed to dial endpoint",
			zap.Error(err),
		)
		return err
	}
	defer conn.Close()

	// Setup 10sec timeout
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()

	// Perform RPC connection to server
	client := pb.NewPxeInitClient(conn)
	_, err = client.RegisterImageMachine(ctx, &pb.RegisterImageMachineRequest{
		Machine: &pb.Machine{
			Mac:   registerImageMachineFlags.Mac,
			Image: registerImageMachineFlags.Image,
		},
	})
	if err != nil {
		logger.Error(
			"Failed to call PxeBoot method",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func addRegisterImageMachineCommand() {
	RootCmd.AddCommand(registerImageMachineCmd)
	addRegisterImageMachineFlags()
}

func addRegisterImageMachineFlags() {
	registerImageMachineCmd.Flags().StringVar(
		&registerImageMachineFlags.Mac,
		ArgMacAddress,
		"",
		"MAC address of machine to register image to",
	)

	registerImageMachineCmd.Flags().StringVar(
		&registerImageMachineFlags.Image,
		ArgImage,
		"",
		"Image name to register to machine",
	)
}
