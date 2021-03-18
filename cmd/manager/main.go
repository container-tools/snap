package main

import (
	"context"
	"fmt"
	"os"

	"github.com/container-tools/snap/pkg/api"
	"github.com/container-tools/snap/pkg/client"
	"github.com/container-tools/snap/pkg/util/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Expected 1 argument, got %d\n", len(os.Args)-1)
		os.Exit(1)
	}
	lib := os.Args[1]

	ctx := context.Background()

	logf.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	cfg, err := client.NewRestConfig()
	if err != nil {
		log.Info("can't get Kubernetes configuration")
		panic(err)
	}

	ns, err := client.GetCurrentNamespace()
	if err != nil {
		log.Info("can't get current Kubernetes namespace")
		panic(err)
	}

	options := api.SnapOptions{}

	snap, err := api.NewSnap(cfg, ns, false, options)
	if err != nil {
		log.Info("can't initialize snap")
		panic(err)
	}

	id, err := snap.Deploy(ctx, lib)
	if err != nil {
		log.Info("error during deployment")
		panic(err)
	}

	log.Infof("Deployed application %s", id)

	log.Info("Terminating...")
}
