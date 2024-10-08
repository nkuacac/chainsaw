package clusters

import (
	"github.com/kyverno/chainsaw/pkg/client"
	"k8s.io/client-go/rest"
)

const DefaultClient = ""

type Registry interface {
	Register(string, Cluster) Registry
	Lookup(string) Cluster
	Build(Cluster) (*rest.Config, client.Client, error)
}

type clientFactory = func(Cluster) (*rest.Config, client.Client, error)

func defaultClientFactory(cluster Cluster) (*rest.Config, client.Client, error) {
	if cluster == nil {
		return nil, nil, nil
	}
	return cluster.Build()
}

type registry struct {
	clientFactory clientFactory
	clusters      map[string]Cluster
}

func NewRegistry(f clientFactory) Registry {
	return registry{
		clientFactory: f,
		clusters:      map[string]Cluster{},
	}
}

func (c registry) Register(name string, cluster Cluster) Registry {
	values := map[string]Cluster{}
	for k, v := range c.clusters {
		values[k] = v
	}
	values[name] = cluster
	return registry{
		clientFactory: c.clientFactory,
		clusters:      values,
	}
}

func (c registry) Lookup(name string) Cluster {
	return c.clusters[name]
}

func (c registry) Build(cluster Cluster) (*rest.Config, client.Client, error) {
	f := c.clientFactory
	if f == nil {
		f = defaultClientFactory
	}
	return f(cluster)
}
