package svc

import (
	"context"
	"fmt"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/Wang-Kai/quotar/pkg/conf"
	"github.com/Wang-Kai/quotar/pkg/xfs"
)

type QuotarService struct {
}

func (q *QuotarService) CreateDir(ctx context.Context, req *pb.CreateDirReq) (*pb.CreateDirResp, error) {

	fmt.Printf("Create %s dir, and with size %s\n", req.Name, req.Quota)

	// create xfs project
	if err := xfs.CreatePrj(req.Name, req.Quota); err != nil {
		return nil, err
	}

	dirPath := fmt.Sprintf("%s/%s", conf.WORKSPACE, req.Name)

	return &pb.CreateDirResp{
		Path: dirPath,
	}, nil
}
