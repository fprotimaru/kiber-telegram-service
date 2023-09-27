package main
// test
import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"telegram/internal/repository"
	"telegram/internal/service"
	"telegram/pb"

	"google.golang.org/grpc"
)

const (
// TOKEN   = "5722224930:AAHzVtMAl-OwftEVbw6Ululaef16o_yaFLA"
// PsqlUrl = "postgres://root:1@localhost:5432/telegram_service?sslmode=disable"
)

func main() {
	token := os.Getenv("TOKEN")
	psqlUrl := os.Getenv("PSQL_URL")
	println(123)
	dbConn, err := repository.NewTelegramUserRepository(psqlUrl)
	if err != nil {
		log.Fatalf("repository.NewTelegramUserRepository error: %v\n", err)
	}

	srv := service.New(token, dbConn)
	go srv.Listen(context.Background())

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
