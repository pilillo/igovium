package main

import (
	"context"
	"log"

	"github.com/pilillo/igovium/service/cachepb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := cachepb.NewCacheServiceClient(cc)

	// put key
	putReq := &cachepb.PutRequest{
		Key:   "key",
		Value: []byte("value"),
		Ttl:   "1m",
	}
	putRes, err := c.Put(context.Background(), putReq)
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}
	log.Printf("put response: %v", putRes)

	getReq := &cachepb.GetRequest{Key: "key"}
	var getRes *cachepb.GetResponse
	getRes, err = c.Get(context.Background(), getReq)
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("get response: %v", getRes)
}
