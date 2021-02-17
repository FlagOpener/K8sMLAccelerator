package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"

	"github.com/fast-ml/nezha/pkg/controller"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultInitializerName    = "hostaliases.initializer.kubernetes.io"
	defaultConfigmapName      = "hostaliases-initializer"
	defaultConfigMapNamespace = "default"
)

var (
	kubeConfig string
	kubeMaster string
)

func main() {
	flag.StringVar(&controller.IntializerConfigmapName, "configmap", defaultConfigmapName, "initializer configuration configmap")
	flag.StringVar(&controller.InitializerName, "initializer-name", defaultInitializerName, "The initializer name")
	flag.StringVar(&controller.IntializerNamespace, "namespace", defaultConfigMapNamespace, "The configuration namespace")
	flag.StringVar(&kubeConfig, "kubeconfig", "", "Absolute path to the kubeconfig")
	flag.StringVar(&kubeMaster, "ku