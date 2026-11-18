package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	. "convoy/proto"

	"google.golang.org/grpc"
)

type Server struct {
	UnimplementedQueueServer
	port int
	db   *sql.DB
}

func (s *Server) Serve(server *grpc.Server, lis net.Listener, wait chan struct{}) {
	log.Println("Listening on " + lis.Addr().String())
	err := server.Serve(lis)
	if err != nil {
		wait <- struct{}{}
		log.Fatal(err)
	}
}

func (s *Server) Push(_ context.Context, link *Link) (*Confirmation, error) {
	err := push(s.db, link.GetUrl())
	if err != nil {
		log.Println("Download error:", err)
	}
	return &Confirmation{
		Added: err == nil,
	}, err
}

func (s *Server) Pop(_ context.Context, worker *WorkerID) (string, error) {
	log.Printf("Worker %s requested a task.", worker)
	link, err := pop(s.db)
	if err != nil {
		log.Println("Download error:", err)
	}
	return link.Url, err
}

func StartServer(port int, db *sql.DB) chan struct{} {
	Srv := Server{
		port: port,
		db:   db,
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", Srv.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterQueueServer(grpcServer, &Srv)

	wait := make(chan struct{})
	go Srv.Serve(grpcServer, lis, wait)
	return wait
}
