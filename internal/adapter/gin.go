package adapter

import (
	"github.com/gin-gonic/gin"
	
	"avito-shop-test/internal/handler"
)


type GinContext struct {
	c *gin.Context
}

func (g *GinContext) ShouldBindJSON(v interface{}) error {
	return g.c.ShouldBindJSON(v)
}

func (g *GinContext) MustGet(key string) interface{} {
	return g.c.MustGet(key)
}

func (g *GinContext) JSON(code int, obj interface{}) {
	g.c.JSON(code, obj)
}

func (g *GinContext) Set(key string, value interface{}) {
	g.c.Set(key, value)
}

func (g *GinContext) Get(key string) (value interface{}, exists bool) {
	return g.c.Get(key) 
}

func (g *GinContext) Param(key string) string {
	return g.c.Param(key) 
}


type GinRouter struct {
	group *gin.RouterGroup
}

func NewGinRouter(group *gin.RouterGroup) *GinRouter {
	return &GinRouter{group: group}
}

func (g *GinRouter) Group(relativePath string) handler.Router {
	return NewGinRouter(g.group.Group(relativePath))
}

func (g *GinRouter) Use(middleware handler.Middleware) {
	g.group.Use(func(c *gin.Context) {
		ctx := &GinContext{c: c}
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			ctx.Set("Authorization", authHeader)
		}

		
		middleware.Handle(func(hCtx handler.Context) {
			
		})(ctx) 
	})
}

func (g *GinRouter) POST(path string, handler func(handler.Context)) {
	g.group.POST(path, func(c *gin.Context) {
		ctx := &GinContext{c: c}
		handler(ctx)
	})
}

func (g *GinRouter) GET(path string, handler func(handler.Context)) {
	g.group.GET(path, func(c *gin.Context) {
		ctx := &GinContext{c: c}
		handler(ctx)
	})
}
