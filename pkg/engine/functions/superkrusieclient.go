package functions

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"code.byted.org/inf/superkruise/pkg/clientcache"
	"github.com/jmespath/go-jmespath"
	"github.com/kyverno/chainsaw/pkg/engine/client"
	"github.com/kyverno/chainsaw/pkg/engine/functions/tracectx"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func createClusterClient(cluster string, f clientcache.ClientFactoryInterface) {
	//fmt.Println("before CreateClient", cluster, time.Now())
	for {
		_, err := f.CreateClient(cluster, &cache.Options{})
		if err != nil {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	//fmt.Println("after CreateClient", cluster, time.Now())
}

func jpAllDataClusterInformerInit(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, apiVersion, kind string
	var namespace string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &kind); err != nil {
		return nil, err
	}
	if len(arguments) >= 5 {
		if err := getArg(arguments, 4, &namespace); err != nil {
			return nil, err
		}
	}
	fmt.Println("jpAllDataClusterInformerInit args:", clusters, apiVersion, kind, namespace)

	f := client.InitDataClusterClient(cfg)
	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterInformerInit"), len(all), len(all), func(i int) error {
		cluster := all[i]
		createClusterClient(cluster, f)
		return nil
	})

	return clusters, err
}

func jpAllDataClusterInformerCleanup(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, apiVersion, kind string
	var namespace string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &kind); err != nil {
		return nil, err
	}
	if len(arguments) >= 5 {
		if err := getArg(arguments, 4, &namespace); err != nil {
			return nil, err
		}
	}
	fmt.Println("jpAllDataClusterInformerCleanup args:", clusters, apiVersion, kind, namespace)

	f := client.InitDataClusterClient(cfg)
	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterInformerCleanup"), len(all), len(all), func(i int) error {
		cluster := all[i]
		f.DestoryClient(cluster)
		return nil
	})

	return clusters, err
}

func jpAllDataClusterList(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, apiVersion, kind string
	var namespace string
	var label string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &kind); err != nil {
		return nil, err
	}
	if len(arguments) >= 5 {
		if err := getArg(arguments, 4, &namespace); err != nil {
			return nil, err
		}
	}
	if len(arguments) >= 6 {
		if err := getArg(arguments, 5, &label); err != nil {
			return nil, err
		}
	}
	fmt.Println("jpAllDataClusterList Args:", clusters, apiVersion, kind, label, namespace)
	var list unstructured.UnstructuredList
	list.SetAPIVersion(apiVersion)
	list.SetKind(kind)

	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterList Args:", clusters, apiVersion, kind, label, namespace),
		len(all), len(all), func(i int) error {
			cluster := all[i]
			client, err := getDataClusterClient(cluster, arguments)
			if err != nil {
				return err
			}

			dataList, err := dataK8sList(client, arguments[1:])
			if err == nil {
				//fmt.Println("jpAllDataClusterList get", Prettify(dataList))
				list.Items = append(list.Items, dataList.Items...)
				return nil
			}
			if apierrors.IsNotFound(err) {
				fmt.Println("jpAllDataClusterList not found", "not found")
				return nil
			}
			return err
		})
	//fmt.Println("jpAllDataClusterList all", Prettify(list))
	return list.UnstructuredContent(), err
}

func jpAllDataClusterPatch(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, apiVersion, kind string
	var namespace string
	var label string
	var newspec string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &kind); err != nil {
		return nil, err
	}
	if len(arguments) >= 5 {
		if err := getArg(arguments, 4, &namespace); err != nil {
			return nil, err
		}
	}
	if len(arguments) >= 6 {
		if err := getArg(arguments, 5, &label); err != nil {
			return nil, err
		}
	}
	if len(arguments) >= 7 {
		if err := getArg(arguments, 6, &newspec); err != nil {
			return nil, err
		}
	}
	fmt.Println("jpAllDataClusterPatch Args:", clusters, apiVersion, kind, label, namespace)

	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterPatch Args:", clusters, apiVersion, kind, label, namespace),
		len(all), len(all), func(i int) error {
			cluster := all[i]
			client, err := getDataClusterClient(cluster, arguments)
			if err != nil {
				return err
			}

			dataList, err := dataK8sList(client, arguments[1:])
			if err == nil {
				err = ParallelRun(ctx.Do().Step("dataClusterPatch:"),
					len(dataList.Items), len(dataList.Items), func(i int) error {
						item := dataList.Items[i]
						return client.Patch(context.Background(), &item, ctrlclient.Merge)
					})
				return err
			}
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		})
	return nil, err
}

func jpAllDataClusterWait(arguments []any) (any, error) {
	var path, expect string
	if len(arguments) >= 7 {
		if err := getArg(arguments, 6, &path); err != nil {
			return nil, err
		}
	}
	if len(arguments) >= 8 {
		if err := getArg(arguments, 7, &expect); err != nil {
			return nil, err
		}
	}

	path = strings.ReplaceAll(path, "`", "'")
	fmt.Println("jpAllDataClusterWait", path, expect)
	err := wait.PollUntilContextTimeout(context.TODO(), 1*time.Second, 60*time.Second, true, func(ctx context.Context) (done bool, err error) {
		list, err := jpAllDataClusterList(arguments)
		if err != nil {
			fmt.Println("jpAllDataClusterWait list err:", err)
			return false, err
		}
		search, err := jmespath.Search(path, list)
		if err != nil {
			fmt.Println("jpAllDataClusterWait search err:", err)
			return false, err
		}
		fmt.Println("jpAllDataClusterWait search result:", search, "expect:", expect)

		return search == expect, nil
	})

	return err == nil, err
}

func jpGetDataClusterClient(arguments []any) (any, error) {
	var cluster string
	if err := getArg(arguments, 1, &cluster); err != nil {
		return nil, err
	}
	return getDataClusterClient(cluster, arguments)
}

