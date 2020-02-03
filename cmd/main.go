package main

import (
	"context"

	"github.com/nicolaferraro/snap/pkg/api"
	"github.com/nicolaferraro/snap/pkg/client"
	"github.com/nicolaferraro/snap/pkg/util/log"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {

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

	id, err := snap.Deploy(ctx, "./example")
	if err != nil {
		log.Info("error during deployment")
		panic(err)
	}

	log.Infof("Deployed application %s", id)

	log.Info("Terminating...")
}
