package main

import (
	"context"

	"github.com/Wang-Kai/quotar/pb"
	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

	log.WithFields(log.Fields{
		"name":  createDirReq.Name,
		"quota": createDirReq.Quota,
	}).Info("Call RPC to create dir")

	resp, err := QuotarClient.CreateDir(context.Background(), createDirReq)
	if err != nil {
		return nil, errors.Wrap(err, "call quotar to create dir")
	}

	log.WithField("path", resp.Path).Info("Create Dir response")

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: opt.PVName,
		},
		Spec: v1.PersistentVolumeSpec{
			AccessModes: opt.PVC.Spec.AccessModes,
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
	log.WithField("name", req.Name).Info("Call RPC to delete dir")

	_, err := QuotarClient.DeleteDir(context.Background(), req)
	if err != nil {
		return errors.Wrap(err, "Call quotar to delete dir")
	}

	return nil
}
