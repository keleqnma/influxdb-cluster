package router

import (
	"influxcluster/server"

	"github.com/gin-gonic/gin"
)

type Router struct {
	clusterSrv *server.ClusterServer
}

func NewRouter() *Router {
	return &Router{
		clusterSrv: server.NewClusterServer(),
	}
}

func (r *Router) PublicServer() *gin.Engine {
	gin.DisableConsoleColor()
	e := gin.New()
	apiV1 := e.Group("api/v1")
	{
		apiV1.GET("query", r.clusterSrv.Query)
		apiV1.POST("write", r.clusterSrv.Write)
	}
	return e
}

func (r *Router) Dispose() {
	//TODO
}
