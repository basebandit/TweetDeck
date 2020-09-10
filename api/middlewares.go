package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"ekraal.org/avatarlysis/business/data/auth"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/api/global"
)

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before")
		h.ServeHTTP(w, r) //call original
		fmt.Println("After")
	})
}

//Authenticate validates a JWT from the `Authorization` header.
func Authenticate(ctx context.Context, a *auth.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {

			ctx, span := global.Tracer("avatarlysis").Start(ctx, "api.middlewares.authenticate")
			defer span.End()

			//Parse the authorization header. Expected header is of the format `Bearer <token>`.
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header format: Bearer <token>")
				render.Render(w, r, ErrUnauthorized(err))
				return
			}
			//Start a span to measure just the time spent in ParseClaims.
			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				render.Render(w, r, ErrUnauthorized(err))
				return
			}

			//Add claims to the context so they can be retrieved later.
			ctx = context.WithValue(ctx, auth.Key, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

//HasRole validates that an authenticated user has at least one role from a
//specified list.This method constructs the actual function that is used.
func HasRole(ctx context.Context, roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx, span := global.Tracer("avatarlysis").Start(ctx, "api.middlewares.hasrole")
			defer span.End()

			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				err := errors.New("claims missing from context")
				render.Render(w, r, ErrForbiddenWithError(err))
				return
			}

			if !claims.HasRole(roles...) {
				render.Render(w, r, ErrForbidden)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
