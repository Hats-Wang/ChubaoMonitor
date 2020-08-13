package chubaomonitor

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
)

// Compare manually because of pointer to prevent update deployment when Recouncil is recalled
func CompareDeployment(a *appsv1.Deployment, b *appsv1.Deployment) bool {

	switch {
	case *a.Spec.Replicas != *b.Spec.Replicas:
		return true
	case !reflect.DeepEqual(a.Spec.Selector, b.Spec.Selector):
		return true
	case a.Spec.Template.Spec.Containers[0].Image != b.Spec.Template.Spec.Containers[0].Image:
		return true
	case a.Spec.Template.Spec.Containers[0].ImagePullPolicy != b.Spec.Template.Spec.Containers[0].ImagePullPolicy:
		return true
	case a.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort != b.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort:
		return true
	case !reflect.DeepEqual(a.Spec.Template.Spec.Containers[0].Env, b.Spec.Template.Spec.Containers[0].Env):
		return true
	case !reflect.DeepEqual(a.Spec.Template.Spec.Containers[0].VolumeMounts, b.Spec.Template.Spec.Containers[0].VolumeMounts):
		return true

	default:
		return false
	}
}

func CompareService(a *corev1.Service, b *corev1.Service) bool {
	switch {
	case !reflect.DeepEqual(a.Spec.Ports, b.Spec.Ports):
		return true
	case a.Spec.Type != corev1.ServiceTypeClusterIP:
		return true
	case !reflect.DeepEqual(a.Spec.Selector, b.Spec.Selector):
		return true

	default:
		return false
	}
}
