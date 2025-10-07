package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type MiddlewareManager struct {
	middlewares []Middleware
}

func NewMiddlewareManager() *MiddlewareManager {
	return &MiddlewareManager{
		middlewares: []Middleware{},
	}
}

var global_middleware []Middleware

func (mngr *MiddlewareManager) Use(middleware ...Middleware) {
	global_middleware = append(global_middleware, middleware...)
}

func (mngr *MiddlewareManager) With(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, m := range middleware {
			next = m(next)
		}
		return next
	}
}

func (mngr *MiddlewareManager) WrappedMux(handler http.Handler) http.Handler {
	next := handler
	for _, m := range global_middleware {
		next = m(next)
	}
	return next
}
