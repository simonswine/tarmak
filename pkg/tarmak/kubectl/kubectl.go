package kubectl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

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
	return filepath.Join(k.tarmak.Context().ConfigPath(), "kubeconfig")
}

func (k *Kubectl) requestNewAdminCert(cluster *api.Cluster, authInfo *api.AuthInfo) error {
	path := fmt.Sprintf("%s/pki/k8s/sign/admin", k.tarmak.Context().ContextName())

	k.log.Infof("request new certificate from vault (%s)", path)

	// read vault root token
	vaultRootToken, err := k.tarmak.Context().Environment().VaultRootToken()
	if err != nil {
		return err
	}

	// init vault statck
	_, err = k.tarmak.Terraform().Output(k.tarmak.Context().Environment().VaultStack())
	if err != nil {
		return err
	}

	vaultTunnel, err := k.tarmak.Context().Environment().VaultTunnel()
	if err != nil {
		return err
	}
	defer vaultTunnel.Stop()

	v := vaultTunnel.VaultClient()
	v.SetToken(vaultRootToken)

	// generate new RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("unable to generate private key: %s", err)
	}
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// define CSR template
	var csrTemplate = x509.CertificateRequest{
		Subject:            pkix.Name{CommonName: "admin"},
		SignatureAlgorithm: x509.SHA512WithRSA,
	}

	// generate the CSR request
	csr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKey)
	if err != nil {
		return err
	}

	// pem encode CSR
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csr,
	})

	inputData := map[string]interface{}{
		"csr":         string(csrPem),
		"common_name": "admin",
	}

	output, err := v.Logical().Write(path, inputData)
	if err != nil {
		return err
	}

	certPemIntf, ok := output.Data["certificate"]
	if !ok {
		return errors.New("key certificate not found")
	}

	certPem, ok := certPemIntf.(string)
	if !ok {
		return fmt.Errorf("certificate has unexpected type %s", certPemIntf)
	}

	caPemIntf, ok := output.Data["issuing_ca"]
	if !ok {
		return errors.New("issuing_ca not found")
	}

	caPem, ok := caPemIntf.(string)
	if !ok {
		return fmt.Errorf("issuing_ca has unexpected type %s", caPemIntf)
	}

	authInfo.ClientKeyData = privateKeyPem
	authInfo.ClientCertificateData = []byte(certPem)
	cluster.CertificateAuthorityData = []byte(caPem)

	return nil
}

func (k *Kubectl) EnsureConfig() error {
	c := api.NewConfig()
	configPath := k.ConfigPath()

	// context name in tamrak is context name in kubeconfig
	key := k.tarmak.Context().ContextName()

	// load an existing config
	if _, err := os.Stat(configPath); err == nil {
		conf, err := clientcmd.LoadFromFile(configPath)
		if err != nil {
			return err
		}
		c = conf
	}

	newTunnel := false

	c.CurrentContext = key

	ctx, ok := c.Contexts[key]
	if !ok {
		ctx = api.NewContext()
		ctx.Namespace = "kube-system"
		ctx.Cluster = key
		ctx.AuthInfo = key
		c.Contexts[key] = ctx
	}

	cluster, ok := c.Clusters[key]
	if !ok {
		newTunnel = true
		cluster = api.NewCluster()
		cluster.CertificateAuthorityData = []byte{}
		cluster.Server = ""
		c.Clusters[key] = cluster
	}

	authInfo, ok := c.AuthInfos[key]
	if !ok {
		authInfo = api.NewAuthInfo()
		authInfo.ClientCertificateData = []byte{}
		authInfo.ClientKeyData = []byte{}
		c.AuthInfos[key] = authInfo
	}

	k.log.Infof("%#+v", c)
	return fmt.Errorf("xx")

	// check if certificates are set
	if len(authInfo.ClientCertificateData) == 0 || len(authInfo.ClientKeyData) == 0 || len(cluster.CertificateAuthorityData) == 0 {
		if err := k.requestNewAdminCert(cluster, authInfo); err != nil {
			return err
		}
	}

	if !newTunnel {
		// test connectivity without setting up tunnel
	}

	// setup kube api tunnel

	// test connectivity

	err := utils.EnsureDirectory(filepath.Dir(configPath), 0700)
	if err != nil {
		return err
	}

	err = clientcmd.WriteToFile(*c, configPath)

	return fmt.Errorf("xx")
}
