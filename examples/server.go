package main

import (
	"flag"
	"fmt"
	"github.com/ramzeng/user-microservice/internal/core"
	"github.com/ramzeng/user-microservice/internal/model"
	"github.com/ramzeng/user-microservice/internal/service"
	"github.com/ramzeng/user-microservice/pb"
	"github.com/ramzeng/user-microservice/pkg/config"
	"github.com/ramzeng/user-microservice/pkg/database"
	"google.golang.org/grpc"
	"net"
)

func main() {
	configPath := flag.String("path", "", "the path of config file")
	migrate := flag.Bool("migrate", false, "auto migrate database")

	flag.Parse()

	core.ReadConfig(*configPath)

	core.Initialize()

	if *migrate {
		if err := database.AutoMigrate(&model.User{}); err != nil {
			panic(fmt.Sprintf("auto migrate database failed: %s", err))
		}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetInt("service.port")))

	if err != nil {
		panic(fmt.Sprintf("tcp listen failed: %s", err))
	}

	server := grpc.NewServer()

	pb.RegisterAuthServer(server, &service.Auth{})
	pb.RegisterUserServer(server, &service.User{})

	if err := server.Serve(listener); err != nil {
		panic(fmt.Sprintf("grpc server start failed: %s", err))
	}
}
