package resources

import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Resource interface {
	v1.Object
}

type ResourceMeta struct {
	Name      string
	Namespace string
	ID        string
	Type      string
}
