package _test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rskvp/qb-lib/qb_auth0/jwt"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/signing"
)

// Example (atypical) using the StandardClaims type by itself to parse a token.
// The StandardClaims type is designed to be embedded into your custom types
// to provide standard validation features.  You can use it alone, but there's
// no way to retrieve other fields after parsing.
// See the CustomClaimsType example for intended usage.
func ExampleNewWithClaims_standardClaims() {
	mySigningKey := []byte("AllYourBase")

	// Create the Claims
	claims := &elements.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(signing.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.QsODzZu3lUZMVdhbO76u3Jv02iYCvEHcYVUI1kOWEU0 <nil>
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ExampleNewWithClaims_customClaimsType() {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		elements.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		"bar",
		elements.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(signing.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", ss, err)
	//Output: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c <nil>
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func TestExampleParseWithClaims_customClaimsType(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		elements.StandardClaims
	}

	// sample token is expired.  override time so it parses as valid
	at(time.Unix(0, 0), func() {
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *elements.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			fmt.Printf("%v %v", claims.Foo, claims.StandardClaims.ExpiresAt)
		} else {
			fmt.Println(err)
		}
	})

	// Output: bar 15000
}

// Override time value for tests.  Restore default value after.
func at(t time.Time, f func()) {
	commons.TimeFunc = func() time.Time {
		return t
	}
	f()
	commons.TimeFunc = time.Now
}

// An example of parsing the error types using bitfield checks
func TestExampleParse_errorChecking(t *testing.T) {
	// Token from another example.  This token is expired
	var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"

	token, err := jwt.Parse(tokenString, func(token *elements.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if nil != err {
		t.Error(err)
	}

	if token.Valid {
		fmt.Println("You look nice today")
	} else if ve, ok := err.(*commons.ValidationError); ok {
		if ve.Errors&commons.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
		} else if ve.Errors&(commons.ValidationErrorExpired|commons.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
	}

	// Output: Timing is everything
}
