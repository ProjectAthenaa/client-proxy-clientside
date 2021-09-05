package main

import (
	"context"
	"flag"
	certificate_module "github.com/ProjectAthenaa/sonic-core/certificate"
	"github.com/ProjectAthenaa/sonic-core/protos/clientProxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

var (
	userID string
)

func main() {
	flag.StringVar(&userID, "user_id", "", "User ID")
	flag.Parse()
	if userID == "" {
		log.Fatal("UserID cannot be empty")
	}

	certs, _ := certificate_module.LoadClientTestCertificate()

	conn, err := grpc.Dial("secure.athenabot.com:443", grpc.WithTransportCredentials(certs))
	if err != nil {
		log.Fatal(err)
	}

	client := client_proxy.NewProxyClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "UserID", userID)

	//ctx, _ = context.WithTimeout(ctx, time.Second * 5)

	stream, err := client.Register(ctx)
	if err != nil{
		log.Fatal(err)
	}


	cmd, err := stream.Recv()
	if err != nil {
		log.Fatal(err)
	}

	if v, ok := cmd.Headers["STOP"]; !ok || v == "1" {
		log.Fatal("Server instructed stop")
	}

	proxy := NewServer(stream)

	go proxy.stopper()

	proxy.Start()
}
