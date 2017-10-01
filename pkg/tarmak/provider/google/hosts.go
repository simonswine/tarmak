package google

import (
	"fmt"
	"strings"

	gce "google.golang.org/api/compute/v1"

	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
)

type host struct {
	id             string
	host           string
	hostnamePublic bool
	hostname       string
	aliases        []string
	roles          []string
	user           string

	cluster interfaces.Cluster
}

var _ interfaces.Host = &host{}

func (h *host) ID() string {
	return h.id
}

func (h *host) Roles() []string {
	return h.roles
}

func (h *host) Aliases() []string {
	return h.aliases
}

func (h *host) Hostname() string {
	return h.hostname
}

func (h *host) HostnamePublic() bool {
	return h.hostnamePublic
}

func (h *host) User() string {
	return h.user
}

func (h *host) Parameters() map[string]string {
	return map[string]string{
		"id":       h.ID(),
		"hostname": h.Hostname(),
		"roles":    strings.Join(h.Roles(), ", "),
	}
}

func (h *host) SSHConfig() string {
	config := fmt.Sprintf(`host %s
    User %s
    Hostname %s

    # use custom host key file per cluster
    UserKnownHostsFile %s
    StrictHostKeyChecking no

    # enable connection multiplexing
    ControlPath %s/ssh-control-%%r@%%h:%%p
    ControlMaster auto
    ControlPersist 10m

    # keep connections alive
    ServerAliveInterval 60
    IdentitiesOnly yes
    IdentityFile %s
`,
		strings.Join(append(h.Aliases(), h.ID()), " "),
		h.User(),
		h.Hostname(),
		h.cluster.SSHHostKeysPath(),
		h.cluster.ConfigPath(),
		h.cluster.Environment().SSHPrivateKeyPath(),
	)

	if !h.HostnamePublic() {
		config += fmt.Sprintf(
			"    ProxyCommand ssh -F %s -W %%h:%%p bastion\n",
			h.cluster.SSHConfigPath(),
		)
	}
	config += "\n"
	return config
}

func (g *Google) ListHosts() ([]interfaces.Host, error) {
	svc, err := gce.New(g.apiClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Compute Engine service: %v", err)
	}

	filter := buildFilter("tarmak-environment", "eq", g.tarmak.Cluster().Environment().Name())
	list, err := svc.Instances.AggregatedList(g.conf.GCP.Project).Filter(filter).Do()
	if err != nil {
		return nil, err
	}

	hosts := []*host{}

	for _, scopedInstances := range list.Items {
		for _, instance := range scopedInstances.Instances {
			host := &host{}
			for _, nic := range instance.NetworkInterfaces {
				if nic.NetworkIP == "" {
					continue
				}
				host.hostname = nic.NetworkIP
				for _, accessConfig := range nic.AccessConfigs {
					if accessConfig.Name == "external-nat" {
						if accessConfig.NatIP != "" {
							host.hostname = accessConfig.NatIP
							host.hostnamePublic = true
						}
					}
				}
				break
			}
			for _, tag := range instance.Metadata.Items {
				if tag.Key == "tarmak_role" {
					host.roles = strings.Split(*tag.Value, ",")
				}
			}
			hosts = append(hosts, host)
		}
	}

	hostsByRole := map[string][]*host{}
	for _, h := range hosts {
		for _, role := range h.roles {
			if _, ok := hostsByRole[role]; !ok {
				hostsByRole[role] = []*host{h}
			} else {
				hostsByRole[role] = append(hostsByRole[role], h)
			}
			h.aliases = append(h.aliases, fmt.Sprintf("%s-%d", role, len(hostsByRole[role])))
		}
	}

	// remove role-1 for single instances
	for role, hosts := range hostsByRole {
		if len(hosts) != 1 {
			continue
		}
		for pos, _ := range hosts[0].aliases {
			if hosts[0].aliases[pos] == fmt.Sprintf("%s-1", role) {
				hosts[0].aliases[pos] = role
			}
		}
	}

	hostsInterfaces := make([]interfaces.Host, len(hosts))

	for pos, _ := range hosts {
		hostsInterfaces[pos] = hosts[pos]
	}

	return hostsInterfaces, nil
}
