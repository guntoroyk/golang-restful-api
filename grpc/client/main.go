package main

import (
	"context"
	"github.com/guntoroyk/golang-restful-api/grpc/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := proto.NewCategoryServiceClient(cc)

	getAllCategory(c)
	getCategory(c)
}

func getAllCategory(c proto.CategoryServiceClient) {
	req := &proto.GetAllCategoryRequest{}
	res, err := c.GetAllCategory(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling Greet RPC: %v", err)
	}
	log.Printf("Categories: %v\n", res)
}

func getCategory(c proto.CategoryServiceClient) {
	req := &proto.GetCategoryRequest{
		CategoryId: 1,
	}
	res, err := c.GetCategory(context.Background(), req)
	if err != nil {
		log.Fatalf("Error when calling Greet RPC: %v", err)
	}
	log.Printf("Category: %v\n", res)
}
