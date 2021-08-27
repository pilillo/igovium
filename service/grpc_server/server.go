package grpc_server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pilillo/igovium/cache"
	"github.com/pilillo/igovium/service/cachepb"
	"github.com/pilillo/igovium/utils"
	"google.golang.org/grpc"
)

type grpcServer struct {
	cachepb.UnimplementedCacheServiceServer
}

func (*grpcServer) Put(ctx context.Context, req *cachepb.PutRequest) (*cachepb.Empty, error) {
	log.Printf("RPC Put request for key %s", req.Key)
	payload := cache.CachePayload(req.Value)
	cacheEntry := &cache.CacheEntry{Key: req.Key, Value: payload, TTL1: &req.Ttl1, TTL2: &req.Ttl2}
	response := cacheService.Put(cacheEntry)
	if response.Error != nil {
		return nil, fmt.Errorf("%s", response.Message)
	}
	return &cachepb.Empty{}, nil
}

func (*grpcServer) Delete(ctx context.Context, req *cachepb.DeleteRequest) (*cachepb.Empty, error) {
	log.Printf("RPC Delete request for key %s", req.Key)
	response := cacheService.Delete(req.Key)
	if response.Error != nil {
		return nil, fmt.Errorf("%s", response.Message)
	}
	return &cachepb.Empty{}, nil
}

func (*grpcServer) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	log.Printf("RPC Get request for key %s", req.Key)
	val, response := cacheService.Get(req.Key)
	if response != nil && response.Error != nil {
		return nil, fmt.Errorf("%s", response.Message)
	}
	cachePbResponse := &cachepb.GetResponse{Value: string(val)}
	return cachePbResponse, nil
}

var cacheService cache.CacheService

// on module import, get singleton cache service
func init() {
	cacheService = cache.GetCacheService()
}

func StartEndpoint(config *utils.Config) {

	err := cacheService.Init(config)
	if err != nil {
		panic(err)
	}
	grpcServerAddress := fmt.Sprintf("0.0.0.0:%d", config.GRPCConfig.Port)
	listener, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	cachepb.RegisterCacheServiceServer(s, &grpcServer{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
