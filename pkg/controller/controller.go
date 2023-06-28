
package controller

import (
	"os"
	"time"

	"github.com/golang/glog"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	IntializerConfigmapName string
	InitializerName         string
	IntializerNamespace     string
)

type Config struct {
	Name    string             `yaml:"name"`
	App     string             `yaml:"app"`
	Label   string             `yaml:"label"`
	Aliases []coreV1.HostAlias `yaml:"hostAliases"`
}

type Controller struct {
	clientset     *kubernetes.Clientset
	podController cache.Controller
	config        *[]Config
}

func NewHostAliasesInitializer(clientset *kubernetes.Clientset, conf *[]Config) *Controller {
	c := &Controller{
		config:    conf,
		clientset: clientset,
	}

	restClient := clientset.CoreV1().RESTClient()
	watchlist := cache.NewListWatchFromClient(restClient, "pods", coreV1.NamespaceAll, fields.Everything())

	// Wrap the returned watchlist to workaround the inability to include
	// the `IncludeUninitialized` list option when setting up watch clients.
	includeUninitializedWatchlist := &cache.ListWatch{
		ListFunc: func(options metaV1.ListOptions) (runtime.Object, error) {
			options.IncludeUninitialized = true
			return watchlist.List(options)
		},
		WatchFunc: func(options metaV1.ListOptions) (watch.Interface, error) {
			options.IncludeUninitialized = true
			return watchlist.Watch(options)
		},
	}

	resyncPeriod := 30 * time.Second

	_, podController := cache.NewInformer(
		includeUninitializedWatchlist,
		&coreV1.Pod{},