package http

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/thomasmortensson/pxe-init/internal/domain/repositories"
)

const (
	paramServerContext = "serverContext"
)

type ServerContext struct {
	Logger        *zap.Logger
	Datastore     repositories.Datastore
	ForwardServer *url.URL
}

type Router struct {
	Engine *gin.Engine
}

func NewServerContext(logger *zap.Logger, db repositories.Datastore, forwardServer *url.URL) *ServerContext {
	return &ServerContext{
		Logger:        logger,
		Datastore:     db,
		ForwardServer: forwardServer,
	}
}

// NewRouter generates a new gin router and populates a logger and datastore connection on the default context
func NewRouter(serverCtx *ServerContext) *Router {
	gin.SetMode(gin.ReleaseMode)

	router := &Router{
		Engine: gin.Default(),
	}

	router.Engine.Use(func(c *gin.Context) {
		c.Set(paramServerContext, serverCtx)
		c.Next()
	})

	return router
}

// AddRoutes adds specified HTTP routes to the provided router struct. The following routes are added to the router:
func (r *Router) AddRoutes() {
	r.Engine.GET("/boot.ipxe", bootHandler)
	r.Engine.GET("/ipxe", ipxeHandler)
	// TODO add proper liveness and readiness endpoints
	r.Engine.GET("/ready", statusHandler)
	r.Engine.GET("/health", statusHandler)
}
