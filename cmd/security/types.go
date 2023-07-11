package security

import "net/http"

type SecurityStrategy interface {
  Wrap(next http.HandlerFunc) http.HandlerFunc
}

type SecurityMiddleware struct {
  securityStrategy SecurityStrategy
}

func (sm SecurityMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
  return sm.securityStrategy.Wrap(next)
}

