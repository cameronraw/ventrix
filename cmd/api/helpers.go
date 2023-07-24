package api

import (
	"github.com/cameronraw/ventrix/cmd/security"
	"github.com/cameronraw/ventrix/cmd/security/strategies"
)

func CreateSecurityMiddleware(strategy security.SecurityStrategy) security.SecurityMiddleware {
  return security.CreateSecurityMiddleware(strategy)
}

func (app *Application) ConfigureSecurityMiddleware() security.SecurityMiddleware {
  strategy := strategies.CreateApiKeySecurityStrategy(app.config.Key)
  return CreateSecurityMiddleware(strategy)
}
