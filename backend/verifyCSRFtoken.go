package main

import (
	"errors"
	"fmt"
	"net/http"
)

func VerifyCSRFToken(r *http.Request) error {
	// Get the CSRF token from the form
	csrfToken := r.FormValue("csrf_token")
	if csrfToken == "" {
		return errors.New("missing CSRF token in the form")
	}

	// Get the CSRF token from the cookies
	csrfCookie, err := r.Cookie("csrf_token")
	if err != nil {
		return errors.New("missing CSRF token in the cookies")
	}

	// Compare the tokens
	if csrfToken != csrfCookie.Value {
		return errors.New("CSRF token mismatch")
	}

	if csrfToken == csrfCookie.Value {
		fmt.Println("CSRF token matched")
	}

	return nil
}
