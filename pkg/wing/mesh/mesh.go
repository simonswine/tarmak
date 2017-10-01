package mesh

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

type Mesh struct {
	router *mesh.Router
	flags  *viper.Viper
	log    *logrus.Entry
	self   mesh.PeerName
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

func New(flags *viper.Viper, log logrus.Entry) *Mesh {
	return &Mesh{
		log:   log.WithField("tier", "mesh"),
		flags: flags,
	}
}

func (m *Mesh) Start() error {
	var err error

	err = m.init()
	if err != nil {
		return fmt.Errorf("unable to initialise mesh: %s", err)
	}

	m.router.Start()
	m.log.Debugf("mesh router started")

	return nil
}

func (m *Mesh) init() error {
	var err error

	m.self, err = mesh.PeerNameFromString(m.flags.GetString(FlagMeshHWAddr))
	if err != nil {
		return fmt.Errorf("could not create peer name: %v", err)
	}

	m.router, err = mesh.NewRouter(
		mesh.Config{
			Host:               m.flags.GetString(FlagMeshListenHost),
			Port:               m.flags.GetInt(FlagMeshListenPort),
			ProtocolMinVersion: mesh.ProtocolMinVersion,
			Password:           []byte(m.flags.GetString(FlagMeshPassword)),
			ConnLimit:          64,
			PeerDiscovery:      true,
			TrustedSubnets:     []*net.IPNet{},
		},
		m.self,
		m.flags.GetString(FlagMeshHostname),
		mesh.NullOverlay{},
		log.New(ioutil.Discard, "", 0),
	)
	if err != nil {
		return fmt.Errorf("could not create router: %v", err)
	}

	/*peer := m.newPeer()
	gossip, err := m.router.NewGossip(m.flags.GetString(FlagMeshChannel), peer)
	if err != nil {
		logger.Fatalf("Could not create gossip: %v", err)
	}
	*/

	return nil
}
