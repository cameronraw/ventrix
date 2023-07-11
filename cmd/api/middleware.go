package api

import (
	"net/http"
)

type Middleware interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

func ConfigureMiddleware(app *Application) {
	middlewareToApply := []Middleware{
		Middleware(app.ConfigureSecurityMiddleware()),
	}

	app.Middleware = append(app.Middleware, middlewareToApply...)
}

func (app *Application) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	if len(app.Middleware) == 0 {
		return handler
	}

	whatToWrap := app.Middleware[len(app.Middleware)-1].Wrap(handler)

	if len(app.Middleware) == 1 {
		return whatToWrap
	}

	for i := len(app.Middleware) - 2; i >= 0; i-- {
		whatToWrap = app.Middleware[i].Wrap(whatToWrap)
	}

	return whatToWrap
}
