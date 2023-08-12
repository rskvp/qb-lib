package jwt

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/signing"
)

var hmacSampleSecret []byte

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func TestCreate(t *testing.T) {

	nbf := time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix()

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := NewWithClaims(signing.SigningMethodHS256, elements.MapClaims{
		"foo": "bar",
		"nbf": nbf,
		"uid":"USER_1234",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	parsed, err := Parse(tokenString, func(token *elements.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*signing.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if nil==parsed{
		t.Error("TOKEN not parsed!")
		t.FailNow()
	}

	if claims, ok := parsed.Claims.(elements.MapClaims); ok && parsed.Valid {
		foo := claims["foo"].(string)
		uid := claims["uid"].(string)
		if foo != "bar" {
			t.Error("foo is not 'bar'")
			t.Fail()
		}
		fmt.Println(foo, claims["nbf"], uid)
	} else {
		t.Error(err)
		t.Fail()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func init() {
	// Load sample key data
	if keyData, e := ioutil.ReadFile("./_test/test/hmacTestKey"); e == nil {
		hmacSampleSecret = keyData
	} else {
		panic(e)
	}
}
