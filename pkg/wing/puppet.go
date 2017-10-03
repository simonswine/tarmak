package wing

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/jetstack/tarmak/pkg/apis/wing/v1alpha1"
)

// This make sure puppet is converged when neccessary
func (w *Wing) convergeLoop() {

	status := &v1alpha1.InstanceStatus{
		Converge: &v1alpha1.InstanceStatusManifest{
			State: "wtf",
		},
	}

	err := w.reportStatus(status)
	if err != nil {
		w.log.Warn("reporting status failed: ", err)
	}
}

// report status to the API server
func (w *Wing) reportStatus(status *v1alpha1.InstanceStatus) error {
	instanceAPI := w.clientset.WingV1alpha1().Instances(w.flags.ClusterName)
	instance, err := instanceAPI.Get(
		w.flags.InstanceName,
		metav1.GetOptions{},
	)
	if err != nil {
		if kerr, ok := err.(*apierrors.StatusError); ok && kerr.ErrStatus.Reason == metav1.StatusReasonNotFound {
			instance = &v1alpha1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name: w.flags.InstanceName,
				},
				Status: status.DeepCopy(),
			}
			_, err := instanceAPI.Create(instance)
			if err != nil {
				return fmt.Errorf("error creating instance: %s", err)
			}
			return nil
		}
		return fmt.Errorf("error get existing instance: %s", err)
	}

	instance.Status = status.DeepCopy()
	_, err = instanceAPI.Update(instance)
	if err != nil {
		return fmt.Errorf("error updating existing instance: %s", err)
		// TODO: handle race for update
	}

	return nil

}
