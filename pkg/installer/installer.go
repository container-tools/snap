package installer

import (
	"context"
	"errors"
	"fmt"
	"github.com/container-tools/snap/deploy"
	snapclient "github.com/container-tools/snap/pkg/client"
	kubeutils "github.com/container-tools/snap/pkg/util/kubernetes"
	"github.com/container-tools/snap/pkg/util/log"
	"github.com/sethvargo/go-password/password"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	logger = log.WithName("installer")

	serverLabelSelector = "snap.container-tools.io/component=server"
)

const (
	AccessKeyEntry = "access-key"
	SecretKeyEntry = "secret-key"
)

type Installer struct {
	client snapclient.Client
	stdOut io.Writer
	stdErr io.Writer
}

type InstallerSnapCredentials struct {
	SecretName     string
	AccessKeyEntry string
	AccessKey      string
	SecretKeyEntry string
	SecretKey      string
}

func NewInstaller(config *restclient.Config, client ctrl.Client, stdOut, stdErr io.Writer) (*Installer, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Installer{
		client: snapclient.Client{
			Interface: clientset,
			Client:    client,
			Config:    config,
		},
		stdOut: stdOut,
		stdErr: stdErr,
	}, nil
}

func (i *Installer) IsInstalled(ctx context.Context, ns string) (bool, error) {
	deploymentList, err := i.client.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{
		LabelSelector: serverLabelSelector,
	})
	if err != nil {
		return false, err
	}
	return len(deploymentList.Items) > 0, nil
}

func (i *Installer) OpenConnection(ctx context.Context, ns string, direct bool) (string, error) {
	if direct {
		return i.GetDirectConnectionHost(ctx, ns)
	}

	logger.Info("Waiting for destination pod to be ready...")
	pod, err := kubeutils.WaitForPodReady(ctx, i.client, ns, serverLabelSelector)
	if err != nil {
		return "", err
	} else if pod == "" {
		return "", errors.New("cannot find server pod")
	}

	logger.Infof("Opening connection to pod %s", pod)
	host, err := kubeutils.PortForward(ctx, i.client.Config, ns, pod, i.stdOut, i.stdErr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", host), nil
}

func (i *Installer) GetDirectConnectionHost(ctx context.Context, ns string) (string, error) {
	serviceList, err := i.client.CoreV1().Services(ns).List(ctx, metav1.ListOptions{
		LabelSelector: serverLabelSelector,
	})
	if err != nil {
		return "", err
	}
	if len(serviceList.Items) == 0 {
		return "", errors.New("no snap server found")
	}
	return fmt.Sprintf("%s:9000", serviceList.Items[0].Name), nil
}

func (i *Installer) EnsureInstalled(ctx context.Context, ns string) error {
	if installed, err := i.IsInstalled(ctx, ns); err != nil {
		return err
	} else if installed {
		logger.Info("Snap is already installed: skipping")
		return nil
	}

	logger.Infof("Installing Snap into the %s namespace...", ns)
	if err := i.installSecret(ctx, ns); err != nil {
		return err
	}
	if err := i.installResource(ctx, ns, "/minio-standalone-pvc.yaml"); err != nil {
		return err
	}
	if err := i.installResource(ctx, ns, "/minio-standalone-deployment.yaml"); err != nil {
		return err
	}
	if err := i.installResource(ctx, ns, "/minio-standalone-service.yaml"); err != nil {
		return err
	}
	logger.Infof("Installation complete in namespace %s", ns)
	return nil
}

func (i *Installer) GetCredentials(ctx context.Context, ns string) (credentials InstallerSnapCredentials, err error) {
	secrets, err := i.client.CoreV1().Secrets(ns).List(ctx, metav1.ListOptions{
		LabelSelector: serverLabelSelector,
	})
	if err != nil {
		return credentials, err
	}

	if len(secrets.Items) == 0 {
		return credentials, errors.New("no credentials found for the server")
	}
	secret := secrets.Items[0]
	key := string(secret.Data[AccessKeyEntry])
	keySecret := string(secret.Data[SecretKeyEntry])
	if len(key) == 0 || len(keySecret) == 0 {
		return credentials, errors.New("empty credentials found")
	}

	credentials.SecretName = secret.Name
	credentials.AccessKey = key
	credentials.AccessKeyEntry = AccessKeyEntry
	credentials.SecretKey = keySecret
	credentials.SecretKeyEntry = SecretKeyEntry
	return credentials, nil
}

func (i *Installer) installSecret(ctx context.Context, ns string) error {
	obj, err := kubeutils.LoadResourceFromYaml(scheme.Scheme, deploy.ResourceAsString("/minio-standalone-secret.yaml"))
	if err != nil {
		return err
	}
	secret := obj.(*corev1.Secret)

	accessKey, err := password.Generate(64, 10, 0, true, true)
	if err != nil {
		return err
	}
	secretKey, err := password.Generate(64, 10, 0, false, true)
	if err != nil {
		return err
	}

	if secret.StringData == nil {
		secret.StringData = make(map[string]string)
	}
	secret.StringData[AccessKeyEntry] = accessKey
	secret.StringData[SecretKeyEntry] = secretKey

	return kubeutils.ReplaceResourceInNamespace(ctx, i.client, secret, ns)
}

func (i *Installer) installResource(ctx context.Context, ns string, name string) error {
	pvc, err := kubeutils.LoadResourceFromYaml(scheme.Scheme, deploy.ResourceAsString(name))
	if err != nil {
		return err
	}

	return kubeutils.ReplaceResourceInNamespace(ctx, i.client, pvc, ns)
}
