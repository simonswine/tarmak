package google

import (
	"fmt"
	"strings"

	dns "google.golang.org/api/dns/v1"
	googleapi "google.golang.org/api/googleapi"
)

func (g *Google) PublicZone() string {
	return g.conf.GCP.PublicZone
}

// this removes an ending . in zone and converts it to lowercase, and converts
// . to -
func normalizeZone(in string) string {
	return strings.Replace(strings.ToLower(strings.TrimRight(in, ".")), ".", "-", -1)
}

func (g *Google) initPublicZone() (*dns.ManagedZone, error) {
	svc, err := dns.New(g.apiClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Cloud DNS service: %v", err)
	}

	zone, err := svc.ManagedZones.Create(g.conf.GCP.Project, &dns.ManagedZone{
		Name:        normalizeZone(g.conf.GCP.PublicZone),
		DnsName:     g.conf.GCP.PublicZone,
		Description: "public zone for tarmak",
	}).Do()

	if err != nil {
		return nil, fmt.Errorf("error creating DNS zone: %s", err.Error())
	}

	return zone, nil
}

func (g *Google) validatePublicZone() error {
	svc, err := dns.New(g.apiClient)
	if err != nil {
		return fmt.Errorf("unable to create Google Cloud DNS service: %v", err)
	}

	_, err = svc.ManagedZones.Get(g.conf.GCP.Project, normalizeZone(g.conf.GCP.PublicZone)).Do()
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 404 {
			_, err := g.initPublicZone()
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("error validating DNS zone: %s", err.Error())
	}
	return nil
}
