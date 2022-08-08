package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetServices(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	services, err := service.GetServices(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, services)
	}
}
