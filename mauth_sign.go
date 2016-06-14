package go_mauth_client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
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
	method := params[1].(string)
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
	return strings.Join([]string{method, urlpath, body, mauth_app.app_id, strconv.FormatInt(epoch, 10)}, "\n")
}

func SignString(mauth_app *MAuthApp, string_to_sign string) (s string, err error) {
	message := []byte(string_to_sign)

	encrypted, err := privateEncrypt(mauth_app, message)
	if err != nil {
		return "", err
	}
	// string needs to be base64 encoded, with the newlines removed
	signed := strings.Replace(base64.StdEncoding.EncodeToString(encrypted), "\n", "", -1)
	return signed, err
}

// privateEncrypt implements OpenSSL's RSA_private_encrypt function
// taken from: https://github.com/marpaia/chef-golang/api.go
func privateEncrypt(mauth_app *MAuthApp, data []byte) (enc []byte, err error) {
	k := (mauth_app.rsa_private_key.N.BitLen() + 7) / 8
	tLen := len(data)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	if tLen > k-11 {
		err = errors.New("Data too long")
		return
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], data)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(mauth_app.rsa_private_key.N) > 0 {
		err = nil
		return
	}
	var m *big.Int
	var ir *big.Int
	if mauth_app.rsa_private_key.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, mauth_app.rsa_private_key.D, mauth_app.rsa_private_key.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, mauth_app.rsa_private_key.Precomputed.Dp, mauth_app.rsa_private_key.Primes[0])
		m2 := new(big.Int).Exp(c, mauth_app.rsa_private_key.Precomputed.Dq, mauth_app.rsa_private_key.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, mauth_app.rsa_private_key.Primes[0])
		}
		m.Mul(m, mauth_app.rsa_private_key.Precomputed.Qinv)
		m.Mod(m, mauth_app.rsa_private_key.Primes[0])
		m.Mul(m, mauth_app.rsa_private_key.Primes[1])
		m.Add(m, m2)

		for i, values := range mauth_app.rsa_private_key.Precomputed.CRTValues {
			prime := mauth_app.rsa_private_key.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, mauth_app.rsa_private_key.N)
	}
	enc = m.Bytes()
	return
}
