package kubernetes

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nicolaferraro/snap/pkg/client"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func TestPF(t *testing.T) {
	t.SkipNow()

	_, err := client.NewClient()
	assert.Nil(t, err)
	conf, err := config.GetConfig()
	assert.Nil(t, err)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var host string
	host, err = PortForward(ctx, conf, "snap2", "minio-7f4cf589bb-6s258")
	assert.Nil(t, err)

	fmt.Printf("Forward address: %s\n", host)

	if !t.Failed() {
		fmt.Printf("Waiting few seconds...\n")
		time.Sleep(3 * time.Second)
	}
}
