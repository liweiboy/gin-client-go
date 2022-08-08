package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfigMaps(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	maps, err := service.GetConfigMaps(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, maps)
	}
}
