package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"telegram/internal/service"
	"telegram/pb"

	"google.golang.org/grpc"
)

func main() {
	token := os.Getenv("TOKEN")

	srv := service.New(token)

	gRPCServer := grpc.NewServer()
	pb.RegisterTelegramServer(gRPCServer, srv)

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("telegram service is running on :5001")
		if err = gRPCServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	<-quit
	log.Println("stopping telegram service...")
	gRPCServer.GracefulStop()
	log.Println("stopped telegram service")

}
