package svc

import (
	"context"
	"fmt"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/Wang-Kai/quotar/pkg/xfs"
)

type QuotarService struct {
}

var workspace = `/data/`

func (q *QuotarService) CreateDir(ctx context.Context, req *pb.CreateDirReq) (*pb.CreateDirResp, error) {

	fmt.Printf("Create %s dir, and with size %d\n", req.Name, req.Size)

	path := fmt.Sprintf("%s/%s", workspace, req.Name)

	// create xfs project
	if err := xfs.CreatePrj(req.Name, path, workspace, req.Size); err != nil {
		return nil, err
	}

	return &pb.CreateDirResp{RetCode: 0}, nil
}
