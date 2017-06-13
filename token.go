package corkboardauth

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//CustomClaims are the the claims that will be in the JWT token returned to the authenticated user
type CustomClaims struct {
	Email string `json:"email"`
	UID   string `json:"uid"`
	jwt.StandardClaims
}

func (cba *CorkboardAuth) generateUserToken(user *User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &CustomClaims{
		Email: user.Email,
		UID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(36 * time.Hour).Unix(), //TODO: Come up with an expiration time
			Issuer:    "CorkboardAuth",
		},
	})
	return token.SignedString(cba.privateKey)
}

func (cba *CorkboardAuth) getPublicPem() ([]byte, error) {
	pubDer, err := x509.MarshalPKIXPublicKey(cba.privateKey.Public())
	if err != nil {
		return nil, err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDer,
	})
	return pubPem, nil
}
