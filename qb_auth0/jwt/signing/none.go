package signing

import (
	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
)

// Implements the none signing method.  This is required by the spec
// but you probably should never use it.
var SigningMethodNone *signingMethodNone

const UnsafeAllowNoneSignatureType unsafeNoneMagicConstant = "none signing method allowed"

var NoneSignatureTypeDisallowedError error

type signingMethodNone struct{}
type unsafeNoneMagicConstant string



func (m *signingMethodNone) Alg() string {
	return "none"
}

// Only allow 'none' alg type if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Verify(signingString, signature string, key interface{}) (err error) {
	// Key must be UnsafeAllowNoneSignatureType to prevent accidentally
	// accepting 'none' signing method
	if _, ok := key.(unsafeNoneMagicConstant); !ok {
		return NoneSignatureTypeDisallowedError
	}
	// If signing method is none, signature must be an empty string
	if signature != "" {
		return commons.NewValidationError(
			"'none' signing method with non-empty signature",
			commons.ValidationErrorSignatureInvalid,
		)
	}

	// Accept 'none' signing method.
	return nil
}

// Only allow 'none' signing if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Sign(signingString string, key interface{}) (string, error) {
	if _, ok := key.(unsafeNoneMagicConstant); ok {
		return "", nil
	}
	return "", NoneSignatureTypeDisallowedError
}

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	SigningMethodNone = &signingMethodNone{}
	NoneSignatureTypeDisallowedError = commons.NewValidationError("'none' signature type is not allowed", commons.ValidationErrorSignatureInvalid)

	commons.RegisterSigningMethod(SigningMethodNone.Alg(), func() commons.SigningMethod {
		return SigningMethodNone
	})
}