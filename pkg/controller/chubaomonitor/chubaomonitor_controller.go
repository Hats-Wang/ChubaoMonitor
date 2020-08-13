package chubaomonitor

import (
	"context"
	"reflect"

	cachev1alpha1 "github.com/ChubaoMonitor/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_chubaomonitor")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ChubaoMonitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileChubaoMonitor{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("chubaomonitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ChubaoMonitor
	err = c.Watch(&source.Kind{Type: &cachev1alpha1.ChubaoMonitor{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.ChubaoMonitor{},
	})
	if err != nil {
		return err
	}
	log.Info("Deployment being watched")

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.ChubaoMonitor{},
	})
	if err != nil {
		return err
	}
	log.Info("Service being watched")
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1alpha1.ChubaoMonitor{},
	})
	if err != nil {
		return err
	}
	log.Info("ConfigMap being watched")

	return nil
}

// blank assignment to verify that ReconcileChubaoMonitor implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileChubaoMonitor{}

// ReconcileChubaoMonitor reconciles a ChubaoMonitor object
type ReconcileChubaoMonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ChubaoMonitor object and makes changes based on the state read
// and what is in the ChubaoMonitor.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileChubaoMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	log.Info("Reconciling ChubaoMonitor")
	ctx := context.Background()

	// Fetch the ChubaoMonitor instance
	chubaomonitor := &cachev1alpha1.ChubaoMonitor{}
	err := r.client.Get(ctx, request.NamespacedName, chubaomonitor)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("ChubaoMonitor resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to fetch ChubaoMonitor")
		return ctrl.Result{}, err
	}
	//fetch the ChubaoMonitor instance successfully

	log.Info("fetch ChubaoMonitor successfully")
	//create desiredDeploymentPrometheus

	//check if the configmap exit. If not,call user's attention to that the configmap configuration file is missing.
	configmapchubaomonitor := &corev1.ConfigMap{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "monitor-config", Namespace: chubaomonitor.Namespace}, configmapchubaomonitor)
	if err != nil && errors.IsNotFound(err) {
		chubaomonitor.Status.Configmapstatus = false
		log.Error(err, "Configmap monitor-config resource not found in namespace", chubaomonitor.Namespace, "Please try 'kubectl apply -f chubaofsmonitor_configmap.yaml -n CHUBAOMONITOR.NAMESPACE'(CHUBAOMONITOR.NAMESPACE is your real namespace, like kube-system)")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get chubaomonitor configmap")
		chubaomonitor.Status.Configmapstatus = false
		return ctrl.Result{}, err
	}
	//fetch chubaomonitor configmap successfully
	chubaomonitor.Status.Configmapstatus = true

	desiredDeploymentPrometheus := r.Deploymentforprometheus(chubaomonitor)
	if err := controllerutil.SetControllerReference(chubaomonitor, desiredDeploymentPrometheus, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	//check if the prometheus deployment exit. If not, create one
	deploymentPrometheus := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "prometheus", Namespace: chubaomonitor.Namespace}, deploymentPrometheus)
	if err != nil && errors.IsNotFound(err) {
		//create prometheus deployment
		log.Info("Creating a new prometheus Deployment", "Deployment.Namespace", desiredDeploymentPrometheus.Namespace, "Deployment.Name", desiredDeploymentPrometheus.Name)
		err = r.client.Create(ctx, desiredDeploymentPrometheus)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Error(err, "Failed to create new prometheus Deployment", "Deployment.Namespace", desiredDeploymentPrometheus.Namespace, "Deployment.Name", desiredDeploymentPrometheus.Name)
			return ctrl.Result{}, err
		}
		//create the deployment successfully.
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get promethues Deployment")
	}
	//fetch the deploymentprometheus successfully

	//check if the deploymentprometheus is right
	if check := CompareDeployment(desiredDeploymentPrometheus, deploymentPrometheus); check {
		deploymentPrometheus.Spec = desiredDeploymentPrometheus.Spec
		if err := controllerutil.SetControllerReference(chubaomonitor, deploymentPrometheus, r.scheme); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Updating deploymentprometheus")

		if err = r.client.Update(ctx, deploymentPrometheus); err != nil && !errors.IsConflict(err) {
			return ctrl.Result{}, err
		}
	}

	//Update chubaomonitor.Status.PrometheusReplicas, if needed.
	if chubaomonitor.Status.PrometheusReplicas != *deploymentPrometheus.Spec.Replicas {
		chubaomonitor.Status.PrometheusReplicas = *deploymentPrometheus.Spec.Replicas
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}
	}

	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(deploymentPrometheus.Namespace),
		client.MatchingLabels(labelsForChubaoMonitor(deploymentPrometheus.Name)),
	}
	if err = r.client.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "deploymentPrometheus.Namespace", deploymentPrometheus.Namespace, "deploymentPrometheus.Name", deploymentPrometheus.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update chubaomonitor.Status.PrometheusPods, if needed.
	if !reflect.DeepEqual(podNames, chubaomonitor.Status.PrometheusPods) {
		chubaomonitor.Status.PrometheusPods = podNames
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}
	}

	//create desiredServicePrometheus
	desiredServicePrometheus := Serviceforprometheus(chubaomonitor)
	if err := controllerutil.SetControllerReference(chubaomonitor, desiredServicePrometheus, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	//check if the prometheus service exit. If not, create one
	servicePrometheus := &corev1.Service{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "prometheus-service", Namespace: chubaomonitor.Namespace}, servicePrometheus)

	if err != nil && errors.IsNotFound(err) {
		//create the prometheus service
		log.Info("Creating a new promethues Service", "Service.Namespace", desiredServicePrometheus.Namespace, "Service.Name", "prometheus-service")
		err = r.client.Create(ctx, desiredServicePrometheus)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Error(err, "Failed to create new promethues Service", "Service.Namespace", desiredServicePrometheus.Namespace, "Service.Name", "prometheus-service")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get prometheus Service")
	}
	//fetch the serviceprometheus successfully

	//check if the service Prometheus is right
	if check := CompareService(servicePrometheus, desiredServicePrometheus); check {
		servicePrometheus.Spec.Ports = desiredServicePrometheus.Spec.Ports
		servicePrometheus.Spec.Type = corev1.ServiceTypeClusterIP
		servicePrometheus.Spec.Selector = desiredServicePrometheus.Spec.Selector
		if err := controllerutil.SetControllerReference(chubaomonitor, servicePrometheus, r.scheme); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Updating serviceprometheus")

		if err := r.client.Update(ctx, servicePrometheus); err != nil && !errors.IsConflict(err) {
			return ctrl.Result{}, err
		}
	}

	if chubaomonitor.Status.PrometheusclusterIP != servicePrometheus.Spec.ClusterIP {
		chubaomonitor.Status.PrometheusclusterIP = servicePrometheus.Spec.ClusterIP
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}
	}

	//create desiredDeploymentGrafana

	desiredDeploymentGrafana := r.Deploymentforgrafana(chubaomonitor)
	if err := controllerutil.SetControllerReference(chubaomonitor, desiredDeploymentGrafana, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	//check whether the grafana deployment exit. If not, create one
	deploymentGrafana := &appsv1.Deployment{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "grafana", Namespace: chubaomonitor.Namespace}, deploymentGrafana)
	if err != nil && errors.IsNotFound(err) {
		//create the grafana deployment
		log.Info("Creating a new grafana Deployment", "Deployment.Namespace", desiredDeploymentGrafana.Namespace, "Deployment.Name", desiredDeploymentGrafana.Name)
		err = r.client.Create(ctx, desiredDeploymentGrafana)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Error(err, "Failed to create new grafana Deployment", "Deployment.Namespace", desiredDeploymentGrafana.Namespace, "Deployment.Name", desiredDeploymentGrafana.Name)
			return ctrl.Result{}, err
		}
		//create the deployment successfully.

		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get grafana Deployment")
		return ctrl.Result{}, err
	}
	//fetch the deploymentgrafana successfully

	//check if the deploymentgrafana is right
	if check := CompareDeployment(desiredDeploymentGrafana, deploymentGrafana); check {
		deploymentGrafana.Spec.Replicas = desiredDeploymentGrafana.Spec.Replicas
		deploymentGrafana.Spec = desiredDeploymentGrafana.Spec
		//		deploymentGrafana.Spec.Selector = desiredDeploymentGrafana.Spec.Selector
		//		deploymentGrafana.Spec.Template = desiredDeploymentGrafana.Spec.Template
		log.Info("Updating deploymentgrafana")

		if err := controllerutil.SetControllerReference(chubaomonitor, deploymentGrafana, r.scheme); err != nil {
			return ctrl.Result{}, err
		}

		if err = r.client.Update(ctx, deploymentGrafana); err != nil && !errors.IsConflict(err) {
			return ctrl.Result{}, err
		}
	}

	//Update chubaomonitor.Status.GrafanaReplicas, if needed.
	if chubaomonitor.Status.GrafanaReplicas != *deploymentGrafana.Spec.Replicas {
		chubaomonitor.Status.GrafanaReplicas = *deploymentGrafana.Spec.Replicas
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}
	}

	podList = &corev1.PodList{}
	listOpts = []client.ListOption{
		client.InNamespace(deploymentGrafana.Namespace),
		client.MatchingLabels(labelsForChubaoMonitor(deploymentGrafana.Name)),
	}
	if err = r.client.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "deploymentGrafana.Namespace", deploymentGrafana.Namespace, "deploymentGrafana.Name", deploymentGrafana.Name)
		return ctrl.Result{}, err
	}
	podNames = getPodNames(podList.Items)

	// Update chubaomonitor.Status.GrafanaPods, if needed.
	if !reflect.DeepEqual(podNames, chubaomonitor.Status.GrafanaPods) {
		chubaomonitor.Status.GrafanaPods = podNames
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}
	}

	//create desiredServiceGrafana

	desiredServiceGrafana := Serviceforgrafana(chubaomonitor)
	if err := controllerutil.SetControllerReference(chubaomonitor, desiredServiceGrafana, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	//check if the grafana service exit. If not, create one
	serviceGrafana := &corev1.Service{}
	err = r.client.Get(ctx, types.NamespacedName{Name: "grafana-service", Namespace: chubaomonitor.Namespace}, serviceGrafana)

	if err != nil && errors.IsNotFound(err) {
		//create the grafana service
		log.Info("Creating a new grafana Service", "Service.Namespace", desiredServiceGrafana.Namespace, "Service.Name", "grafana-service")
		err = r.client.Create(ctx, desiredServiceGrafana)
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Error(err, "Failed to create new grafana Service", "Service.Namespace", desiredServiceGrafana.Namespace, "Service.Name", "grafana-service")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get grafana Service")
	}
	//fetch the servicegrafana successful

	//check if the grafana service is right
	if check := CompareService(serviceGrafana, desiredServiceGrafana); check {
		serviceGrafana.Spec.Ports = desiredServiceGrafana.Spec.Ports
		serviceGrafana.Spec.Type = corev1.ServiceTypeClusterIP
		serviceGrafana.Spec.Selector = desiredServiceGrafana.Spec.Selector
		if err := controllerutil.SetControllerReference(chubaomonitor, serviceGrafana, r.scheme); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Updating servicegrafana")

		if err := r.client.Update(ctx, serviceGrafana); err != nil && !errors.IsConflict(err) {
			return ctrl.Result{}, err
		}
	}

	if chubaomonitor.Status.GrafanaclusterIP != serviceGrafana.Spec.ClusterIP {
		chubaomonitor.Status.GrafanaclusterIP = serviceGrafana.Spec.ClusterIP
		err := r.client.Status().Update(ctx, chubaomonitor)
		if err != nil && !errors.IsConflict(err) {
			log.Error(err, "Failed to update ChubaoMonitor status")
			return ctrl.Result{}, err
		}

	}

	return reconcile.Result{}, nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
