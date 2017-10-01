package google

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/oauth2/google"
	gce "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"

	clusterv1alpha1 "github.com/jetstack/tarmak/pkg/apis/cluster/v1alpha1"
	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/interfaces"
	"github.com/jetstack/tarmak/pkg/tarmak/utils/input"
)

type Google struct {
	conf *tarmakv1alpha1.Provider

	tarmak interfaces.Tarmak
	log    *logrus.Entry

	gcsClient *storage.Client
	apiClient *http.Client
}

var _ interfaces.Provider = &Google{}

func NewFromConfig(tarmak interfaces.Tarmak, conf *tarmakv1alpha1.Provider) (*Google, error) {
	client, err := google.DefaultClient(context.Background(),
		gce.ComputeReadonlyScope,
		dns.NdevClouddnsReadwriteScope)

	if err != nil {
		return nil, fmt.Errorf("Unable to get Google Cloud client: %v", err)
	}

	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	g := &Google{
		conf:      conf,
		log:       tarmak.Log().WithField("provider_name", conf.ObjectMeta.Name),
		tarmak:    tarmak,
		gcsClient: gcsClient,
		apiClient: client,
	}
	return g, nil
}

func (g *Google) Name() string {
	return g.conf.Name
}

func (g *Google) Cloud() string {
	return clusterv1alpha1.CloudGoogle
}

func (g *Google) Region() string {
	region := g.tarmak.Environment().Location()
	if region == "" {
		return "europe-west1"
	}
	return g.tarmak.Environment().Location()
}

// This parameters should include non sensitive information to identify a provider
func (g *Google) Parameters() map[string]string {
	p := map[string]string{
		"name":          g.Name(),
		"cloud":         g.Cloud(),
		"project":       g.conf.GCP.Project,
		"public_zone":   g.conf.GCP.PublicZone,
		"bucket_prefix": g.conf.GCP.BucketPrefix,
	}
	return p
}

func (g *Google) Validate() error {
	var result error
	var err error

	// These checks only make sense with an environment given
	if g.tarmak.Environment() != nil {
		err = g.validateRemoteStateBucket()
		if err != nil {
			result = multierror.Append(result, err)
		}

		// 	err = a.validateRemoteStateDynamoDB()
		// 	if err != nil {
		// 		result = multierror.Append(err)
		// 	}

		// 	err = a.validateAvailabilityZones()
		// 	if err != nil {
		// 		result = multierror.Append(err)
		// 	}

		// 	err = a.validateAWSKeyPair()
		// 	if err != nil {
		// 		result = multierror.Append(err)
		// 	}

	}

	err = g.validatePublicZone()
	if err != nil {
		result = multierror.Append(err)
	}

	if result != nil {
		return result
	}
	return nil
}

func (g *Google) AskEnvironmentLocation(init interfaces.Initialize) (location string, err error) {
	regions, err := g.ListRegions()
	if err != nil {
		return "", err
	}

	regionPos, err := init.Input().AskSelection(&input.AskSelection{
		Query:   "In which region should this environment reside?",
		Choices: regions,
		Default: -1,
	})
	if err != nil {
		return "", err
	}

	return regions[regionPos], nil
}

func (g *Google) AskInstancePoolLocation(init interfaces.Initialize) (zones []string, err error) {
	// TODO: implement zone Choice in multiSel
	// zones, err := g.AZsForRegion(init.CurrentEnvironment().Location())
	// if err != nil {
	// 	return nil, err
	// }

	multiSel := &input.AskMultipleSelection{
		AskOpen: &input.AskOpen{
			Query:      "Please enter an availability zone(s)",
			AllowEmpty: false,
		},
		Query:   "How many availability zones in the cluster? Availability zones will be added to each instance pool in the cluster.",
		Default: 1,
	}

	return init.Input().AskMultipleSelection(multiSel)
}

func (g *Google) ListRegions() ([]string, error) {
	svc, err := gce.New(g.apiClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Compute Engine service: %v", err)
	}
	list, err := svc.Regions.List(g.conf.GCP.Project).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to list GCE regions: %s", err.Error())
	}
	regions := make([]string, len(list.Items))
	for i, region := range list.Items {
		regions[i] = region.Name
	}
	return regions, nil
}

func (g *Google) AZsForRegion(r string) ([]string, error) {
	svc, err := gce.New(g.apiClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Compute Engine service: %v", err)
	}
	region, err := svc.Regions.Get(g.conf.GCP.Project, r).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to get GCE region '%s': %s", r, err.Error())
	}
	return region.Zones, nil
}

// This will return necessary environment variables
func (g *Google) Environment() ([]string, error) {
	// TODO: correctly set credentials location
	return []string{
		fmt.Sprintf("GOOGLE_CREDENTIALS=%s", ""),
		fmt.Sprintf("GOOGLE_REGION=%s", g.Region()),
	}, nil
}

// This methods converts and possibly validates a generic instance type to a
// provider specifc
func (g *Google) InstanceType(typeIn string) (typeOut string, err error) {
	if typeIn == clusterv1alpha1.InstancePoolSizeTiny {
		return "n1-standard-1", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeSmall {
		return "n1-standard-4", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeMedium {
		return "n1-standard-8", nil
	}
	if typeIn == clusterv1alpha1.InstancePoolSizeLarge {
		return "n1-standard-16", nil
	}

	// TODO: Validate custom instance type here
	return typeIn, nil
}

func (g *Google) String() string {
	return fmt.Sprintf("%s[%s]", g.Cloud(), g.Name())
}

// TODO: implement
func (a *Google) Variables() map[string]interface{} {
	output := map[string]interface{}{}
	// output["key_name"] = a.KeyName()
	// if len(a.conf.Amazon.AllowedAccountIDs) > 0 {
	// 	output["allowed_account_ids"] = a.conf.Amazon.AllowedAccountIDs
	// }
	// output["availability_zones"] = a.AvailabilityZones()
	// output["region"] = a.Region()

	// output["public_zone"] = a.conf.Amazon.PublicZone
	// output["public_zone_id"] = a.conf.Amazon.PublicHostedZoneID
	// output["bucket_prefix"] = a.conf.Amazon.BucketPrefix

	return output
}

// This methods converts and possibly validates a generic volume type to a
// provider specifc
func (g *Google) VolumeType(typeIn string) (typeOut string, err error) {
	if typeIn == clusterv1alpha1.VolumeTypeHDD {
		return "pd-standard", nil
	}
	if typeIn == clusterv1alpha1.VolumeTypeSSD {
		return "pd-ssd", nil
	}
	// TODO: Validate custom instance type here
	return typeIn, nil
}
