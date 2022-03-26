package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"net/http"
)

func Login(conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				state := uuid.NewString() + uuid.NewString()

				cookie := &http.Cookie{
					Name:   "state",
					Value:  state,
					MaxAge: 100,
				}

				http.SetCookie(w, cookie)

				url := conf.AuthCodeURL(state)

				http.Redirect(w, r, url, http.StatusSeeOther)

				return
			}
			inner.ServeHTTP(w, r)
		})
	}
}

func Redirect(logger log.Logger, conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/google/redirect" {
				// verify CSRF token
				cookie, err := r.Cookie("state")
				if err != nil {
					logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if state := cookie.Value; state != r.URL.Query().Get("state") {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				// exchange grant code for token
				tok, err := conf.Exchange(context.Background(), r.URL.Query().Get("code"))
				if err != nil {
					logger.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"missing param code"}}`)
					return
				}

				idToken := tok.Extra("id_token").(string)

				// set authorization cookie
				idCookie := &http.Cookie{
					Name:    "auth",
					Value:   idToken,
					Expires: tok.Expiry,
					Path:    "/",
				}

				accessCookie := &http.Cookie{
					Name:    "access",
					Value:   tok.AccessToken,
					Expires: tok.Expiry,
					Path:    "/",
				}

				http.SetCookie(w, idCookie)
				http.SetCookie(w, accessCookie)
				http.Redirect(w, r, "http://localhost:8000/register", http.StatusFound)

				return
			}

			inner.ServeHTTP(w, r)
		})
	}
}

func ValidateToken(logger log.Logger, conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.URL.Path == "/register" {
				inner.ServeHTTP(w, r)
				return
			}

			tokenAuth, err := r.Cookie("auth")
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"User Not Authenticated"}}`)
				return
			}

			payload, err := idtoken.Validate(context.Background(), tokenAuth.Value, conf.ClientID)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"user not authorized"}}`)
				return
			}

			if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"user not authorized"}}`)
				return
			}

			inner.ServeHTTP(w, r)
		})
	}
}

func Logout(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logout" {
			client := http.DefaultClient

			token, _ := r.Cookie("access")

			req, _ := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/revoke?token="+token.Value, nil)

			res, _ := client.Do(req)
			if res.StatusCode != http.StatusOK {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, errors.InternalServerErr{})
				return
			}

			idCookie := &http.Cookie{
				Name:   "auth",
				MaxAge: -1,
			}

			accessCookie := &http.Cookie{
				Name:   "access",
				MaxAge: -1,
			}

			http.SetCookie(w, idCookie)
			http.SetCookie(w, accessCookie)

			w.WriteHeader(http.StatusNoContent)
			return
		}

		inner.ServeHTTP(w, r)
	})
}
