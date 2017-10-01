package google

import (
	"fmt"

	"github.com/jetstack-experimental/vault-unsealer/pkg/kv"
)

func (g *Google) VaultKV() (kv.Service, error) {
	return nil, fmt.Errorf("not implemented")
}
