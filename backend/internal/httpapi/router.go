package httpapi

import (
	"net/http"
	"strings"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
}

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Handle(route Route) {
	if route.Handler == nil {
		panic("httpapi route handler is required")
	}
	if strings.TrimSpace(route.Method) == "" {
		panic("httpapi route method is required")
	}
	if strings.TrimSpace(route.Pattern) == "" {
		panic("httpapi route pattern is required")
	}

	r.mux.Handle(route.Pattern, methodHandler(route.Method, route.Handler))
}

func (r *Router) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	r.Handle(Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	})
}

func (r *Router) Register(routes ...Route) {
	for _, route := range routes {
		r.Handle(route)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func methodHandler(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			w.Header().Set("Allow", method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, req)
	})
}
