package kubectl

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/jetstack/tarmak/pkg/tarmak/errors"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils"
)

var _ interfaces.Kubectl = &Kubectl{}

type Kubectl struct {
	tarmak interfaces.Tarmak
	log    *logrus.Entry
}

func New(tarmak interfaces.Tarmak) *Kubectl {
	k := &Kubectl{
		tarmak: tarmak,
		log:    tarmak.Log(),
	}

	return k
}

func (k *Kubectl) ConfigPath() string {
	return filepath.Join(k.tarmak.ConfigPath(), "kubeconfig")
}

func (k *Kubectl) EnsureConfig() error {
	c := api.NewConfig()
	configPath := k.ConfigPath()

	// context name in tamrak is context name in kubeconfig
	key := k.tarmak.Context().ContextName()

	// load an existing config
	if _, err := os.Stat(configPath); err != nil {
		conf, err := clientcmd.LoadFromFile(configPath)
		if err != nil {
			return err
		}
		c = conf
	}

	c.CurrentContext = key

	ctx := api.NewContext()
	ctx.Namespace = "kube-system"
	ctx.Cluster = key
	ctx.AuthInfo = key
	c.Contexts[key] = ctx

	cluster := api.NewCluster()
	cluster.CertificateAuthorityData = []byte{}
	cluster.Server = ""
	c.Clusters[key] = cluster

	authInfo := api.NewAuthInfo()
	authInfo.ClientCertificateData = []byte{}
	authInfo.ClientKeyData = []byte{}
	c.AuthInfos[key] = authInfo

	// check if certificates are set
	if len(authInfo.ClientCertificateData) == 0 || len(authInfo.ClientKeyData) == 0 || len(cluster.CertificateAuthorityData) == 0 {
		// TODO: Ask vault for new certificate
		return errors.NotImplemented
	}

	err := utils.EnsureDirectory(filepath.Dir(configPath), 0700)
	if err != nil {
		return err
	}

	err = clientcmd.WriteToFile(*c, configPath)
	return err
}
