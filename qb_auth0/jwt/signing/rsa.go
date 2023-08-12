package signing

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
)

// Implements the RSA family of signing methods signing methods
// Expects *rsa.PrivateKey for signing and *rsa.PublicKey for validation
type SigningMethodRSA struct {
	Name string
	Hash crypto.Hash
}

// Specific instances for RS256 and company
var (
	SigningMethodRS256 *SigningMethodRSA
	SigningMethodRS384 *SigningMethodRSA
	SigningMethodRS512 *SigningMethodRSA
)

func (m *SigningMethodRSA) Alg() string {
	return m.Name
}

// Implements the Verify method from SigningMethod
// For this signing method, must be an *rsa.PublicKey structure.
func (m *SigningMethodRSA) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = commons.DecodeSegment(signature); err != nil {
		return err
	}

	var rsaKey *rsa.PublicKey
	var ok bool

	if rsaKey, ok = key.(*rsa.PublicKey); !ok {
		return commons.ErrInvalidKeyType
	}

	// Create hasher
	if !m.Hash.Available() {
		return commons.ErrHashUnavailable
	}
	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Verify the signature
	return rsa.VerifyPKCS1v15(rsaKey, m.Hash, hasher.Sum(nil), sig)
}

// Implements the Sign method from SigningMethod
// For this signing method, must be an *rsa.PrivateKey structure.
func (m *SigningMethodRSA) Sign(signingString string, key interface{}) (string, error) {
	var rsaKey *rsa.PrivateKey
	var ok bool

	// Validate type of key
	if rsaKey, ok = key.(*rsa.PrivateKey); !ok {
		return "", commons.ErrInvalidKey
	}

	// Create the hasher
	if !m.Hash.Available() {
		return "", commons.ErrHashUnavailable
	}

	hasher := m.Hash.New()
	hasher.Write([]byte(signingString))

	// Sign the string and return the encoded bytes
	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, m.Hash, hasher.Sum(nil)); err == nil {
		return commons.EncodeSegment(sigBytes), nil
	} else {
		return "", err
	}
}


//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	// RS256
	SigningMethodRS256 = &SigningMethodRSA{"RS256", crypto.SHA256}
	commons.RegisterSigningMethod(SigningMethodRS256.Alg(), func() commons.SigningMethod {
		return SigningMethodRS256
	})

	// RS384
	SigningMethodRS384 = &SigningMethodRSA{"RS384", crypto.SHA384}
	commons.RegisterSigningMethod(SigningMethodRS384.Alg(), func() commons.SigningMethod {
		return SigningMethodRS384
	})

	// RS512
	SigningMethodRS512 = &SigningMethodRSA{"RS512", crypto.SHA512}
	commons.RegisterSigningMethod(SigningMethodRS512.Alg(), func() commons.SigningMethod {
		return SigningMethodRS512
	})
}
