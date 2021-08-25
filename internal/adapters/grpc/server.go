package grpc

import (
	"context"
	"net/url"

	"go.uber.org/zap"

	pb "github.com/thomasmortensson/pxe-init/internal/adapters/grpc/v1"
	"github.com/thomasmortensson/pxe-init/internal/domain/repositories"
	"github.com/thomasmortensson/pxe-init/internal/usecases"
)

type PxeInitServer struct {
	pb.UnimplementedPxeInitServer
	logger        *zap.Logger
	db            repositories.Datastore
	forwardServer *url.URL
}

func NewPxeInitServer(logger *zap.Logger, db repositories.Datastore, forwardServer *url.URL) *PxeInitServer {
	return &PxeInitServer{
		logger:        logger,
		db:            db,
		forwardServer: forwardServer,
	}
}

func (s *PxeInitServer) ListImages(
	ctx context.Context,
	request *pb.ListImagesRequest,
) (*pb.ListImagesResponse, error) {
	s.logger.Debug("Entering ListImages")
	defer s.logger.Debug("Exiting ListImages")

	usecase := usecases.NewListImages(s.logger, s.db)

	images, err := usecase.Execute()
	if err != nil {
		s.logger.Error(
			"Failed in ListImages usecase",
			zap.Error(err),
		)
		return &pb.ListImagesResponse{}, err
	}

	response := &pb.ListImagesResponse{}
	response.Images = make([]*pb.Image, len(images))
	for i, image := range images {
		response.Images[i] = &pb.Image{
			Name:   image.Name,
			Kernel: image.GetKernel(s.forwardServer.Host),
			Initrd: image.GetInitrd(s.forwardServer.Host),
			Rootfs: image.GetRootfs(s.forwardServer.Host),
		}
	}

	return response, nil
}

func (s *PxeInitServer) RegisterImageMachine(
	ctx context.Context,
	request *pb.RegisterImageMachineRequest,
) (*pb.RegisterImageMachineResponse, error) {
	s.logger.Debug("Entering RegisterImageMachine")
	defer s.logger.Debug("Exiting RegisterImageMachine")

	mac := request.GetMachine().GetMac()
	image := request.GetMachine().GetImage()

	usecase := usecases.NewRegisterMachineImage(s.logger, s.db)
	err := usecase.Execute(mac, image)
	if err != nil {
		s.logger.Error(
			"Failed in RegisterMachineImage usecase",
			zap.Error(err),
		)
	}

	// TODO UNIMPLEMENTED

	return &pb.RegisterImageMachineResponse{}, nil
}
