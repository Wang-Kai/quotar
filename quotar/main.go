package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/Wang-Kai/quotar/pkg/conf"
	"github.com/Wang-Kai/quotar/pkg/svc"
	_ "github.com/Wang-Kai/quotar/pkg/xfs"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := fmt.Sprintf(":%d", conf.PORT)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	// for grpc debug
	reflection.Register(s)
	pb.RegisterQuotarServer(s, &svc.QuotarService{})

	log.WithField("Port", port).Info("gRPC server running ...")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
