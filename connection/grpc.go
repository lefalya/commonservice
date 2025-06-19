package connection

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func EstablishGRPC(grpcServerAddr string) (*grpc.ClientConn, context.Context, context.CancelFunc) {

	grpcConnection, errorDialGRPC := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if errorDialGRPC != nil {

		log.Fatalf("did not connect: %v", errorDialGRPC)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return grpcConnection, ctx, cancel
}
