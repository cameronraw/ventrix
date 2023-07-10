package api

import (
	"net/http"
)

func (app *Application) ConfigureMiddleware() {
	securityMiddleware := Middleware(app.ConfigureSecurityMiddleware())

	app.middleware = append(app.middleware, securityMiddleware)
}

type Middleware interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

type SecurityStrategy interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}
