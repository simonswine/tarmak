package puppet

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"

	clusterv1alpha1 "github.com/jetstack/tarmak/pkg/apis/cluster/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
)

type Puppet struct {
	log    *logrus.Entry
	tarmak interfaces.Tarmak
}

func New(tarmak interfaces.Tarmak) *Puppet {
	log := tarmak.Log().WithField("module", "puppet")

	return &Puppet{
		log:    log,
		tarmak: tarmak,
	}
}

func (p *Puppet) TarGz(writer io.Writer) error {

	rootPath, err := p.tarmak.RootPath()
	if err != nil {
		return fmt.Errorf("error getting rootPath: %s", err)
	}

	path := filepath.Join(rootPath, "puppet")

	reader, err := archive.Tar(
		path,
		archive.Gzip,
	)
	if err != nil {
		return fmt.Errorf("error creating tar from path '%s': %s", path, err)
	}

	if _, err := io.Copy(writer, reader); err != nil {
		return fmt.Errorf("error writing tar: %s", err)
	}

	return nil
}

func kubernetesConfig(conf *clusterv1alpha1.Kubernetes) (lines []string) {
	if conf == nil {
		return
	}
	if conf.Version != "" {
		lines = append(lines, fmt.Sprintf(`tarmak::kubernetes_version: "%s"`, conf.Kubernetes.Version))
	}
}

func contentGlobalConfig(conf *clusterv1alpha1.Cluster) (lines []string) {
	lines = append(lines, kubernetesConfig(conf.Kubernetes)...)
	return lines
}

func contentInstancePoolConfig(conf *clusterv1alpha1.ServerPool) (lines []string) {
	lines = append(lines, kubernetesConfig(conf.Kubernetes)...)
	return lines
}

func (p *Puppet) HieraData(context interfaces.Context) {

	// get global cluster config
	clusterCfg := context.Config()

	// loop through instance pools
	for _, instancePool := range context.NodeGroups() {
		cfg := instancePool.Config()
	}

}
