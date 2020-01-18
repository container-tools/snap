package api

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/nicolaferraro/snap/pkg/deployer"
	"github.com/nicolaferraro/snap/pkg/deployer/java"
	"github.com/nicolaferraro/snap/pkg/installer"
	"github.com/nicolaferraro/snap/pkg/publisher"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DefaultBucketName = "snap"

	DefaultTimeout = 10 * time.Minute
)

type Snap struct {
	deployerModule  deployer.Deployer
	installerModule *installer.Installer
	publisherModule *publisher.Publisher

	namespace string
	direct    bool

	options SnapOptions
}

type SnapOptions struct {
	Bucket string

	Timeout time.Duration

	StdOut io.Writer
	StdErr io.Writer
}

func NewSnap(config *rest.Config, namespace string, direct bool, options SnapOptions) (*Snap, error) {
	if options.Bucket == "" {
		options.Bucket = DefaultBucketName
	}

	if options.Timeout <= 0 {
		options.Timeout = DefaultTimeout
	}

	if options.StdOut == nil {
		options.StdOut = ioutil.Discard
	}
	if options.StdErr == nil {
		options.StdErr = ioutil.Discard
	}

	client, err := ctrl.New(config, ctrl.Options{})
	if err != nil {
		return nil, err
	}
	return &Snap{
		deployerModule:  java.NewJavaDeployer(options.StdOut, options.StdErr),
		installerModule: installer.NewInstaller(config, client, options.StdOut, options.StdErr),
		publisherModule: publisher.NewPublisher(),

		namespace: namespace,
		direct:    direct,

		options: options,
	}, nil
}

func (s *Snap) Deploy(ctx context.Context, libraryDir string) error {
	deployCtx, cancel := context.WithTimeout(ctx, s.options.Timeout)
	defer cancel()

	// ensure installation
	if err := s.installerModule.EnsureInstalled(deployCtx, s.namespace); err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "snap-")
	if err != nil {
		return errors.Wrap(err, "cannot create a temporary dir")
	}
	defer os.RemoveAll(dir)

	if err := s.deployerModule.Deploy(libraryDir, dir); err != nil {
		return errors.Wrap(err, "error while creating deployment for source code")
	}

	host, err := s.installerModule.OpenConnection(deployCtx, s.namespace, s.direct)
	if err != nil {
		return err
	}

	publishDestination := publisher.NewPublishDestination(host, "minio", "minio123", false)

	if err := s.publisherModule.Publish(dir, s.options.Bucket, publishDestination); err != nil {
		return errors.Wrap(err, "cannot publish to server")
	}

	return nil
}
