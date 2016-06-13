package go_mauth_client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func MakeAuthenticationHeaders(mauth_app *MAuthApp, signed_string string, seconds_since_epoch int64) map[string]string {
	headers := map[string]string{
		"X-MWS-Authentication": fmt.Sprintf("MWS %s:%s", mauth_app.app_id, signed_string),
		"X-MWS-Time":           string(seconds_since_epoch),
	}
	return headers
}

func MakeSignatureString(params ...interface{}) string {
	mauth_app := params[0].(*MAuthApp)
	verb := params[1].(string)
	rawurlstring := params[2].(string)
	body := params[3].(string)
	var epoch int64
	if len(params) == 5 {
		// if supplied, use it
		epoch = params[4].(int64)
	} else {
		// default if not supplied
		epoch = time.Now().Unix()
	}
	// remove the query strings
	urlpath := strings.Split(rawurlstring, "?")[0]
	return strings.Join([]string{verb, urlpath, body, mauth_app.app_id, strconv.FormatInt(epoch, 10)}, "\n")
}

func SignString(mauth_app *MAuthApp, string_to_sign string) (s []byte, err error) {
	rng := rand.Reader
	message := []byte(string_to_sign)

	hashed := sha512.Sum512(message)

	signature, err := rsa.SignPKCS1v15(rng, mauth_app.rsa_private_key, crypto.SHA512, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, err
}
