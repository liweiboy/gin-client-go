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
	r.GET("/namespace/:namespaceName/pod/:podName/container/:containerName", apis.ExecContainer)
}
