package qb_auth0

import (
	"fmt"
	"testing"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_auth0/jwt"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
	"github.com/rskvp/qb-lib/qb_auth0/storage"
)

var DELEGATE_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzg0YmI4MWQyZmFiOTJkYWRiYzdiMjUwYzVjM2M1NzEiLCJwYXlsb2FkIjp7ImltcGVyc29uYXRlIjp0cnVlLCJyb2xlIjoic3VwZXIiLCJ0aW1lIjoiMjAyMC0xMC0yOSAxNzo0NDo1NC40NjYxMTIgKzAxMDAgQ0VUIG09KzAuMDExNDg1OTA3IiwidXNlcl9pZCI6IjM4NGJiODFkMmZhYjkyZGFkYmM3YjI1MGM1YzNjNTcxIn0sInNlY3JldF90eXBlIjoiYWNjZXNzIiwiZXhwIjo0Njc5ODI5OTQxLCJqdGkiOiJkMGY5ZjE1ZjViMGExNjBhMmE3NDQzMTI4YzE2NDM5OSJ9.gsdYg8FwqFHrIUF8EZZGvHdFlfuomAS-lsBEvUSzR-A"
var TEST_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZTYwMTIxODEyMDgyNTVhMzRkNTFjZjQwYWY1MDFiYTQiLCJwYXlsb2FkIjp7InJlZnJlc2hfdXVpZCI6IjMxOTJlODgwZjk1ZTFjYjE0YWFlZmZhZTM2YmIyMDUxIn0sInNlY3JldF90eXBlIjoicmVmcmVzaCIsImV4cCI6MTYwNTE5NDU4OCwianRpIjoiMzk5YTBiYjA4ZGY1NWU1NmM5MDZkZWRiZmNjZTZiNWMifQ.wH8gZt39L0TW22-HYj1NLKYAcuHNd1cks9JZ-d2c6Gc"
var REFRESH_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjlmMDU4ZWQzZjY4NmUxZGI5MDZlM2JiMjdmM2UyMDgiLCJwYXlsb2FkIjp7InJlZnJlc2hfdXVpZCI6ImQyZWMzOGUxODAzNTQ2YTE3OTFhMTdhZDE4NDQ5MWNlIn0sInNlY3JldF90eXBlIjoicmVmcmVzaCIsImV4cCI6MTY3NzkzMTA5NCwianRpIjoiNTc1YmZhNzQ3OWI2YzVhZjc2NGU3NWQzOTNmMWFiN2QifQ.c9jWPbF3P_qmtV1yHpqAmxIx3SxUNsvDFIWFIpQOFLI"

func Test_RegisterDouble(t *testing.T) {
	auth0 := getAuth("gorm")

	// open auth0 service
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	username := "user_to_remove_after_test"
	password := "password"

	// do register
	authResp := auth0.AuthSignUp(username, password, map[string]interface{}{"field1": "test value in payload"})
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}
	token := authResp.AccessToken
	defer removeUserAndClose(t, auth0, token)

	// do register
	authResp = auth0.AuthSignUp(username, password, map[string]interface{}{"field1": "test value in payload"})
	if len(authResp.Error) == 0 {
		t.Error("Expected 'already registered' error")
		t.FailNow()
	}
}

func removeUserAndClose(t *testing.T, auth0 *Auth0, token string) {
	err := auth0.AuthRemove(token)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	err = auth0.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	return
}

