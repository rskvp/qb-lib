package jwt

import (
	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func Parse(tokenString string, keyFunc elements.Keyfunc) (*elements.Token, error) {
	return new(elements.Parser).Parse(tokenString, keyFunc)
}

func ParseWithClaims(tokenString string, claims elements.Claims, keyFunc elements.Keyfunc) (*elements.Token, error) {
	return new(elements.Parser).ParseWithClaims(tokenString, claims, keyFunc)
}

// Create a new Token.  Takes a signing method
func New(method commons.SigningMethod) *elements.Token {
	return NewWithClaims(method, elements.MapClaims{})
}

func NewWithClaims(method commons.SigningMethod, claims elements.Claims) *elements.Token {
	return &elements.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}
