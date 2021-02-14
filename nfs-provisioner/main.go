package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Wang-Kai/quotar/pb"

	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var (
	NFS_SERVER       string
	QUOTAR_SERVER    string
	PROVISIONER_NAME string
	QuotarClient     pb.QuotarClient
)

func init() {
	/*
		collect envs, including:
		1. NFS_SERVER
		2. QUOTAR_SERVER
		3. STORAGECLASS_NAME
	*/
	var requiredEnv = []string{"NFS_SERVER", "QUOTAR_SERVER", "PROVISIONER_NAME"}
	for _, envName := range requiredEnv {
		if os.Getenv(envName) == "" {
			msg := fmt.Sprintf("Require %s, but not found", envName)
			panic(msg)
		}
	}

	NFS_SERVER = os.Getenv("NFS_SERVER")
	QUOTAR_SERVER = os.Getenv("QUOTAR_SERVER")
	PROVISIONER_NAME = os.Getenv("PROVISIONER_NAME")

	// init gRPC client for call quotar server
	cc, err := grpc.Dial(QUOTAR_SERVER, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	QuotarClient = pb.NewQuotarClient(cc)
}

func main() {
	clientset := genClientset()
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		klog.Fatalf("Error getting server version: %v", err)
	}

	provisioner := &nfsProvisioner{}
	pc := controller.NewProvisionController(
		clientset,
		PROVISIONER_NAME,
		provisioner,
		serverVersion.GitVersion,
	)

	done := make(chan struct{})
	pc.Run(done)
}

func genClientset() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
