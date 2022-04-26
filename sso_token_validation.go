package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func ValidateSSOToken(tokenStr string, ssoJwksUri string) (bool, error) {

	token, _ := jwt.Parse(tokenStr, nil)
	if token == nil {
		return false, errors.New("invalid_token")
	}

	cert, err := GetSSOPemCert(token, ssoJwksUri)
	if err != nil {
		return false, err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return false, err
	}

	parts := strings.Split(tokenStr, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetSSOPemCert(token *jwt.Token, ssoJwksUri string) (string, error) {
	cert := ""
	resp, err := http.Get(ssoJwksUri)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
