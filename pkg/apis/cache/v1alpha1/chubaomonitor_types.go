package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ChubaoMonitorSpec defines the desired state of ChubaoMonitor
// +k8s:openapi-gen=true

type Prometheus struct {
	SizeProm            int32                        `json:"sizeprom"`
	ImageProm           string                       `json:"imageprom"`
	PortProm            int32                        `json:"portprom"`
	ImagePullPolicyProm corev1.PullPolicy            `json:"imagePullPolicyprom,omitempty"`
	ResourcesProm       corev1.ResourceRequirements  `json:"resourcesprom,omitempty"`
	HostPath            *corev1.HostPathVolumeSource `json:"hostPath,omitempty"`
}

type Grafana struct {
	SizeGrafana            int32                       `json:"sizegrafana"`
	ImageGrafana           string                      `json:"imagegrafana"`
	PortGrafana            int32                       `json:"portgrafana"`
	ImagePullPolicyGrafana corev1.PullPolicy           `json:"imagePullPolicygrafana,omitempty"`
	ResourcesGrafana       corev1.ResourceRequirements `json:"resourcesgrafana,omitempty"`
}

type ChubaoMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Prometheus Prometheus `json:"prometheus,omitempty"`
	Grafana    Grafana    `json:"grafana,omitempty"`
}

// ChubaoMonitorStatus defines the observed state of ChubaoMonitor
// +k8s:openapi-gen=true
type ChubaoMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	PrometheusReplicas  int32    `json:"prometheusreplicas,omitempty"`
	GrafanaReplicas     int32    `json:"grafanareplicas,omitempty"`
	PrometheusPods      []string `json:"prometiheuspods,omitempty"`
	GrafanaPods         []string `json:"grafanapods,omitempty"`
	GrafanaclusterIP    string   `json:"grafanaip,omitempty"`
	PrometheusclusterIP string   `json:"prometheusip,omitempty"`
	Configmapstatus     bool     `json:"configmapstatus,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChubaoMonitor is the Schema for the chubaomonitors API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=chubaomonitors,scope=Namespaced
type ChubaoMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChubaoMonitorSpec   `json:"spec,omitempty"`
	Status ChubaoMonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ChubaoMonitorList contains a list of ChubaoMonitor
type ChubaoMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ChubaoMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ChubaoMonitor{}, &ChubaoMonitorList{})
}
