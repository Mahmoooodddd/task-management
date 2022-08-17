package platform

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"time"
)

type JWTHandler interface {
	GetToken(params GetTokenParams) (string, error)
	GetUsernameFromToken(signedToken string) (string, error)
}

type jwtHandler struct {
	PrivateKey []byte
	PublicKey  []byte
	Passphrase string
}

type claims struct {
	Username string `json:"username"`
	Iat      int64  `json:"iat"`
	Exp      int64  `json:"exp"`
}

type GetTokenParams struct {
	Email string
}

func (c *claims) Valid() error {
	now := time.Now().Unix()
	if c.Exp < now {
		return fmt.Errorf("token expired")
	}
	return nil
}

func (jwtHandler *jwtHandler) GetToken(params GetTokenParams) (string, error) {
	rsaPrivateKey, err := ssh.ParseRawPrivateKeyWithPassphrase(jwtHandler.PrivateKey, []byte(jwtHandler.Passphrase))

	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(time.Minute * 600).Unix()
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &claims{
		Username: params.Email,
		Iat:      now.Unix(),
		Exp:      expiresAt,
	}
	tokenStr, err := t.SignedString(rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (jwtHandler *jwtHandler) GetUsernameFromToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &claims{}, func(token *jwt.Token) (i interface{}, e error) {
		rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(jwtHandler.PublicKey)
		return rsaPublicKey, err
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		if claims.Valid() != nil {
			return "", err
		}
		return claims.Username, nil
	}
	return "", err
}

func NewJWTHandler() JWTHandler {
	privateKey, err := ioutil.ReadFile("config/jwt/private.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err := ioutil.ReadFile("config/jwt/public.pem")
	if err != nil {
		panic(err)
	}
	return &jwtHandler{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Passphrase: "123456789",
	}
}
