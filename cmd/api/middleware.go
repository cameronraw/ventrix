package api

import (
	"net/http"
)

type Middleware interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

type SecurityStrategy interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

func (app *Application) ConfigureMiddleware() {
  middlewareToApply := []Middleware{
    Middleware(app.ConfigureSecurityMiddleware()),
  }

  app.middleware = append(app.middleware, middlewareToApply...)
}

func (app *Application) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	if len(app.middleware) == 0 {
		return handler
	}

	whatToWrap := app.middleware[len(app.middleware)-1].Wrap(handler)

	if len(app.middleware) == 1 {
		return whatToWrap
	}

	for i := len(app.middleware) - 2; i >= 0; i-- {
		whatToWrap = app.middleware[i].Wrap(whatToWrap)
	}

	return whatToWrap
}

