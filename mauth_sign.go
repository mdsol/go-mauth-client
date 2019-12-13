package go_mauth_client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
Wraps the functions around signing a request and generating the headers
*/

// MakeAuthenticationHeaders generates the formatted headers as a map for
// insertion into the request headers.
func MakeAuthenticationHeaders(mauthApp *MAuthApp, signed_string string, seconds_since_epoch int64) map[string]string {
	headers := map[string]string{
		"X-MWS-Authentication": fmt.Sprintf("MWS %s:%s", mauthApp.AppId, signed_string),
		"X-MWS-Time":           strconv.FormatInt(seconds_since_epoch, 10),
	}
	return headers
}

// MakeAuthenticationHeadersV2 generates the formatted headers as a map for
// insertion into the request headers.
func MakeAuthenticationHeadersV2(mauthApp *MAuthApp, signed_string string, seconds_since_epoch int64) map[string]string {
	headers := map[string]string{
		"MCC-Authentication": fmt.Sprintf("MWSV2 %s:%s;", mauthApp.AppId, signed_string),
		"MCC-Time":           strconv.FormatInt(seconds_since_epoch, 10),
	}
	return headers
}

// MakeSignatureString generates the string to be signed as part of the MWS header
func MakeSignatureString(mauthApp *MAuthApp, method string, url string, body string, epoch int64) string {
	if epoch == -1 {
		epoch = time.Now().Unix()
	}
	// remove the query strings
	return strings.Join([]string{method, strings.Split(url, "?")[0],
		body, mauthApp.AppId, strconv.FormatInt(epoch, 10)}, "\n")
}

// MakeSignatureStringV2 generates the string to be signed as part of the MWS header
func MakeSignatureStringV2(mauthApp *MAuthApp, method string, url string, body string, epoch int64) string {
	if epoch == -1 {
		epoch = time.Now().Unix()
	}

	//hash body and query string
	bodyHasher := sha512.New()
	bodyHasher.Write([]byte(body))
	hashedBody := hex.EncodeToString(bodyHasher.Sum(nil))

	urlParts := strings.Split(url, "?")

	encodedQueryParams := ""
	if len(urlParts) > 1 {
		encodedQueryParams = buildEncodedQueryParams(urlParts[1])
	}

	// remove the query strings
	return strings.Join([]string{method, urlParts[0],
		hashedBody, mauthApp.AppId, strconv.FormatInt(epoch, 10), encodedQueryParams},
		"\n")
}

func buildEncodedQueryParams(queryString string) string {
	encodedQueryStrings := []string{}
	var queryArray = strings.Split(queryString, "&")
	sort.Strings(queryArray)
	for _, x := range queryArray {
		keyValue := strings.Split(x, "=")
		var escapedKey = url.QueryEscape(keyValue[0])
		var escapedValue = url.QueryEscape(keyValue[1])
		encodedQueryStrings = append(encodedQueryStrings, escapedKey+"="+escapedValue)
	}
	return strings.Join(encodedQueryStrings, "&")
}

// SignString encrypts and encodes the string to sign
func SignString(mauthApp *MAuthApp, stringToSign string) (s string, err error) {
	// create a hasher
	h := sha512.New()
	h.Write([]byte(stringToSign))
	hashed := hex.EncodeToString(h.Sum(nil))

	// thanks to https://github.com/johnduhart for this
	encrypted, err := rsa.SignPKCS1v15(rand.Reader, mauthApp.RsaPrivateKey, 0, []byte(hashed))
	if err != nil {
		return "", err
	}
	// string needs to be base64 encoded, with the newlines removed
	signed := strings.Replace(base64.StdEncoding.EncodeToString(encrypted), "\n", "", -1)
	return signed, err
}

// SignStringV2 encrypts and encodes the string to sign
func SignStringV2(mauthApp *MAuthApp, stringToSign string) (s string, err error) {
	// create a hasher
	hashed := sha512.Sum512([]byte(stringToSign))

	// thanks to https://github.com/johnduhart for this
	encrypted, err := rsa.SignPKCS1v15(rand.Reader, mauthApp.RsaPrivateKey, crypto.SHA512, hashed[:])
	if err != nil {
		return "", err
	}
	// string needs to be base64 encoded, with the newlines removed
	signed := strings.Replace(base64.StdEncoding.EncodeToString(encrypted), "\n", "", -1)
	return signed, err
}
