package service

import (
	"context"
	"gin-client-go/gin-client-go/pkg/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetServices(namespaceName string) ([]v1.Service, error) {
	clientset, err := client.GetClientset()
	if err != nil {
		return nil, err
	}
	list, err := clientset.CoreV1().Services(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
