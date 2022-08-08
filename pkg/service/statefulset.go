package service

import (
	"context"
	"gin-client-go/gin-client-go/pkg/client"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetStatefulSets(namespaceName string) ([]v1.StatefulSet, error) {
	clientset, err := client.GetClientset()
	if err != nil {
		return nil, err
	}
	list, err := clientset.AppsV1().StatefulSets(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
