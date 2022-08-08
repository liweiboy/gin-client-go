package service

import (
	"context"
	"gin-client-go/gin-client-go/pkg/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMaps(namespaceName string) ([]v1.ConfigMap, error) {
	clientset, err := client.GetClientset()
	if err != nil {
		return nil, err
	}
	list, err := clientset.CoreV1().ConfigMaps(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
