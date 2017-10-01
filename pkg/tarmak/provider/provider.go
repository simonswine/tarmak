package provider

import (
	"fmt"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/provider/amazon"
	"github.com/jetstack/tarmak/pkg/tarmak/provider/google"
)

func NewFromConfig(tarmak interfaces.Tarmak, conf *tarmakv1alpha1.Provider) (interfaces.Provider, error) {
	var provider interfaces.Provider
	var err error

	switch {
	case conf.Amazon != nil:
		provider, err = amazon.NewFromConfig(tarmak, conf)
	case conf.GCP != nil:
		provider, err = google.NewFromConfig(tarmak, conf)
	default:
		return nil, fmt.Errorf("unknown provider: '%s'", conf.Name)
	}

	if err != nil {
		return provider, fmt.Errorf("error creating provider: %s", err.Error())
	}

	return provider, nil
}
