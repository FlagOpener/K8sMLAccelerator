
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

// Config contains the server (the webhook) cert and key.
type certConfig struct {
	CertFile string
	KeyFile  string
}

func (c *certConfig) addFlags() {
	flag.StringVar(&configFile, "config-file", "", "path to hostAliases configuration config file")
	flag.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func configTLS(config certConfig) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		glog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		// TODO: uses mutual tls after we agree on what cert the apiserver should use.
		// ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

func mutateDeployments(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	glog.V(2).Info("mutating deployments")
	dpResource := metav1.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "deployments"}
	if ar.Request.Resource != dpResource {
		glog.Errorf("expect resource to be %s", dpResource)
		return nil
	}
