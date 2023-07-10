package api

import (
	"montecristo/cmd/api/security"
	"montecristo/cmd/api/security/strategies"
)

func CreateSecurityMiddleware(strategy SecurityStrategy) security.SecurityMiddleware {
  return security.CreateSecurityMiddleware(strategy)
}

func (app *Application) ConfigureSecurityMiddleware() security.SecurityMiddleware {
  strategy := strategies.CreateApiKeySecurityStrategy(app.config.Key)
  return CreateSecurityMiddleware(strategy)
}
