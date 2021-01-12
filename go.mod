module github.com/tmax-cloud/hypercloud-server

go 1.14

require (
	github.com/gorilla/mux v1.7.3
	github.com/tmax-cloud/claim-operator v0.0.0-20201224052434-e7800d54877c
	github.com/tmax-cloud/cluster-manager-operator v0.0.0-20210109105810-9bf04e3af331
	github.com/tmax-cloud/efk-operator v0.0.0-20201207030412-fd9c02a3e1c2
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	sigs.k8s.io/controller-runtime v0.6.2
)
