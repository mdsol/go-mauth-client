package go_mauth_client

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// MakeAuthenticationHeaders generates the formatted headers as a map for
// insertion into the request headers.
func MakeAuthenticationHeaders(mauth_app *MAuthApp, signed_string string, seconds_since_epoch int64) map[string]string {
	headers := map[string]string{
		"X-MWS-Authentication": fmt.Sprintf("MWS %s:%s", mauth_app.App_ID, signed_string),
		"X-MWS-Time":           strconv.FormatInt(seconds_since_epoch, 10),
	}
	return headers
}

// MakeSignatureString generates the string to be signed as part of the MWS header
func MakeSignatureString(mauth_app *MAuthApp, method string, url string, body string, epoch int64) string {
	if epoch == -1 {
		epoch = time.Now().Unix()
	}
	// remove the query strings
	return strings.Join([]string{method, strings.Split(url, "?")[0],
		body, mauth_app.App_ID, strconv.FormatInt(epoch, 10)}, "\n")
}

// SignString encrypts and encodes the string to sign
func SignString(mauth_app *MAuthApp, string_to_sign string) (s string, err error) {
	// create a hasher
	h := sha512.New()
	h.Write([]byte(string_to_sign))
	hashed := hex.EncodeToString(h.Sum(nil))
	encrypted, err := privateEncrypt(mauth_app, []byte(hashed))
	if err != nil {
		return "", err
	}
	// string needs to be base64 encoded, with the newlines removed
	signed := strings.Replace(base64.StdEncoding.EncodeToString(encrypted), "\n", "", -1)
	return signed, err
}

// privateEncrypt implements OpenSSL's RSA_private_encrypt function
// taken from: https://github.com/marpaia/chef-golang/api.go (MIT License)
func privateEncrypt(mauth_app *MAuthApp, data []byte) (enc []byte, err error) {
	k := (mauth_app.RSA_Private_Key.N.BitLen() + 7) / 8
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
	if c.Cmp(mauth_app.RSA_Private_Key.N) > 0 {
		err = nil
		return
	}
	var m *big.Int
	var ir *big.Int
	if mauth_app.RSA_Private_Key.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, mauth_app.RSA_Private_Key.D, mauth_app.RSA_Private_Key.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, mauth_app.RSA_Private_Key.Precomputed.Dp, mauth_app.RSA_Private_Key.Primes[0])
		m2 := new(big.Int).Exp(c, mauth_app.RSA_Private_Key.Precomputed.Dq, mauth_app.RSA_Private_Key.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, mauth_app.RSA_Private_Key.Primes[0])
		}
		m.Mul(m, mauth_app.RSA_Private_Key.Precomputed.Qinv)
		m.Mod(m, mauth_app.RSA_Private_Key.Primes[0])
		m.Mul(m, mauth_app.RSA_Private_Key.Primes[1])
		m.Add(m, m2)

		for i, values := range mauth_app.RSA_Private_Key.Precomputed.CRTValues {
			prime := mauth_app.RSA_Private_Key.Primes[2+i]
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
		m.Mod(m, mauth_app.RSA_Private_Key.N)
	}
	enc = m.Bytes()
	return
}
