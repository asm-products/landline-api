package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	secret := os.Args[1]

	fmt.Println("starting with secret", secret)
	r := gin.Default()
	r.GET("/sso", func(c *gin.Context) {
		nonce, err := ExtractNonce(secret, c.Request)
		if err != nil {
			panic(err)
		}

		v := url.Values{}
		v.Set("nonce", nonce)
		v.Set("team", "test-dev")
		v.Set("id", "1")
		v.Set("avatar_url", "http://i.imgur.com/gYfoZpY.gif")
		v.Set("username", "finn")
		v.Set("email", "finn@ooo.com")
		v.Set("profile_url", "http://ooo.com/finn")
		v.Set("real_name", "Finn Mertens")

		raw := v.Encode()
		payload := base64.StdEncoding.EncodeToString([]byte(raw))
		url := "localhost:3000/sessions/sso?payload=" + url.QueryEscape(payload) + "&sig=" + Sign([]byte(secret), []byte(payload))

		c.Redirect(302, url)
	})

	r.Run(":8989")
}

func ExtractNonce(secret string, r *http.Request) (string, error) {
	r.ParseForm()
	payload := r.Form.Get("payload")
	sig := r.Form.Get("sig")

	if Sign([]byte(secret), []byte(payload)) != sig {
		return "", errors.New("bad signature")
	}

	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(string(raw))
	if err != nil {
		return "", err
	}

	return values.Get("nonce"), nil
}

func Sign(secret, payload []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
