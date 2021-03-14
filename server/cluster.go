package server

import (
	"errors"
	"influxcluster/conf"
	"influxcluster/logging"
	"influxcluster/service"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoDataBase = errors.New("database parameters missing")
	ErrNoQuery    = errors.New("q parameters missing")
	logger        = logging.GetLogger("server")
)

type ClusterServer struct {
	*service.ClusterService
}

func NewClusterServer(appCfg conf.APPConfig) *ClusterServer {
	return &ClusterServer{
		service.NewInfluxCluster(appCfg),
	}
}

func (cs *ClusterServer) Query(c *gin.Context) {
	db, ok := c.GetQuery("db")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, ErrNoDataBase)
		return
	}
	q, ok := c.GetQuery("q")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, ErrNoQuery)
		return
	}
	res, err := cs.ClusterService.Query(db, q)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (cs *ClusterServer) Write(c *gin.Context) {
	db, ok := c.GetQuery("db")
	if !ok {
		c.AbortWithError(http.StatusBadRequest, ErrNoDataBase)
		return
	}
	lp, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = cs.ClusterService.Write(db, lp)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
