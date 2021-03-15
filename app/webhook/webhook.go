
package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/glog"

	"github.com/fast-ml/nezha/pkg/controller"
	"k8s.io/api/admission/v1beta1"
	batch "k8s.io/api/batch/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	configFile    string
	useTLS        *bool
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
	hostAliasConf *[]controller.Config
	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter                  = runtime.ObjectDefaulter(runtimeScheme)
	addHostAliasesPatch string = `[{"op": "add", "path": "/spec/template/spec/hostAliases", "value": %s }]`
)
