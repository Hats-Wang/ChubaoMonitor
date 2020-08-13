package chubaomonitor

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilintstr "k8s.io/apimachinery/pkg/util/intstr"

	cachev1alpha1 "github.com/ChubaoMonitor/pkg/apis/cache/v1alpha1"
)

func (r *ReconcileChubaoMonitor) Deploymentforprometheus(m *cachev1alpha1.ChubaoMonitor) *appsv1.Deployment {
	name := "prometheus"
	labels := labelsForChubaoMonitor(name)
	selector := &metav1.LabelSelector{MatchLabels: labels}

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &m.Spec.Prometheus.SizeProm,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: containerforprometheus(m),
					Volumes:    volumeforprometheus(m),
				},
			},
			Selector: selector,
		},
	}
	return dep
}

func containerforprometheus(m *cachev1alpha1.ChubaoMonitor) []corev1.Container {

	return []corev1.Container{
		{
			Name:  "prometheus-pod",
			Image: m.Spec.Prometheus.ImageProm,
			Ports: []corev1.ContainerPort{
				{ContainerPort: m.Spec.Prometheus.PortProm},
			},
			ImagePullPolicy: m.Spec.Prometheus.ImagePullPolicyProm,
			Env: []corev1.EnvVar{
				{Name: "TZ", Value: "Asia/Shanghai"},
			},
			VolumeMounts: volumemountsforprometheus(),
			Resources:    m.Spec.Prometheus.ResourcesProm,
		},
	}
}

func Serviceforprometheus(m *cachev1alpha1.ChubaoMonitor) *corev1.Service {
	name := "prometheus"
	labels := labelsForChubaoMonitor(name)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       m.Spec.Prometheus.PortProm,
					TargetPort: utilintstr.IntOrString{IntVal: m.Spec.Prometheus.PortProm, Type: utilintstr.Int},
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func volumeforprometheus(m *cachev1alpha1.ChubaoMonitor) []corev1.Volume {
	var defaultmode int32 = 0555
	return []corev1.Volume{
		{
			Name: "monitor-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "monitor-config",
					},
					DefaultMode: &defaultmode,
				},
			},
		},
		{
			Name: "prometheus-data",
			VolumeSource: corev1.VolumeSource{
				HostPath: m.Spec.Prometheus.HostPath,
			},
		},
	}
}

func volumemountsforprometheus() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "monitor-config",
			MountPath: "/etc/prometheus/prometheus.yml",
			SubPath:   "prometheus.yml",
		},
		{
			Name:      "prometheus-data",
			MountPath: "/prometheus-data",
		},
	}
}

func (r *ReconcileChubaoMonitor) Deploymentforgrafana(m *cachev1alpha1.ChubaoMonitor) *appsv1.Deployment {
	name := "grafana"
	labels := labelsForChubaoMonitor(name)
	selector := &metav1.LabelSelector{MatchLabels: labels}
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana",
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &m.Spec.Grafana.SizeGrafana,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: containerforgrafana(m),
					Volumes:    volumesforgrafana(m),
				},
			},
			Selector: selector,
		},
	}
	return dep
}

func containerforgrafana(m *cachev1alpha1.ChubaoMonitor) []corev1.Container {

	var privileged bool = true
	return []corev1.Container{
		{
			Name:  "grafana-pod",
			Image: m.Spec.Grafana.ImageGrafana,
			Ports: []corev1.ContainerPort{
				{ContainerPort: m.Spec.Grafana.PortGrafana},
			},
			ImagePullPolicy: m.Spec.Grafana.ImagePullPolicyGrafana,
			SecurityContext: &corev1.SecurityContext{
				Privileged: &privileged,
			},
			Env: envforgrafana(),
			// If grafana pod show the err "back-off restarting failed container", run this command to keep the container running ang then run ./run.sh in the container to check the really error.
			//          Command:        []string{"/bin/bash", "-ce", "tail -f /dev/null"},
			ReadinessProbe: readinessforgrafana(),
			VolumeMounts:   volumemountsforgrafana(),
			Resources:      m.Spec.Grafana.ResourcesGrafana,
		},
	}
}

func Serviceforgrafana(m *cachev1alpha1.ChubaoMonitor) *corev1.Service {
	name := "grafana"
	labels := labelsForChubaoMonitor(name)
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       m.Spec.Grafana.PortGrafana,
					TargetPort: utilintstr.IntOrString{IntVal: m.Spec.Grafana.PortGrafana, Type: utilintstr.Int},
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func volumesforgrafana(m *cachev1alpha1.ChubaoMonitor) []corev1.Volume {
	var defaultmode int32 = 0555

	return []corev1.Volume{
		{
			Name: "monitor-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "monitor-config",
					},
					DefaultMode: &defaultmode,
				},
			},
		},
		{Name: "grafana-persistent-storage",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
}

func volumemountsforgrafana() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "grafana-persistent-storage",
			MountPath: "/var/lib/grafana",
		},
		{
			Name:      "monitor-config",
			MountPath: "/grafana/init.sh",
			SubPath:   "init.sh",
		},
		{
			Name:      "monitor-config",
			MountPath: "/etc/grafana/grafana.ini",
			SubPath:   "grafana.ini",
		},
		{
			Name:      "monitor-config",
			MountPath: "/etc/grafana/provisioning/dashboards/chubaofs.json",
			SubPath:   "chubaofs.json",
		},
		{
			Name:      "monitor-config",
			MountPath: "/etc/grafana/provisioning/dashboards/dashboard.yml",
			SubPath:   "dashboard.yml",
		},
		{
			Name:      "monitor-config",
			MountPath: "/etc/grafana/provisioning/datasources/datasource.yml",
			SubPath:   "datasource.yml",
		},
	}
}

func envforgrafana() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "GF_AUTH_BASIC_ENABLED",
			Value: "true",
		},
		{
			Name:  "GF_AUTH_ANONYMOUS_ENABLED",
			Value: "false",
		},
		{
			Name:  "GF_SECURITY_ADMIN_PASSWORD",
			Value: "123456",
		},
	}
}

func readinessforgrafana() *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/login",
				Port: utilintstr.IntOrString{
					IntVal: 3000,
				},
			},
		},
	}
}

func labelsForChubaoMonitor(name string) map[string]string {
	return map[string]string{"app": name}
}
