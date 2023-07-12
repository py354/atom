package main

import (
	"atom/internal/app"
	"atom/internal/config"
	"atom/internal/db"
	"atom/internal/repository"
	pb "atom/pkg/api"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("I'm server")

	ctx := context.Background()

	store, err := db.New(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	service := app.New(repository.New(store))
	addr := fmt.Sprintf("%s:%s", conf.Grpc.Host, conf.Grpc.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	server := grpc.NewServer()
	pb.RegisterFleetServer(server, service)
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
