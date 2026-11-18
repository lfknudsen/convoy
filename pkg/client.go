package pkg

import (
	"context"
	"fmt"
	"log"
	"uuid"

	. "convoy/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConvoyClient struct {
	conn    *grpc.ClientConn
	speaker QueueClient
	id      string
}

func NewConvoyClient(port int) (client *ConvoyClient) {
	client = new(ConvoyClient{})

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	connString := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.NewClient(connString, opts...)
	if err != nil {
		log.Fatal(err)
	}
	client.conn = conn
	client.speaker = NewQueueClient(conn)
	client.id = uuid.New().String()
	return client
}

func (client *ConvoyClient) Push(url string) {
	id := uuid.New().String()
	link := Link{
		Id:  id,
		Url: url,
	}
	var callOpts []grpc.CallOption
	reply, err := client.speaker.Push(context.Background(), &link, callOpts...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Is the download finished? ", reply.Added)
}
