package google

import (
	"fmt"
	"time"

	tarmakv1alpha1 "github.com/jetstack/tarmak/pkg/apis/tarmak/v1alpha1"

	gce "google.golang.org/api/compute/v0.beta"
)

func buildFilter(name, op, expr string) string {
	return fmt.Sprintf("%s %s %s", name, op, expr)
}

func (g *Google) QueryImages(tags map[string]string) (images []tarmakv1alpha1.Image, err error) {
	svc, err := gce.New(g.apiClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Compute Engine service: %v", err)
	}

	listCall := svc.Images.List(g.conf.GCP.Project)
	for key, val := range tags {
		listCall = listCall.Filter(buildFilter(key, "eq", val))
	}

	list, err := listCall.Do()
	if err != nil {
		return nil, fmt.Errorf("error listing GCE images: %s", err.Error())
	}

	formatRFC3339 := "2006-01-02T15:04:05.999Z07:00"

	for _, img := range list.Items {
		image := tarmakv1alpha1.Image{}
		image.Annotations = map[string]string{}
		image.Labels = img.Labels
		if baseImage, ok := image.Labels[tarmakv1alpha1.ImageTagBaseImageName]; ok {
			image.BaseImage = baseImage
		}
		creationTimestamp, err := time.Parse(formatRFC3339, img.CreationTimestamp)
		if err != nil {
			return images, fmt.Errorf("error parsing time stamp '%s'", err)
		}
		image.CreationTimestamp.Time = creationTimestamp
		image.Name = img.Name
		images = append(images, image)
	}

	return images, nil
}
