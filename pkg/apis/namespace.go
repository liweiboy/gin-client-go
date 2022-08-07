package apis

import (
	"gin-client-go/gin-client-go/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNamespace(r *gin.Context) {
	namespaces, err := service.GetNamespaces()
	if err != nil {
		r.JSON(http.StatusInternalServerError, err.Error())
	} else {
		r.JSON(http.StatusOK, namespaces)
	}
}
