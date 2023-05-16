package apis

import (
	"wadary/database/redis"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

type Router struct {
	Port  string
	Redis redis.IRedisCache
}

func NewRouter(port string, redis redis.IRedisCache) *Router {
	return &Router{
		Port:  port,
		Redis: redis,
	}
}

func (r *Router) Start() error {
	var err error
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.AllowAll())

	relayer := router.Group("/relayer")
	relayer.GET("/data", r.GetData)
	relayer.POST("/data", r.PostData)
	relayer.GET("/result", r.GetResult)
	relayer.POST("/result", r.PostResult)

	err = router.Run("0.0.0.0:" + r.Port)
	if err != nil {
		return err
	}
	return nil
}
