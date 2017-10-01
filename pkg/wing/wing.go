package wing

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weaveworks/mesh"
)

const (
	FlagMeshListenHost = "mesh.listen_host"
	FlagMeshListenPort = "mesh.listen_port"
	FlagMeshPassword   = "mesh.password"
	FlagMeshChannel    = "mesh.channel"
	FlagMeshHWAddr     = "mesh.hw_addr"
	FlagMeshHostname   = "mesh.hostname"
)

type Wing struct {
	log   *logrus.Entry
	flags *viper.Viper

	router *mesh.Router
}

type Provider interface {
	ID() string
}

func getHWAddr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if s := iface.HardwareAddr.String(); s != "" {
			return s, nil
		}
	}
	return "", errors.New("no valid network interface found")
}

func DefaultFlags() (keys []string, values []interface{}, descriptions []string, err error) {
	hwAddr, err := getHWAddr()
	if err != nil {
		return keys, values, descriptions, err
	}
	keys = append(keys, FlagMeshHWAddr)
	values = append(values, hwAddr)
	descriptions = append(descriptions, "ID used for mesh network")

	hostname, err := os.Hostname()
	if err != nil {
		return keys, values, descriptions, err
	}
	keys = append(keys, FlagMeshHostname)
	values = append(values, hostname)
	descriptions = append(descriptions, "Hostname used for mesh alias")

	keys = append(keys, FlagMeshListenHost)
	values = append(values, "0.0.0.0")
	descriptions = append(descriptions, "Listen address for mesh TCP port")

	keys = append(keys, FlagMeshListenPort)
	values = append(values, mesh.Port)
	descriptions = append(descriptions, "Listen port for mesh TCP port")

	keys = append(keys, FlagMeshChannel)
	values = append(values, "tarmak")
	descriptions = append(descriptions, "Gossip channel")

	return keys, values, descriptions, nil
}

func New(flags *viper.Viper) *Wing {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	t := &Wing{
		log:   logger.WithField("app", "wing"),
		flags: flags,
	}
	return t
}

func (w *Wing) Must(err error) *Wing {
	if err != nil {
		w.log.Fatal(err)
	}
	return w
}

func (w *Wing) Start() error {
	var err error

	err = w.initMesh()
	if err != nil {
		return fmt.Errorf("unable to initialise mesh: %s", err)
	}

	w.router.Start()
	w.log.Debugf("mesh router started")

	return nil
}

func (w *Wing) initMesh() error {
	var err error

	peerName, err := mesh.PeerNameFromString(w.flags.GetString(FlagMeshHWAddr))
	if err != nil {
		return fmt.Errorf("could not create peer name: %v", err)
	}

	w.router, err = mesh.NewRouter(
		mesh.Config{
			Host:               w.flags.GetString(FlagMeshListenHost),
			Port:               w.flags.GetInt(FlagMeshListenPort),
			ProtocolMinVersion: mesh.ProtocolMinVersion,
			Password:           []byte(w.flags.GetString(FlagMeshPassword)),
			ConnLimit:          64,
			PeerDiscovery:      true,
			TrustedSubnets:     []*net.IPNet{},
		},
		peerName,
		w.flags.GetString(FlagMeshHostname),
		mesh.NullOverlay{},
		log.New(ioutil.Discard, "", 0),
	)
	if err != nil {
		return fmt.Errorf("could not create router: %v", err)
	}

	return nil
}
