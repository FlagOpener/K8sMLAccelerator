
package controller

import (
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	coreV1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetAliases(app string, config []Config) []coreV1.HostAlias {
	for _, conf := range config {
		glog.V(5).Infof("looking for %s using %s", app, conf.Label)
		if conf.Label == app {
			return conf.Aliases
		}
	}
	return nil
}

func GetAliasesByKV(k, v string, config []Config) []coreV1.HostAlias {
	for _, conf := range config {
		glog.V(5).Infof("looking for %s, %s using %s", k, v, conf.Label)