package engine

import (
	"context"
	"github.com/kyverno/chainsaw/pkg/mutate"
	"path/filepath"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/kyverno/chainsaw/pkg/engine/bindings"
	"github.com/kyverno/chainsaw/pkg/engine/clusters"
)

func WithBindings(ctx context.Context, tc Context, variables ...v1alpha1.Binding) (Context, error) {
	for _, variable := range variables {
		name, value, err := bindings.ResolveBinding(ctx, tc.Bindings(), nil, variable)
		if err != nil {
			return tc, err
		}
		tc = tc.WithBinding(ctx, name, value)
	}
	return tc, nil
}

func WithClusters(ctx context.Context, tc Context, basePath string, c map[string]v1alpha1.Cluster) Context {
	for name, cluster := range c {
		if len(cluster.VClusterName) > 0 && len(cluster.VClusterNameSpace) > 0 {
			vClusterName, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, cluster.VClusterName), nil, tc.Bindings())
			if err != nil {
				panic(err)
			}
			vClusterNameSpace, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, cluster.VClusterNameSpace), nil, tc.Bindings())
			if err != nil {
				panic(err)
			}
			ck := clusters.NewVCluster(ctx, vClusterName.(string), vClusterNameSpace.(string))
			tc = tc.WithCluster(ctx, name, ck)
		} else if len(cluster.DataClusterName) > 0 {
			get, err := mutate.Mutate(ctx, nil, mutate.Parse(ctx, cluster.DataClusterName), nil, tc.Bindings())
			if err != nil {
				panic(err)
			}
			controlClusterName := cluster.ControlClusterName
			ck := clusters.NewDataCluster(tc.Cluster(controlClusterName), get.(string))
			tc = tc.WithCluster(ctx, name, ck)
		} else {
			kubeconfig := filepath.Join(basePath, cluster.Kubeconfig)
			cluster := clusters.NewClusterFromKubeconfig(kubeconfig, cluster.Context)
			tc = tc.WithCluster(ctx, name, cluster)
		}
	}
	return tc
}

func WithCurrentCluster(ctx context.Context, tc Context, name string) (Context, error) {
	tc = tc.WithCurrentCluster(ctx, name)
	config, client, err := tc.CurrentClusterClient()
	if err != nil {
		return tc, err
	}
	tc = tc.WithBinding(ctx, "client", client)
	tc = tc.WithBinding(ctx, "config", config)
	return tc, nil
}

func WithNamespace(ctx context.Context, tc Context, namespace string) Context {
	return tc.WithBinding(ctx, "namespace", namespace)
}

func WithValues(ctx context.Context, tc Context, values any) Context {
	return tc.WithBinding(ctx, "values", values)
}
