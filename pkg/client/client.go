package client

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/user"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	restclient "k8s.io/client-go/rest"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Client struct {
	kubernetes.Interface
	ctrl.Client
	Config *restclient.Config
}

// NewClient creates a new k8s client that can be used from outside or in the cluster
func NewRestConfig() (*rest.Config, error) {
	initialize()
	// Get a config to talk to the apiserver
	return config.GetConfig()
}

// init initialize the k8s client for usage outside the cluster
func initialize() {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		var err error
		kubeconfig, err = getDefaultKubeConfigFile()
		if err != nil {
			panic(err)
		}
	}
	os.Setenv("KUBECONFIG", kubeconfig)
}

func getDefaultKubeConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".kube", "config"), nil
}

// GetCurrentNamespace --
func GetCurrentNamespace() (string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		var err error
		kubeconfig, err = getDefaultKubeConfigFile()
		if err != nil {
			logrus.Errorf("Cannot get information about current user: %v", err)
		}
	}
	if kubeconfig == "" {
		return "default", nil
	}

	data, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return "", err
	}
	conf := clientcmdapi.NewConfig()
	if len(data) == 0 {
		return "", errors.New("kubernetes config file is empty")
	}

	decoded, _, err := clientcmdlatest.Codec.Decode(data, &schema.GroupVersionKind{Version: clientcmdlatest.Version, Kind: "Config"}, conf)
	if err != nil {
		return "", err
	}

	clientcmdconfig := decoded.(*clientcmdapi.Config)

	cc := clientcmd.NewDefaultClientConfig(*clientcmdconfig, &clientcmd.ConfigOverrides{})
	ns, _, err := cc.Namespace()
	return ns, err
}
