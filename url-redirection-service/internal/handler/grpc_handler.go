package handler

import (
	"context"
	"log"

	"github.com/hamza4253/tiny-url/redirect/internal/service"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
	"google.golang.org/grpc"
)

type GRPCRedirectionHandler struct {
	pb.UnimplementedRedirectionServiceServer
	service *service.RedirectionService
}

func NewGRPCRedirectionHandler(grpcServer *grpc.Server, s *service.RedirectionService) {
	handler := &GRPCRedirectionHandler{
		service: s,
	}
	pb.RegisterRedirectionServiceServer(grpcServer, handler)
}

func (h *GRPCRedirectionHandler) LookupURL(ctx context.Context, req *pb.LookupRequest) (*pb.LookupResponse, error) {
	shortUrl := req.ShortUrl
	longUrl, err := h.service.GetURLByShortURL(ctx, shortUrl)
	if err != nil {
		log.Printf("Error getting long url for %v. Error: %v", shortUrl, err)
		return nil, err
	}
	response := &pb.LookupResponse{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}
	return response, nil
}