func getDataClusterClient(cluster string, arguments []any) (clientcache.ClientInterface, error) {
	var cfg *rest.Config
	var apiVersion, kind string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	//if err := getArg(arguments, 1, &cluster); err != nil {
	//	return nil, err
	//}
	if err := getArg(arguments, 2, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &kind); err != nil {
		return nil, err
	}
	fmt.Println("getDataClusterClient args:", cluster, apiVersion, kind)

	f := client.InitDataClusterClient(cfg)
	var client clientcache.ClientInterface
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do(), 1, 1, func(i int) error {
		var dataList unstructured.UnstructuredList
		dataList.SetAPIVersion(apiVersion)
		dataList.SetKind(kind)
		client = f.GetClient(cluster)

		if client == nil {
			return fmt.Errorf("%s %s %s client not found", cluster, apiVersion, kind)
		}
		return nil
	})

	return client, err
}

func jpAllDataClusterServerVersion(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("rest config is nil")
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}

	var versions []string
	f := client.InitDataClusterClient(cfg)
	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterServerVersion"), len(all), len(all), func(i int) error {
		cluster := all[i]
		//config := f.GetRestConfigByClusterName(cluster)
		//if config == nil {
		//	return fmt.Errorf("jpAllDataClusterServerVersion get GetRestConfigByClusterName of %s return nil", cluster)
		//}
		//client, err := discovery.NewDiscoveryClientForConfig(config)
		//if err != nil {
		//	return err
		//}

		client, err := getDataClientSet(f, cluster)
		if err != nil {
			return err
		}
		version, err := client.ServerVersion()
		if err != nil {
			fmt.Println("jpAllDataClusterServerVersion", cluster, err.Error())
			return err
		}

		fmt.Println("jpAllDataClusterServerVersion", cluster, Prettify(version))
		versions = append(versions, cluster, Prettify(version))
		return nil
	})

	//fmt.Println("jpAllDataClusterServerVersion", Prettify(versions))
	return Prettify(versions), err
}

func jpAllDataClusterCreateNamespace(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, namespace string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("rest config is nil")
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}

	if err := getArg(arguments, 2, &namespace); err != nil {
		return nil, err
	}

	f := client.InitDataClusterClient(cfg)
	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterCreateNamespace args:", clusters, namespace),
		len(all), len(all), func(i int) error {
			cluster := all[i]
			client, err := getDataClientSet(f, cluster)
			if err != nil {
				return err
			}
			ns := &v1.Namespace{}
			ns.Name = namespace
			_, err = client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})

			if err != nil {
				return fmt.Errorf("cluster: %s, err: %w", cluster, err)
			}
			return err
		})

	return clusters, err
}

func getDataClientSet(f clientcache.ClientFactoryInterface, cluster string) (*clientset.Clientset, error) {
	config := f.GetRestConfigByClusterName(cluster)
	if config == nil {
		return nil, fmt.Errorf("GetRestConfigByClusterName of %s return nil", cluster)
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func jpAllDataClusterDeleteNamespace(arguments []any) (any, error) {
	var cfg *rest.Config
	var clusters, namespace string
	if err := getArg(arguments, 0, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("rest config is nil")
	}
	if err := getArg(arguments, 1, &clusters); err != nil {
		return nil, err
	}

	if err := getArg(arguments, 2, &namespace); err != nil {
		return nil, err
	}

	f := client.InitDataClusterClient(cfg)
	all := strings.Split(clusters, ",")
	ctx := tracectx.Context{}
	err := ParallelRun(ctx.Do().Step("jpAllDataClusterDeleteNamespace args:", clusters, namespace),
		len(all), len(all), func(i int) error {
			cluster := all[i]
			client, err := getDataClientSet(f, cluster)
			if err != nil {
				return err
			}

			err = client.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("delete ns in cluster: %s, err: %w", cluster, err)
			}
			return nil
		})

	return nil, err
}

func jpDataKubernetesGet(arguments []any) (any, error) {
	var client clientcache.ClientInterface
	var apiVersion, kind string
	var key ctrlclient.ObjectKey
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &key.Namespace); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 4, &key.Name); err != nil {
		return nil, err
	}
	var obj unstructured.Unstructured
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)
	if err := client.Get(context.Background(), key, &obj); err != nil {
		return nil, err
	}
	return obj.UnstructuredContent(), nil
}

func jpDataKubernetesList(arguments []any) (any, error) {
	var client clientcache.ClientInterface
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	list, err := dataK8sList(client, arguments)
	if err == nil {
		return list.UnstructuredContent(), err
	} else {
		return nil, err
	}
}

func dataK8sList(client clientcache.ClientInterface, arguments []any) (*unstructured.UnstructuredList, error) {
	var apiVersion, kind, label, namespace string
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if len(arguments) >= 4 {
		if err := getArg(arguments, 3, &namespace); err != nil {
			return nil, err
		}
	}
	if len(arguments) == 5 {
		if err := getArg(arguments, 4, &label); err != nil {
			return nil, err
		}
	}
	var list unstructured.UnstructuredList
	list.SetAPIVersion(apiVersion)
	list.SetKind(kind)
	var listOptions []ctrlclient.ListOption
	if namespace != "" {
		listOptions = append(listOptions, ctrlclient.InNamespace(namespace))
	}
	selector := NewLabelSelector(strings.Split(label, ",")...)
	listOptions = append(listOptions, ctrlclient.MatchingLabelsSelector{Selector: selector})

	err := client.List(context.Background(), &list, listOptions...)
	if err == nil {
		//fmt.Println("jpDataKubernetesList", Prettify(list))
		return &list, nil
	}
	return nil, err
}
