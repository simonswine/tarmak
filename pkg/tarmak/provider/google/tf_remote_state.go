package google

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// TODO: remove me, deprecated
func (g *Google) RemoteStateBucketName() string {
	return g.RemoteStateName()
}

func (g *Google) RemoteStateName() string {
	return fmt.Sprintf(
		"%s%s-terraform-state",
		g.conf.GCP.BucketPrefix,
		g.Region(),
	)
}

func (g *Google) RemoteState(namespace string, clusterName string, stackName string) string {
	return fmt.Sprintf(`terraform {
  backend "gcs" {
    bucket  = "%s"
    path    = "%s"
    project = "%s"
  }
}`,
		g.RemoteStateName(),
		fmt.Sprintf("%s/%s/%s.tfstate", namespace, clusterName, stackName),
		g.conf.GCP.Project,
	)
}

func (g *Google) RemoteStateBucketAvailable() (bool, error) {
	ctx := context.Background()
	_, err := g.gcsClient.Bucket(g.RemoteStateName()).Attrs(ctx)

	if err == storage.ErrBucketNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return false, fmt.Errorf("error while checking if remote state is available: %s", err)
}

func (g *Google) validateRemoteStateBucket() error {
	ctx := context.Background()
	attrs, err := g.gcsClient.Bucket(g.RemoteStateName()).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		return g.initRemoteStateBucket()
	}
	if err != nil {
		return fmt.Errorf("error looking for terraform state bucket: %s", err)
	}

	if bucketRegion, myRegion := attrs.Location, g.Region(); bucketRegion != myRegion {
		return fmt.Errorf("bucket region is wrong, actual: %s expected: %s", bucketRegion, myRegion)
	}

	if !attrs.VersioningEnabled {
		g.log.Warnf("state bucket %s has versioning disabled", g.RemoteStateName())
	}

	return nil
}

func (g *Google) initRemoteStateBucket() error {
	bucket := g.gcsClient.Bucket(g.RemoteStateName())

	ctx := context.Background()
	// TODO: Set up ACLs
	return bucket.Create(ctx, g.conf.GCP.Project, &storage.BucketAttrs{
		VersioningEnabled: true,
	})
}
