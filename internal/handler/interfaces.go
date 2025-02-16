package handler

type Context interface {
	ShouldBindJSON(v interface{}) error
	MustGet(key string) interface{}
	JSON(code int, obj interface{})
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	Param(key string) string
}

type Router interface {
	Group(string) Router
	Use(Middleware)
	POST(path string, handler func(Context))
	GET(path string, handler func(Context))
}

type Middleware interface {
	Handle(handlerFunc func(Context)) func(Context)
}
