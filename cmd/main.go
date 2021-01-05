package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/Wang-Kai/quotar/pkg/svc"
)

const (
	PORT = ":10013"
)

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterQuotarServer(s, &svc.QuotarService{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
