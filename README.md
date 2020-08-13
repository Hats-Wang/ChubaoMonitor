# ChubaoMonitorOperator

Prerequisites:
```
go version v1.13+.
docker version 17.03+.
kubectl version v1.11.3+.
kustomize v3.1.0+
operator-SDK v0.12.0
Access to a Kubernetes v1.11.3+ cluster
```

Build and run ChubaoMonitorOperator in master node locally outside the cluster:
```
git clone https://github.com/Hats-Wang/ChubaoMonitor.git
cd ChubaoMonitor
kubectl create -f ./deploy/crds/cache.example.com_chubaomonitors_crd.yaml
operator-sdk up local --namespace=default
```


Edit file ./deploy/crds/cache.example.com_v1alpha1_chubaomonitor_cr.yaml to customize your own ChubaoMonitor instance.


Create your own ChubaoMonitor instance:
```
kubectl create -f ./deploy/crds/cache.example.com_v1alpha1_chubaomonitor_cr.yaml
```

CRD definition: ./deploy/crds/cache.example.com_chubaomonitors_crd.yaml.

Business logic code are located in folder ./pkg/controller/chubaomonitor.

Reference: https://docs.openshift.com/container-platform/4.3/operators/operator_sdk/osdk-getting-started.html
