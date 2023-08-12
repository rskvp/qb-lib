package test

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/rskvp/qb-lib/qb_auth0/jwt"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/signing"
)

func LoadRSAPrivateKeyFromDisk(location string) *rsa.PrivateKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := signing.ParseRSAPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := signing.ParseRSAPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func MakeSampleToken(c elements.Claims, key interface{}) string {
	token := jwt.NewWithClaims(signing.SigningMethodRS256, c)
	s, e := token.SignedString(key)

	if e != nil {
		panic(e.Error())
	}

	return s
}
