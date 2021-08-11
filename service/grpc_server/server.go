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
	cacheEntry := &cache.CacheEntry{Key: req.Key, Value: req.Value, TTL: req.Ttl}
	err := cacheService.Put(cacheEntry)
	return &cachepb.Empty{}, err
}

func (*grpcServer) Delete(ctx context.Context, req *cachepb.DeleteRequest) (*cachepb.Empty, error) {
	log.Printf("RPC Delete request for key %s", req.Key)
	err := cacheService.Delete(req.Key)
	return &cachepb.Empty{}, err
}

func (*grpcServer) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	log.Printf("RPC Get request for key %s", req.Key)
	val, err := cacheService.Get(req.Key)
	if err != nil {
		return nil, err
	}
	byteVal, err := utils.GetBytes(val)
	if err != nil {
		return nil, err
	}
	response := &cachepb.GetResponse{Value: byteVal}
	return response, nil
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
