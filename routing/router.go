package routing

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/hoopra/GoAuthServer/authorization"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc negroni.HandlerFunc
	Secure      bool
}

// GetRouting returns the complete routing
// for this server
func GetRouting() *mux.Router {

	routers := newRouter(Routes)
	return routers
}

func newRouter(routes []route) *mux.Router {

	router := mux.NewRouter()
	for _, route := range routes {

		handler := negroni.New(negroni.HandlerFunc(route.HandlerFunc))
		if route.Secure {
			handler = negroni.New(
				negroni.HandlerFunc(authorization.RequireTokenAuthentication),
				negroni.HandlerFunc(route.HandlerFunc),
			)
		}

		router.Handle(route.Pattern, handler).Methods(route.Method).Name(route.Name)

	}

	return router
}
