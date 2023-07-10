package strategies

import "net/http"

type ApiKeySecurityStrategy struct {
	key string
}

func CreateApiKeySecurityStrategy(key string) ApiKeySecurityStrategy {
  return ApiKeySecurityStrategy {
    key,
  }
}

func (akss ApiKeySecurityStrategy) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerKey := r.Header.Get("X-API-Key")
		if headerKey != akss.key {
			http.Error(w, "Not Allowed", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
