package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStatefulSets(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	sts, err := service.GetStatefulSets(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, sts)
	}
}
