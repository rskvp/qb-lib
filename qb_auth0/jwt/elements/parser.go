package elements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Parser struct {
	ValidMethods         []string // If populated, only these methods will be considered valid
	UseJSONNumber        bool     // Use JSON Number format in JSON decoder
	SkipClaimsValidation bool     // Skip claims validation during token parsing
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return p.ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}

func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	token, parts, err := p.ParseUnverified(tokenString, claims)
	if err != nil {
		return token, err
	}

	// Verify signing method is in the required set
	if p.ValidMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.ValidMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return token, commons.NewValidationError(fmt.Sprintf("signing method %v is invalid", alg), commons.ValidationErrorSignatureInvalid)
		}
	}

	// Lookup key
	var key interface{}
	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, commons.NewValidationError("no Keyfunc was provided.", commons.ValidationErrorUnverifiable)
	}
	if key, err = keyFunc(token); err != nil {
		// keyFunc returned an error
		if ve, ok := err.(*commons.ValidationError); ok {
			return token, ve
		}
		return token, &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorUnverifiable}
	}

	vErr := &commons.ValidationError{}

	// Validate Claims
	if !p.SkipClaimsValidation {
		if err := token.Claims.Valid(); err != nil {

			// If the Claims Valid returned an error, check if it is a validation error,
			// If it was another error type, create a ValidationError with a generic ClaimsInvalid flag set
			if e, ok := err.(*commons.ValidationError); !ok {
				vErr = &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorClaimsInvalid}
			} else {
				vErr = e
			}
		}
	}

	// Perform validation
	token.Signature = parts[2]
	if err = token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, key); err != nil {
		vErr.Inner = err
		vErr.Errors |= commons.ValidationErrorSignatureInvalid
	}

	if vErr.Valid() {
		token.Valid = true
		return token, nil
	}

	return token, vErr
}

// WARNING: Don't use this method unless you know what you're doing
//
// This method parses the token but doesn't validate the signature. It's only
// ever useful in cases where you know the signature is valid (because it has
// been checked previously in the stack) and you want to extract values from
// it.
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error) {
	parts = strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, parts, commons.NewValidationError("token contains an invalid number of segments", commons.ValidationErrorMalformed)
	}

	token = &Token{Raw: tokenString}

	// parse Header
	var headerBytes []byte
	if headerBytes, err = commons.DecodeSegment(parts[0]); err != nil {
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			return token, parts, commons.NewValidationError("tokenstring should not contain 'bearer '", commons.ValidationErrorMalformed)
		}
		return token, parts, &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorMalformed}
	}
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		return token, parts, &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorMalformed}
	}

	// parse Claims
	var claimBytes []byte
	token.Claims = claims

	if claimBytes, err = commons.DecodeSegment(parts[1]); err != nil {
		return token, parts, &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorMalformed}
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	if p.UseJSONNumber {
		dec.UseNumber()
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	if c, ok := token.Claims.(MapClaims); ok {
		err = dec.Decode(&c)
	} else {
		err = dec.Decode(&claims)
	}
	// Handle decode error
	if err != nil {
		return token, parts, &commons.ValidationError{Inner: err, Errors: commons.ValidationErrorMalformed}
	}

	// Lookup signature method
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = commons.GetSigningMethod(method); token.Method == nil {
			return token, parts, commons.NewValidationError("signing method (alg) is unavailable.", commons.ValidationErrorUnverifiable)
		}
	} else {
		return token, parts, commons.NewValidationError("signing method (alg) is unspecified.", commons.ValidationErrorUnverifiable)
	}

	return token, parts, nil
}
