package _test

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/rskvp/qb-lib/qb_auth0/jwt"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/_test/test"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/commons"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
)

var keyFuncError error = fmt.Errorf("error loading key")

var (
	jwtTestDefaultKey *rsa.PublicKey
	defaultKeyFunc    elements.Keyfunc = func(t *elements.Token) (interface{}, error) { return jwtTestDefaultKey, nil }
	emptyKeyFunc      elements.Keyfunc = func(t *elements.Token) (interface{}, error) { return nil, nil }
	errorKeyFunc      elements.Keyfunc = func(t *elements.Token) (interface{}, error) { return nil, keyFuncError }
	nilKeyFunc        elements.Keyfunc = nil
)

func init() {
	jwtTestDefaultKey = test.LoadRSAPublicKeyFromDisk("test/sample_key.pub")
}

var jwtTestData = []struct {
	name        string
	tokenString string
	keyfunc     elements.Keyfunc
	claims      elements.Claims
	valid       bool
	errors      uint32
	parser      *elements.Parser
}{
	{
		"basic",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar"},
		true,
		0,
		nil,
	},
	{
		"basic expired",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		commons.ValidationErrorExpired,
		nil,
	},
	{
		"basic nbf",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		commons.ValidationErrorNotValidYet,
		nil,
	},
	{
		"expired and nbf",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		commons.ValidationErrorNotValidYet | commons.ValidationErrorExpired,
		nil,
	},
	{
		"basic invalid",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar"},
		false,
		commons.ValidationErrorSignatureInvalid,
		nil,
	},
	{
		"basic nokeyfunc",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		nilKeyFunc,
		elements.MapClaims{"foo": "bar"},
		false,
		commons.ValidationErrorUnverifiable,
		nil,
	},
	{
		"basic nokey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		emptyKeyFunc,
		elements.MapClaims{"foo": "bar"},
		false,
		commons.ValidationErrorSignatureInvalid,
		nil,
	},
	{
		"basic errorkey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		errorKeyFunc,
		elements.MapClaims{"foo": "bar"},
		false,
		commons.ValidationErrorUnverifiable,
		nil,
	},
	{
		"invalid signing method",
		"",
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar"},
		false,
		commons.ValidationErrorSignatureInvalid,
		&elements.Parser{ValidMethods: []string{"HS256"}},
	},
	{
		"valid signing method",
		"",
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar"},
		true,
		0,
		&elements.Parser{ValidMethods: []string{"RS256", "HS256"}},
	},
	{
		"JSON Number",
		"",
		defaultKeyFunc,
		elements.MapClaims{"foo": json.Number("123.4")},
		true,
		0,
		&elements.Parser{UseJSONNumber: true},
	},
	{
		"Standard Claims",
		"",
		defaultKeyFunc,
		&elements.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
		true,
		0,
		&elements.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic expired",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		commons.ValidationErrorExpired,
		&elements.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic nbf",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		false,
		commons.ValidationErrorNotValidYet,
		&elements.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - expired and nbf",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100)), "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		commons.ValidationErrorNotValidYet | commons.ValidationErrorExpired,
		&elements.Parser{UseJSONNumber: true},
	},
	{
		"SkipClaimsValidation during token parsing",
		"", // autogen
		defaultKeyFunc,
		elements.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		true,
		0,
		&elements.Parser{UseJSONNumber: true, SkipClaimsValidation: true},
	},
}

func TestParser_Parse(t *testing.T) {
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = test.MakeSampleToken(data.claims, privateKey)
		}

		// Parse the token
		var token *elements.Token
		var err error
		var parser = data.parser
		if parser == nil {
			parser = new(elements.Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case elements.MapClaims:
			token, err = parser.ParseWithClaims(data.tokenString, elements.MapClaims{}, data.keyfunc)
		case *elements.StandardClaims:
			token, err = parser.ParseWithClaims(data.tokenString, &elements.StandardClaims{}, data.keyfunc)
		}

		// Verify result matches expectation
		if !reflect.DeepEqual(data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}

		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}

		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}

		if (err == nil && !token.Valid) || (err != nil && token.Valid) {
			t.Errorf("[%v] Inconsistent behavior between returned error and token.Valid", data.name)
		}

		if data.errors != 0 {
			if err == nil {
				t.Errorf("[%v] Expecting error.  Didn't get one.", data.name)
			} else {

				ve := err.(*commons.ValidationError)
				// compare the bitfield part of the error
				if e := ve.Errors; e != data.errors {
					t.Errorf("[%v] Errors don't match expectation.  %v != %v", data.name, e, data.errors)
				}

				if err.Error() == keyFuncError.Error() && ve.Inner != keyFuncError {
					t.Errorf("[%v] Inner error does not match expectation.  %v != %v", data.name, ve.Inner, keyFuncError)
				}
			}
		}
		if data.valid && token.Signature == "" {
			t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
		}
	}
}

func TestParser_ParseUnverified(t *testing.T) {
	privateKey := test.LoadRSAPrivateKeyFromDisk("test/sample_key")

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = test.MakeSampleToken(data.claims, privateKey)
		}

		// Parse the token
		var token *elements.Token
		var err error
		var parser = data.parser
		if parser == nil {
			parser = new(elements.Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case elements.MapClaims:
			token, _, err = parser.ParseUnverified(data.tokenString, elements.MapClaims{})
		case *elements.StandardClaims:
			token, _, err = parser.ParseUnverified(data.tokenString, &elements.StandardClaims{})
		}

		if err != nil {
			t.Errorf("[%v] Invalid token", data.name)
		}

		// Verify result matches expectation
		if !reflect.DeepEqual(data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}

		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}
	}
}

// Helper method for benchmarking various methods
func benchmarkSigning(b *testing.B, method commons.SigningMethod, key interface{}) {
	t := jwt.New(method)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := t.SignedString(key); err != nil {
				b.Fatal(err)
			}
		}
	})

}
