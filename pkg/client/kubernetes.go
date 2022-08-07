package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"sync"
)

var onceClient = sync.Once{}
var onceConfig = sync.Once{}
var clientset *kubernetes.Clientset
var kubeconfig *rest.Config

func GetClientset() (*kubernetes.Clientset, error) {
	onceClient.Do(func() {
		kubeconfig, err := GetKubeconfig()
		if err != nil {
			klog.Errorln(err)
			return
		}
		clientset, err = kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			klog.Errorln(err)
			return
		}
	})
	return clientset, nil
}

func GetKubeconfig() (*rest.Config, error) {
	onceConfig.Do(func() {
		var err error
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			kubeconfig, err = rest.InClusterConfig()
		}
	})
	return kubeconfig, nil
}
