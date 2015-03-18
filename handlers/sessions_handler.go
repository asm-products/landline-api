package handlers

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/asm-products/landline-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SessionsNew(c *gin.Context) {
	c.Request.ParseForm()
	team := models.FindTeamBySlug(c.Request.Form.Get("team"))
	if team == nil {
		c.String(404, "Not found")
		return
	}

	nonce, err := models.CreateNonce()
	if err != nil {
		panic(err)
	}
	raw := "nonce=" + nonce.Nonce + "&uid=" + c.Request.Form.Get("uid")
	payload := base64.StdEncoding.EncodeToString([]byte(raw))

	url := team.SSOUrl + "?payload=" + url.QueryEscape(payload) + "&sig=" + models.Sign([]byte(team.SSOSecret), []byte(payload))

	c.Redirect(302, url)
}

func SessionsLoginSSO(c *gin.Context) {
	r, err := ExtractSSORequest(c.Request)
	if err != nil {
		panic(err)
	}

	if !models.NonceValid(r.Nonce) {
		c.String(403, "Invalid nonce")
		return
	}

	team := models.FindTeamBySlug(r.TeamSlug)

	if !r.IsValid(team.SSOSecret) {
		c.String(403, "Not authorized")
		return
	}

	u := &models.User{
		TeamId:     team.Id,
		AvatarUrl:  r.AvatarUrl,
		Email:      r.Email,
		ExternalId: r.ExternalId,
		ProfileUrl: r.ProfileUrl,
		RealName:   r.RealName,
		Username:   r.Username,
	}

	u, err = models.FindOrCreateUserByExternalId(u)
	if err != nil {
		panic(err)
	}

	token := GenerateToken(u.Id)

	c.JSON(200, gin.H{"token": token})
}

func ExtractSSORequest(r *http.Request) (*models.SSORequest, error) {
	r.ParseForm()
	payload := r.Form.Get("payload")
	sig := r.Form.Get("sig")

	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(raw))
	if err != nil {
		return nil, err
	}

	sso := &models.SSORequest{
		Payload:    payload,
		Signature:  sig,
		Nonce:      values.Get("nonce"),
		TeamSlug:   values.Get("team"),
		ExternalId: values.Get("id"),
		AvatarUrl:  values.Get("avatar_url"),
		Username:   values.Get("username"),
		Email:      values.Get("email"),
		ProfileUrl: values.Get("profile_url"),
		RealName:   values.Get("real_name"),
	}

	return sso, nil
}

func GenerateToken(id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = id
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		panic(err)
	}
	return tokenString
}
