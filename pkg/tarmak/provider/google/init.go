package google

import (
	"fmt"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"
	"github.com/jetstack/tarmak/pkg/tarmak/utils/input"
)

func Init(in *input.Input, provider *tarmakv1alpha1.Provider) error {
	if provider.GCP == nil {
		provider.GCP = &tarmakv1alpha1.ProviderGCP{}
	}
	err := initCredentials(in, provider)
	if err != nil {
		return err
	}

	err = initBucketPrefix(in, provider)
	if err != nil {
		return err
	}

	err = initPublicZone(in, provider)
	if err != nil {
		return err
	}

	return nil
}

func initBucketPrefix(in *input.Input, provider *tarmakv1alpha1.Provider) error {
	for {
		bucketPrefix, err := in.AskOpen(&input.AskOpen{
			Query:   "Which prefix should be used for the state buckets? ([a-z0-9-]+, should be globally unique)",
			Default: fmt.Sprintf("%s-tarmak-", provider.Name),
		})
		if err != nil {
			return err
		}

		nameValid := input.RegexpName.MatchString(bucketPrefix)

		if !nameValid {
			in.Warnf("bucket/table prefix '%s' is not valid", bucketPrefix)
		} else {
			provider.GCP.BucketPrefix = bucketPrefix
			break
		}
	}

	return nil
}

func initPublicZone(in *input.Input, provider *tarmakv1alpha1.Provider) error {
	for {
		publicZone, err := in.AskOpen(&input.AskOpen{
			Query: "Which public DNS zone should be used? (DNS zone will be created if not existing, it needs to be delegated from the Root)",
		})
		if err != nil {
			return err
		}

		zoneValid := input.RegexpDNS.MatchString(publicZone)

		if !zoneValid {
			in.Warnf("Public DNS zone '%s' is not valid", publicZone)
		} else {
			provider.GCP.PublicZone = publicZone
			break
		}
	}

	return nil
}

func initCredentials(in *input.Input, provider *tarmakv1alpha1.Provider) error {
	for {
		projectName, err := in.AskOpen(&input.AskOpen{
			Query: "What is the Google Cloud project ID?",
		})
		if err != nil {
			return err
		}

		provider.GCP.Project = projectName
		break
	}

	return nil
}
