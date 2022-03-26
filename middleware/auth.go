package middleware

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
)

func Login(conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" {
				state := uuid.NewString() + uuid.NewString()

				const CookieLife = 30

				cookie := &http.Cookie{
					Name:    "state",
					Value:   state,
					Expires: time.Now().Add(CookieLife * time.Second),
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

func Redirect(conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/google/redirect" {
				inner.ServeHTTP(w, r)
				return
			}

			// verify CSRF token
			if ok := verifyCSRFToken(w, r); !ok {
				return
			}

			// exchange grant code for token
			tok, err := conf.Exchange(context.Background(), r.URL.Query().Get("code"))
			if err != nil {
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
		})
	}
}

func ValidateToken(conf *oauth2.Config) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/register" {
				inner.ServeHTTP(w, r)
				return
			}

			// get ID token from cookies
			tokenAuth, err := r.Cookie("auth")
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"User Not Authenticated"}}`)
				return
			}

			// validate ID token
			ok := validateIDToken(tokenAuth.Value, conf.ClientID, w)
			if !ok {
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

			req, _ := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/revoke?token="+token.Value, http.NoBody)

			res, err := client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, errors.InternalServerErr{})
				return
			}

			defer res.Body.Close()

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

func verifyCSRFToken(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("state")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if cookie.Value != r.URL.Query().Get("state") {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func validateIDToken(token, clientID string, w http.ResponseWriter) bool {
	payload, err := idtoken.Validate(context.Background(), token, clientID)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error":{"code":"invalid_request","reason":"user not authorized"}}`)

		return false
	}

	if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
		return false
	}

	return true
}
