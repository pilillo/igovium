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

	// put req
	putReq := &cachepb.PutRequest{
		Key:   "key",
		Value: `{"mykey":"this-is-my-test-value"}`,
		Ttl1:  "1m",
		Ttl2:  "10m",
	}
	log.Printf("putting: k='%s', v='%s'", putReq.Key, putReq.Value)
	putRes, err := c.Put(context.Background(), putReq)
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}
	log.Printf("put response: res='%v', err='%v'", putRes, err)

	getReq := &cachepb.GetRequest{Key: "key"}
	var getRes *cachepb.GetResponse
	getRes, err = c.Get(context.Background(), getReq)
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("get response: %v", getRes)
}
