package client

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/thomasmortensson/pxe-init/internal/adapters/grpc/v1"
	"github.com/thomasmortensson/pxe-init/internal/drivers/cli/common"
)

var (
	listImagesCmd = &cobra.Command{
		Use:   "list-images",
		Short: "List available PXE images for boot",
		RunE:  listImages,
	}
)

func listImages(cmd *cobra.Command, args []string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		common.ErrorPrint("Failed to setup logger", err)
		return err
	}
	defer logger.Sync() // nolint:errcheck // In defer function

	// Get args from cmd
	endpoint := endpointFromCmd(cmd)

	logger.Debug("In list-images command",
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
	response, err := client.ListImages(ctx, &pb.ListImagesRequest{})
	if err != nil {
		logger.Error(
			"Failed to call ListImages method",
			zap.Error(err),
		)
		return err
	}

	// Format output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Kernel", "Initrd", "Rootfs"})

	for _, image := range response.Images {
		table.Append([]string{image.Name, image.Kernel, image.Initrd, image.Rootfs})
	}
	table.Render() // Send output

	return nil
}

func addListImagesCommand() {
	RootCmd.AddCommand(listImagesCmd)
}
