package middleware

import (
	"html"
	"net/http"
	"strings"
)

// SanitizeInput is a simple middleware that escapes HTML in form values to prevent basic XSS.
func SanitizeInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if err := r.ParseForm(); err == nil {
				for key, values := range r.PostForm {
					for i, value := range values {
						r.PostForm[key][i] = html.EscapeString(strings.TrimSpace(value))
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
