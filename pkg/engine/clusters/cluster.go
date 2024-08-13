package clusters

import (
	"context"
	"sync"

	appsv1alpha1 "code.byted.org/inf/superkruise/api/apps/v1alpha1"
	skutil "code.byted.org/inf/superkruise/pkg/superkruiseutil"
	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/client/simple"
	engineclient "github.com/kyverno/chainsaw/pkg/engine/client"
	restutils "github.com/kyverno/chainsaw/pkg/utils/rest"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Cluster interface {
	Config() (*rest.Config, error)
	Build() (*rest.Config, client.Client, error)
}

type fromConfig struct {
	config *rest.Config
}

func NewClusterFromConfig(config *rest.Config) (Cluster, error) {
	return &fromConfig{
		config: config,
	}, nil
}

func (c *fromConfig) Config() (*rest.Config, error) {
	return c.config, nil
}
func (c *fromConfig) Build() (*rest.Config, client.Client, error) {
	config, err := c.Config()
	if err != nil {
		return nil, nil, err
	}
	client, err := simple.New(config)
	if err != nil {
		return nil, nil, err
	}
	client = engineclient.New(client)
	return config, client, nil
}

type fromKubeconfig struct {
	resolver func() (*rest.Config, error)
}

func NewClusterFromKubeconfig(kubeconfig string, context string) Cluster {
	resolver := sync.OnceValues(func() (*rest.Config, error) {
		return restutils.Config(kubeconfig, clientcmd.ConfigOverrides{
			CurrentContext: context,
		})
	})
	return &fromKubeconfig{
		resolver: resolver,
	}
}

func (c *fromKubeconfig) Config() (*rest.Config, error) {
	return c.resolver()
}
func (c *fromKubeconfig) Build() (*rest.Config, client.Client, error) {
	config, err := c.Config()
	if err != nil {
		return nil, nil, err
	}
	client, err := simple.New(config)
	if err != nil {
		return nil, nil, err
	}
	client = engineclient.New(client)
	return config, client, nil
}

type fromDataCluster struct {
	resolver func() (*rest.Config, error)
	client   func() (*rest.Config, client.Client, error)
}

func NewDataCluster(defaultCluster Cluster, clusterName string) Cluster {
	cfg, err := defaultCluster.Config()
	if err != nil {
		panic(err)
	}
	f := engineclient.InitDataClusterClient(cfg)
	resolver := func() (*rest.Config, error) {
		return f.GetRestConfigByClusterName(clusterName), nil
	}
	client := func() (*rest.Config, client.Client, error) {
		cc, err := f.CreateClient(clusterName, &cache.Options{})
		if err != nil {
			return nil, nil, err
		}
		return f.GetRestConfigByClusterName(clusterName), cc, nil
	}
	return &fromDataCluster{
		resolver: resolver,
		client:   client,
	}
}

func (c *fromDataCluster) Config() (*rest.Config, error) {
	return c.resolver()
}
func (c *fromDataCluster) Build() (*rest.Config, client.Client, error) {
	return c.client()
}

type fromVCluster struct {
	resolver func() (*rest.Config, error)
}

func NewVCluster(ctx context.Context, clusterName, clusterNamespace string) Cluster {
	resolver := sync.OnceValues(func() (*rest.Config, error) {
		return skutil.GetRestConfig(ctx, skutil.SecretNamePrefix+clusterName, clusterNamespace)
	})

	return &fromVCluster{
		resolver: resolver,
	}
}

func (c *fromVCluster) Config() (*rest.Config, error) {
	return c.resolver()
}
func (c *fromVCluster) Build() (*rest.Config, client.Client, error) {
	cfg, err := c.Config()
	if err != nil {
		return nil, nil, err
	}
	cli, err := GetVClusterClient(cfg)
	return cfg, cli, err
}

func GetVClusterClient(config *rest.Config) (client.Client, error) {
	return GetVClusterClientWithScheme(config, newDefaultScheme())
}

func GetVClusterClientWithScheme(config *rest.Config, scheme *runtime.Scheme) (client.Client, error) {
	return runtimeclient.New(config, runtimeclient.Options{
		Scheme: scheme,
		Cache: &runtimeclient.CacheOptions{
			DisableFor: []runtimeclient.Object{
				&corev1.ConfigMap{},
				&corev1.Secret{},
			},
		},
	})
}

func newDefaultScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))

	return scheme
}
