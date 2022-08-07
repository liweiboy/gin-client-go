package main

import (
	"fmt"
	"gin-client-go/gin-client-go/pkg/config"
	"gin-client-go/gin-client-go/pkg/router"
	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func main() {
	engine := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.InitRouter(engine)

	err := engine.Run(fmt.Sprintf("%s:%s", config.AppConfig.Host, config.AppConfig.Port))
	if err != nil {
		klog.Fatal(err)
	}

}
