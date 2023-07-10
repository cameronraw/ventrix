package security

func CreateSecurityMiddleware(strategy SecurityStrategy) SecurityMiddleware {
  return SecurityMiddleware{strategy}
}