func Test_Token(t *testing.T) {
	auth0 := getAuth("gorm")

	// open auth0 service
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	username := "user_to_remove_after_test"
	password := "password"

	// do register
	authResp := auth0.AuthSignUp(username, password, map[string]interface{}{"field1": "test value in payload"})
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}
	fmt.Println("Registration:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	confirmToken := authResp.ConfirmToken

	// confirm
	authResp = auth0.AuthConfirm(confirmToken)
	fmt.Println("Confirm:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// signin
	authResp = auth0.AuthSignIn(username, password)
	fmt.Println("SignIn:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// refresh
	time.Sleep(5 * time.Second)
	authResp = auth0.TokenRefresh(authResp.RefreshToken)
	fmt.Println("Refresh:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// update
	payload := authResp.ItemPayload
	payload["timestamp"] = time.Now().Unix()
	payload["field2"] = "Added after update"
	newId, err := auth0.AuthUpdate(authResp.ItemId, payload)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	if newId != authResp.ItemId {
		t.Error("Expected same userId")
		t.FailNow()
	}

	// check password
	var currUser, currPassword string
	currUser, currPassword, err = auth0.AuthGetCredentials(newId)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Username:", currUser, "Password:", currPassword)

	// signin after update
	authResp = auth0.AuthSignIn(username, password)
	fmt.Println("SignIn After Payload Changed:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	oldId := authResp.ItemId
	// update password
	password = "this is new"
	authResp = auth0.AuthChangeLogin(authResp.ItemId, username, password)
	if len(authResp.Error) > 0 {
		t.Error(err)
		t.FailNow()
	}
	if oldId == authResp.ItemId {
		t.Error("Expected different userId")
		t.FailNow()
	}

	// check password
	currUser, currPassword, err = auth0.AuthGetCredentials(authResp.ItemId)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Username:", currUser, "Password:", currPassword)

	// signin after update
	authResp = auth0.AuthSignIn(username, password)
	fmt.Println("SignIn After Password Changed:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	//-- remove --//
	err = auth0.AuthRemove(authResp.AccessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// finally close
	err = auth0.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func Test_UserRegisterButForgetPasswordBeforeConfirm(t *testing.T) {
	auth0 := getAuth("gorm")

	// open auth0 service
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	username := "user_to_remove_after_test"
	password := "password"

	// do register
	authResp := auth0.AuthSignUp(username, password, map[string]interface{}{"field1": "test value in payload"})
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}
	fmt.Println("Registration:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// update password
	password = "this is new"
	authResp = auth0.AuthChangeLogin(authResp.ItemId, username, password)
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}

	// confirm
	confirmToken := authResp.ConfirmToken
	authResp = auth0.AuthConfirm(confirmToken)
	fmt.Println("Confirm:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// signin
	authResp = auth0.AuthSignIn(username, password)
	fmt.Println("SignIn:", authResp)
	err = parse(auth0, authResp)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	//-- remove --//
	err = auth0.AuthRemove(authResp.AccessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// finally close
	err = auth0.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func TestAuth0_Token(t *testing.T) {
	auth0 := getAuth("gorm")

	// add os env variable
	// _ = os.Setenv("foo", "Hello, this is a foo secret") // this variable is available also as secret
	// fmt.Println("Secrets", auth0.Secrets().String(), auth0.Secrets().Get("foo")) // request a secret from env
	// fmt.Println("Secrets after Get env", auth0.Secrets().String())
	// auth0.Secrets().Remove("foo") // remove foo secret
	// fmt.Println("Secrets after Remove env", auth0.Secrets().String())

	// open auth0 service
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	username := "user11"
	password := "password"
	new_password := "new_password"

	// do login
	authResp := auth0.AuthSignIn(username, password)
	if len(authResp.Error) > 0 && authResp.Error != storage.ErrorEntityDoesNotExists.Error() {
		t.Error(authResp.Error)
		t.FailNow()
	}
	fmt.Println("Found:", authResp)

	// parse tokens
	if len(authResp.AccessToken) > 0 {
		fmt.Println("-----------------------------")
		contentAccess, err := auth0.TokenParse(authResp.AccessToken)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println("CONTENT ACCESS: ", contentAccess, time.Unix(int64(contentAccess["exp"].(float64)), 0))
		// fmt.Println("EXPIRE:", time.Unix(int64(contentAccess["exp"].(float64)), 0))
		contentRefresh, err := auth0.TokenParse(authResp.RefreshToken)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println("CONTENT REFRESH: ", contentRefresh, time.Unix(int64(contentRefresh["exp"].(float64)), 0))

		// refresh
		authResp = auth0.TokenRefresh(authResp.RefreshToken)
		contentAccess, err = auth0.TokenParse(authResp.AccessToken)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println("CONTENT ACCESS: ", contentAccess, time.Unix(int64(contentAccess["exp"].(float64)), 0))

		fmt.Println("-----------------------------")
	}

	// register if does not exists
	if len(authResp.ItemPayload) == 0 {
		payload := map[string]interface{}{
			"role": "user",
		}
		// signup
		authResp = auth0.AuthSignUp(username, password, payload)
		if len(authResp.Error) > 0 {
			t.Error(authResp.Error)
			t.FailNow()
		}
		// confirm
		authResp = auth0.AuthConfirm(authResp.ConfirmToken)
		if len(authResp.Error) > 0 {
			t.Error(authResp.Error)
			t.FailNow()
		}
	}

	// update
	payload := map[string]interface{}{
		"role": "super",
		"time": time.Now().String(),
	}
	_, err = auth0.AuthUpdate("wrong_ID", payload)
	if nil == err {
		t.Error("Expected entity_dos_not_exists error")
		t.FailNow()
	}

	_, err = auth0.AuthUpdate(authResp.ItemId, payload)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	authResp = auth0.AuthSignIn(username, new_password)
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}

	_, err = auth0.AuthUpdate(authResp.ItemId, payload)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	authResp = auth0.AuthSignIn(username, password)
	if len(authResp.Error) > 0 {
		t.Error(authResp.Error)
		t.FailNow()
	}

	fmt.Println("Found after Update:", authResp)
	accessToken := authResp.AccessToken
	fmt.Println("Access Token:", accessToken)
	fmt.Println("Refresh Token:", authResp.RefreshToken)
	accessKey := []byte(auth0.Secrets().GetNotEmpty(AccessSecretName))
	token, err := jwt.Parse(accessToken, func(token *elements.Token) (interface{}, error) {

		return accessKey, nil
	})
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Parsed Access Token VALID:", token.Valid)
	claims := token.GetMapClaims()
	fmt.Println("Parsed Access Token CLAIMS:", claims)
	tpayload := claims["payload"]
	fmt.Println("Parsed Access Token CLAIMS.payload:", tpayload)

	// VALIDATE THE TOKEN
	valid, err := auth0.TokenValidate(accessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("VALIDATE:", valid)

	// finally close
	err = auth0.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
}

func TestParse(t *testing.T) {
	tokenToValidate := TEST_TOKEN //"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmY0MzE5ZWMwZTViNDA2ODMyNDNjYWNiZjFhODllMmEiLCJwYXlsb2FkIjp7InJlZnJlc2hfdXVpZCI6ImFhN2JhMTRlYjAxYjVjZGQ2ODJmM2YwYTgxNzhkYTg4In0sInNlY3JldF90eXBlIjoicmVmcmVzaCIsImV4cCI6MTYwNTEwNjY2MSwianRpIjoiMzg4MTJiMWNmNWEwYWQxNWEwNzI1MWM3ODhjODA5MDgifQ.RERPYDSsleDo0vbXXRW0k_k9doQ8EyZuSQP9udL4pZQ"

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// new(elements.Parser).Parse(tokenToValidate, elements.Keyfunc)

	// VALIDATE THE TOKEN
	data, err := auth0.TokenParse(tokenToValidate)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("PARSE:", data)
	fmt.Println("EXPIRE:", time.Unix(int64(data["exp"].(float64)), 0))
}

func TestDelegate(t *testing.T) {
	ownerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzg0YmI4MWQyZmFiOTJkYWRiYzdiMjUwYzVjM2M1NzEiLCJwYXlsb2FkIjp7InJvbGUiOiJzdXBlciIsInRpbWUiOiIyMDIwLTEwLTI5IDE3OjQ0OjU0LjQ2NjExMiArMDEwMCBDRVQgbT0rMC4wMTE0ODU5MDcifSwic2VjcmV0X3R5cGUiOiJhY2Nlc3MiLCJleHAiOjE2MDM5OTAxOTQsImp0aSI6ImFiOGE0YTgwN2IzYWEzNzliNmM2Y2Y3YmJhYjY4ZmFhIn0.TxtqIvTCpfSp8tlxJ3cnLisNO5IW0moCjSK1E3wA8Qs"

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// VALIDATE THE TOKEN
	response := auth0.AuthGrantDelegation(ownerToken)
	if len(response.Error) > 0 {
		t.Error(response.Error)
		t.FailNow()
	}
	fmt.Println("DELEGATE RESPONSE:", response)
	fmt.Println("Access Token:", response.AccessToken)
	fmt.Println("Refresh Token:", response.RefreshToken)
	fmt.Println("UserId:", response.ItemId)
	fmt.Println("Payload:", response.ItemPayload)

	// VALIDATE THE TOKEN
	valid, err := auth0.TokenValidate(response.AccessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("VALIDATE:", valid)
}

func TestDelegateAndRevoke(t *testing.T) {
	ownerToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzg0YmI4MWQyZmFiOTJkYWRiYzdiMjUwYzVjM2M1NzEiLCJwYXlsb2FkIjp7InJvbGUiOiJzdXBlciIsInRpbWUiOiIyMDIwLTEwLTI5IDExOjUzOjU2LjE0NDM5NSArMDEwMCBDRVQgbT0rMC4wMTI1ODcyMzUifSwic2VjcmV0X3R5cGUiOiJhY2Nlc3MiLCJleHAiOjE2MDM5NjkxMzYsImp0aSI6ImFiOGE0YTgwN2IzYWEzNzliNmM2Y2Y3YmJhYjY4ZmFhIn0.ICVelWBg6iZdkdh8-XW2FFD-i5skRg6Z00RgmJtFmR0"

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// VALIDATE THE TOKEN
	response := auth0.AuthGrantDelegation(ownerToken)
	if len(response.Error) > 0 {
		t.Error(response.Error)
		t.FailNow()
	}
	fmt.Println("DELEGATE RESPONSE:", response)
	fmt.Println("Access Token:", response.AccessToken)
	fmt.Println("Refresh Token:", response.RefreshToken)
	fmt.Println("UserId:", response.ItemId)
	fmt.Println("Payload:", response.ItemPayload)

	// VALIDATE THE TOKEN
	valid, err := auth0.TokenValidate(response.AccessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("VALIDATE:", valid)

	// REVOKE THE TOKEN
	err = auth0.AuthRevokeDelegation(response.AccessToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("TOKEN REVOKED:", DELEGATE_TOKEN)

	// VALIDATE THE TOKEN
	valid, err = auth0.TokenValidate(response.AccessToken)
	if nil == err {
		t.Error("Expected a non valid token")
		t.FailNow()
	}
	fmt.Println("DELEGATE TOKEN REVOKED: ", response.AccessToken)
}

func TestRefresh(t *testing.T) {
	refreshToken := REFRESH_TOKEN // "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzg0YmI4MWQyZmFiOTJkYWRiYzdiMjUwYzVjM2M1NzEiLCJwYXlsb2FkIjp7InJvbGUiOiJzdXBlciIsInRpbWUiOiIyMDIwLTEwLTI5IDA4OjUyOjA4LjA0MzY3MiArMDEwMCBDRVQgbT0rMC4wMTE3MzM3MzgifSwic2VjcmV0X3R5cGUiOiJhY2Nlc3MiLCJleHAiOjE2MDM5NTgyMjgsImp0aSI6ImVlODlkNmJjY2MzMzUzMGZjZWY0MDVhMzliZjY2NzM1In0.CZZfjY6nCJ4V0L1m30-kCkjxoaf36Eb_2BsNixbtmUE"

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// VALIDATE THE TOKEN
	response := auth0.TokenRefresh(refreshToken)
	if len(response.Error) > 0 {
		t.Error(response.Error)
		t.FailNow()
	}
	fmt.Println("RESPONSE:", response)

	parse(auth0, response)
}

func TestRefreshAccess(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjlmMDU4ZWQzZjY4NmUxZGI5MDZlM2JiMjdmM2UyMDgiLCJwYXlsb2FkIjp7ImNvbmZpcm1lZCI6dHJ1ZSwidXNlcl9uYW1lIjoiZW5jLTY3ME1scFF3QW9UNTMrTGV5U3JUYVVCYms4ZHVETGpYS1NuMUh4cm0xeUJ6dW5mOGlUNldjekNaYWtWSUVKUk5VeEhKKyt3SSIsInVzZXJfcHN3IjoiZW5jLXNueGUrUnJPR0dFQ053ZmQ0TERraGI5SjExYlRCZEFvVVB5MlhLaXovVTNORnc9PSIsInVzZXJfcHN3X3RpbWVzdGFtcCI6MTY3MzUxMjcxN30sInNlY3JldF90eXBlIjoiYWNjZXNzIiwiZXhwIjoxNjc3ODQ2NDk0LCJqdGkiOiJkMmVjMzhlMTgwMzU0NmExNzkxYTE3YWQxODQ0OTFjZSJ9.-8FIy0au3oLrsfJU-T5MB-0G67j4a2LgQ9fWF9L7xeo"
	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZjlmMDU4ZWQzZjY4NmUxZGI5MDZlM2JiMjdmM2UyMDgiLCJwYXlsb2FkIjp7InJlZnJlc2hfdXVpZCI6ImQyZWMzOGUxODAzNTQ2YTE3OTFhMTdhZDE4NDQ5MWNlIn0sInNlY3JldF90eXBlIjoicmVmcmVzaCIsImV4cCI6MTY3NzkzMTA5NCwianRpIjoiNTc1YmZhNzQ3OWI2YzVhZjc2NGU3NWQzOTNmMWFiN2QifQ.c9jWPbF3P_qmtV1yHpqAmxIx3SxUNsvDFIWFIpQOFLI"

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// VALIDATE THE TOKEN
	response := auth0.TokenRefreshAccess(accessToken, refreshToken)
	if len(response.Error) > 0 {
		t.Error(response.Error)
		t.FailNow()
	}
	fmt.Println("RESPONSE:", response)

	parse(auth0, response)
}

func TestValidate(t *testing.T) {
	tokenToValidate := DELEGATE_TOKEN

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// new(elements.Parser).Parse(tokenToValidate, elements.Keyfunc)

	// VALIDATE THE TOKEN
	valid, err := auth0.TokenValidate(tokenToValidate)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("VALIDATE:", valid)
}

func TestRevoke(t *testing.T) {
	delegateToken := DELEGATE_TOKEN

	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	// REVOKE THE TOKEN
	err = auth0.AuthRevokeDelegation(delegateToken)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("TOKEN REVOKED:", DELEGATE_TOKEN)
}

func TestAuth0_Config(t *testing.T) {
	config, err := Auth0ConfigLoad("auth0.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("configuration", config)

	dsn := config.CacheStorage.Dsn
	db := storage.NewDriverBolt(dsn)
	if !db.Enabled() {
		t.Error("expected database was enabled")
		t.FailNow()
	}

	txt, err := qbc.IO.ReadTextFromFile("auth0.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	config = Auth0ConfigParse(txt)
	fmt.Println("configuration parsed: ", config)
}

func TestPasswordExpire(t *testing.T) {

	timestamp := time.Now().Unix() - 60*60*24 // 1 day ago
	payload := map[string]interface{}{
		FLD_USERPASSWORD_TIMESTAMP: timestamp,
	}
	expired := isPasswordExpired(payload, 2)
	fmt.Println(expired)
}

func TestParseClaims(t *testing.T) {
	auth0 := getAuth("gorm")
	err := auth0.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZWQ0M2IzNGE0M2YyNDYwZjJjYWVhNjkzYzI4Mzc2MzYiLCJwYXlsb2FkIjp7InVpZCI6ImVkNDNiMzRhNDNmMjQ2MGYyY2FlYTY5M2MyODM3NjM2In0sImV4cCI6NDgxODY2MTEzMywianRpIjoiODdiOWQ5YWNmYWU1YzIyNjM4ZmY5MGM1NDM4OTRhOTAifQ.kKNQFj_Csua2PPl6JehrEVg7fnTvnJz26xBQ5h1NvXs"

	data := auth0.TokenClaimsNoValidate(token)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(data)
}

// ----------------------------------------------------------------------------------------------------------------------
//
//	p r i v a t e
//
// ----------------------------------------------------------------------------------------------------------------------
func getAuth(dbName string) *Auth0 {
	qbc.Paths.SetWorkspacePath("./")

	auth0 := NewAuth0(getConfig(dbName))
	auth0.Secrets().Put(AuthSecretName, "this-is-token-to-authenticate")
	auth0.Secrets().Put(AccessSecretName, "hsdfuhksdhf5435khjsd")
	auth0.Secrets().Put(RefreshSecretName, "hsdfuhqswe34qwksdhfkhjsd")
	return auth0
}

func getConfig(name string) *Auth0Config {
	switch name {
	case "bolt":
		return getConfigBolt()
	case "arango":
		return getConfigArango()
	default:
		return getConfigGorm()
	}
}

func getConfigGorm() *Auth0Config {
	config := Auth0ConfigNew()
	config.CacheStorage.Driver = "sqlite"
	config.CacheStorage.Dsn = "../_test/auth0/data/auth.db"
	config.AuthStorage.Driver = "sqlite"
	config.AuthStorage.Dsn = "../_test/auth0/data/auth.db"
	return config
}

func getConfigBolt() *Auth0Config {
	config := Auth0ConfigNew()
	config.CacheStorage.Driver = "bolt"
	config.CacheStorage.Dsn = "root:root@file:../_test/auth0/data/cacheDb.dat"
	config.AuthStorage.Driver = "bolt"
	config.AuthStorage.Dsn = "root:root@file:../_test/auth0/data/usersDb.dat"
	return config
}

func getConfigArango() *Auth0Config {
	config := Auth0ConfigNew()
	config.CacheStorage.Driver = "arango"
	config.CacheStorage.Dsn = "root:!qaz2WSX098@tcp(localhost:8529)/test"
	config.AuthStorage.Driver = "arango"
	config.AuthStorage.Dsn = "root:!qaz2WSX098@tcp(localhost:8529)/test"
	return config
}

func parse(auth0 *Auth0, authResp *Auth0Response) error {
	fmt.Println("-----------------------------")

	fmt.Println("PAYLOAD: ", qbc.JSON.Stringify(authResp.ItemPayload))

	contentAccess, err := auth0.TokenParse(authResp.AccessToken)
	if nil != err {
		return err
	}
	expireAccess := time.Unix(int64(contentAccess["exp"].(float64)), 0)
	fmt.Println("\tACCESS:\t\t", expireAccess, "claims: ", qbc.JSON.Stringify(contentAccess))

	contentRefresh, err := auth0.TokenParse(authResp.RefreshToken)
	if nil != err {
		return err
	}
	expireRefresh := time.Unix(int64(contentRefresh["exp"].(float64)), 0)
	fmt.Println("\tREFRESH:\t", expireRefresh, "claims: ", qbc.JSON.Stringify(contentRefresh))

	if len(authResp.ConfirmToken) > 0 {
		contentConfirm, err := auth0.TokenParse(authResp.ConfirmToken)
		if nil != err {
			return err
		}
		expireConfirm := time.Unix(int64(contentConfirm["exp"].(float64)), 0)
		fmt.Println("\tCONFIRM:\t", expireConfirm, "claims: ", qbc.JSON.Stringify(contentConfirm))
	}

	fmt.Println("-----------------------------")

	return nil
}
