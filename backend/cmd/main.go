package main
import (
	spot_exchange "git.gendocu.com/gendocu/SpotExchange.git/sdk/go"
	"github.com/gendocu-com-examples/spot-exchange/backend/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	log.Println("starting container")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8010"
		log.Printf("Defaulting to port %s", port)
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Got error: %+v", err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer) //not required, but useful for debugging
	spot_exchange.RegisterSpotExchangeServer(grpcServer, internal.NewService())
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("got an error", err)
	}
}
