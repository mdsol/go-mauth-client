package go_mauth_client

import (
	"net/http"
	"time"
)

func MAuthSignerHandler(fn http.HandlerFunc, mauth_app MAuthApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seconds_since_epoch := time.Now().Unix()
		string_to_sign := MakeSignatureString(mauth_app, r.Method, r.URL.Path, r.Body, seconds_since_epoch)
		signed_string, err := SignString(mauth_app, string_to_sign)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		made_headers := MakeAuthenticationHeaders(mauth_app, signed_string, seconds_since_epoch)
		for header, value := range made_headers {
			w.Header().Set(header, value)
		}
		fn(w, r)
	}
}
