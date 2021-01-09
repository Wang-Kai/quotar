package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Wang-Kai/quotar/pb"
	_ "github.com/Wang-Kai/quotar/pkg/conf"
	"github.com/Wang-Kai/quotar/pkg/svc"
	_ "github.com/Wang-Kai/quotar/pkg/xfs"
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

	// for grpc debug
	reflection.Register(s)
	pb.RegisterQuotarServer(s, &svc.QuotarService{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
