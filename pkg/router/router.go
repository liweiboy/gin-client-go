package router

import (
	"gin-client-go/gin-client-go/pkg/apis"
	"gin-client-go/gin-client-go/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	middleware.InitMiddleware(r)
	r.GET("/ping", apis.Ping)
	// namespace
	r.GET("/namespaces", apis.GetNamespace)
	// pod
	r.GET("/namespace/:namespaceName/pods", apis.GetPods)
	r.DELETE("/namespace/:namespaceName/pods", apis.DeletePods)
	r.GET("/namespace/:namespaceName/pod/:podName/container/:containerName", apis.ExecContainer)

	// service
	r.GET("/namespace/:namespaceName/services", apis.GetServices)

	// deployment
	r.GET("/namespace/:namespaceName/deployments", apis.GetDeployments)

	// statefulset
	r.GET("/namespace/:namespaceName/statefulsets", apis.GetStatefulSets)

	// configmap
	r.GET("/namespace/:namespaceName/configmaps", apis.GetConfigMaps)

	// secret
	r.GET("/namespace/:namespaceName/secrets", apis.GetSecrets)

}
