package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func extractTokenMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		auth := r.Header.Get("authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			ctx = context.WithValue(ctx, tokenKey, auth[7:])
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractCompanyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		companyId := r.Header.Get("company")
		ctx := r.Context()
		if companyId != "" {
			ctx = context.WithValue(ctx, companyIdKey, companyId)
		}
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func stripApiPrefix(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.HasPrefix(path, "/api") {
			http.Error(w, `{"error": "not found"}`, http.StatusForbidden)
			return
		}
		path = strings.TrimPrefix(path, "/api")
		if path == "" {
			path = "/"
		}
		r.URL.Path = path
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization", "Company", "Origin"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "POST"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

// allowCORS allows Cross Origin Resource Sharing from certain origins.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3002" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func handleStatusChecks(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			fmt.Fprintf(w, `{"ok": "true"}`)
			return
		}
		h.ServeHTTP(w, r)
	})
}
