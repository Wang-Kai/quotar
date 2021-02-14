package main

import (
	"context"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type nfsProvisioner struct {
}

// Provision
func (n *nfsProvisioner) Provision(opt controller.ProvisionOptions) (*v1.PersistentVolume, error) {
	// get pvc request storage
	resourceQuantity := opt.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	quota := resourceQuantity.String()

	// create project dir
	var createDirReq = &pb.CreateDirReq{
		Name:  opt.PVName,
		Quota: quota,
	}

	resp, err := QuotarClient.CreateDir(context.Background(), createDirReq)
	if err != nil {
		return nil, errors.Wrap(err, "call quotar to create dir")
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: opt.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resourceQuantity,
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Server: NFS_SERVER,
					Path:   resp.Path,
				},
			},
		},
	}

	return pv, nil
}

func (n *nfsProvisioner) Delete(pv *v1.PersistentVolume) error {
	req := &pb.DeleteDirReq{
		Name: pv.Name,
	}

	_, err := QuotarClient.DeleteDir(context.Background(), req)
	if err != nil {
		return errors.Wrap(err, "Call quotar to delete dir")
	}

	return nil
}
