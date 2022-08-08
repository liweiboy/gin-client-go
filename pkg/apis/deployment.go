package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDeployments(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	deployments, err := service.GetDeployments(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, deployments)
	}
}
