package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetPods(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	pods, err := service.GetPods(namespaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, pods)
	}
}

type podDeleteDto struct {
	Names []string `json:"names"`
}

func DeletePods(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	podDeleteDto := &podDeleteDto{}
	if err := c.ShouldBind(&podDeleteDto); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err := service.DeletePods(namespaceName, podDeleteDto.Names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, "success")
	}
}

func ExecContainer(c *gin.Context) {
	namespaceName := c.Param("namespaceName")
	podName := c.Param("podName")
	containerName := c.Param("containerName")
	method := c.DefaultQuery("action", "sh")
	err := service.WebSSH(namespaceName, podName, containerName, method, c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
}
