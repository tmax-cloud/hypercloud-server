package v1

import (
	// v1 "github.com/tmax-cloud/efk-operator/api/v1"
	v1 "github.com/tmax-cloud/hypercloud-server/external/hyper/v1"

	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

type HyperV1Interface interface {
	RESTClient() rest.Interface
	HyperClusterResourceGetter
}

// HyperV1Client is used to interact with features provided by the  group.
type HyperV1Client struct {
	restClient rest.Interface
}

func (c *HyperV1Client) HyperClusterResources(namespace string) HyperClusterResourceInterface {
	return newHyperClusterResources(c, namespace)
}

// NewForConfig creates a new HyperV1Client for the given config.
func NewForConfig(c *rest.Config) (*HyperV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &HyperV1Client{client}, nil
}

func setConfigDefaults(config *rest.Config) error {
	v1.AddToScheme(scheme.Scheme)
	gv := v1.GroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *HyperV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
