package client

import (
	"code.byted.org/inf/superkruise/pkg/clientcache"
	"k8s.io/client-go/rest"
	"sync"
	"time"
)

var once sync.Once

func InitDataClusterClient(cfg *rest.Config) clientcache.ClientFactoryInterface {
	client := new(clientcache.ClientFactoryShop)
	f := client.Generate(clientcache.ClientFactoryNameDefault)
	once.Do(func() {
		go f.Start(cfg, 0)
		time.Sleep(5 * time.Second)
	})
	return f
}
