package qb_auth0

import (
	"errors"
	"fmt"
	"time"

	"github.com/rskvp/qb-lib/qb_auth0/jwt"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/elements"
	"github.com/rskvp/qb-lib/qb_auth0/jwt/signing"
	"github.com/rskvp/qb-lib/qb_auth0/storage"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrorMissingSecureKey     = errors.New("missing_secure_key")
	ErrorMissingClaims        = errors.New("missing_claims")
	ErrorUnauthorized         = errors.New("unauthorized_401")
	ErrorNotConfirmed         = errors.New("not_confirmed")
	ErrorMalformedAccountData = errors.New("malformed_account_data")
	ErrorPasswordExpired      = errors.New("password_expired")

	accessTD         = 30 * time.Minute
	accessCD         = 24 * time.Hour // same as refresh token
	refreshTD        = 24 * time.Hour
	infiniteDuration = 100 * 356 * 24 * time.Hour
)

const (
	CACHE_KEY                  = "jti"
	FLD_USERID                 = "user_id"
	FLD_PAYLOAD                = "payload"
	FLD_CONFIRMED              = "confirmed" // user confirmed account
	FLD_USERNAME               = "user_name"
	FLD_USERPASSWORD           = "user_psw"
	FLD_USERPASSWORD_TIMESTAMP = "user_psw_timestamp" // last change timestamp
	FLD_SECRET_TYPE            = "secret_type"
	FLD_EXP                    = "exp"
)

var PROTECTED_FIELDS = []string{
	CACHE_KEY, FLD_USERNAME, FLD_USERPASSWORD, FLD_SECRET_TYPE, FLD_CONFIRMED,
}

const (
	TAccess = iota
	TRefresh
	TConfirm
	TDelegate
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Auth0Claims struct {
	UserId     string                 `json:"user_id,omitempty"`
	Payload    map[string]interface{} `json:"payload,omitempty"`
	SecretType string                 `json:"secret_type,omitempty"`
	elements.StandardClaims
}

type Auth0 struct {
	PasswordDurationDays  int
	AccessTokenDuration   time.Duration
	AccessCacheDuration   time.Duration
	RefreshTokenDuration  time.Duration
	ConfirmTokenDuration  time.Duration
	DelegateTokenDuration time.Duration

	config  *Auth0Config
	secrets Auth0ConfigSecrets
	cacheDb storage.IDatabase // cache database
	authDb  storage.IDatabase // authentication database
}

func NewAuth0(config ...interface{}) *Auth0 {
	instance := new(Auth0)

	instance.PasswordDurationDays = 0 // never expires
	instance.AccessTokenDuration = accessTD
	instance.AccessCacheDuration = accessCD
	instance.RefreshTokenDuration = refreshTD
	instance.ConfirmTokenDuration = infiniteDuration
	instance.DelegateTokenDuration = infiniteDuration

	instance.config = Auth0ConfigNew()
	instance.secrets = Auth0ConfigSecrets{}

	if len(config) == 1 {
		if c, b := config[0].(*Auth0Config); b {
			instance.config = c
		} else if c, b := config[0].(Auth0Config); b {
			instance.config = &c
		} else if dsn, b := config[0].(storage.Dsn); b {
			instance.config.CacheStorage.Dsn = dsn.String()
		} else if dsn, b := config[0].(*storage.Dsn); b {
			instance.config.CacheStorage.Dsn = dsn.String()
		}
	}

	if nil != instance.config.Secrets {
		instance.secrets = instance.config.Secrets
	}

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Auth0) Secrets() Auth0ConfigSecrets {
	if nil != instance && nil != instance.secrets {
		return instance.secrets
	}
	return nil
}

func (instance *Auth0) Open() (err error) {
	if nil != instance && nil != instance.config {
		// init cache database
		if len(instance.config.CacheStorage.Dsn) > 0 && len(instance.config.CacheStorage.Driver) > 0 {
			instance.cacheDb, err = storage.NewDatabase(instance.config.CacheStorage.Driver,
				instance.config.CacheStorage.Dsn, true)
		}
		// init auth database
		if len(instance.config.AuthStorage.Dsn) > 0 && len(instance.config.AuthStorage.Driver) > 0 {
			instance.authDb, err = storage.NewDatabase(instance.config.AuthStorage.Driver,
				instance.config.AuthStorage.Dsn, false)
		}
	}
	return err
}

func (instance *Auth0) Close() (err error) {
	if nil != instance {
		if instance.isCacheEnabled() {
			err = instance.cacheDb.Close()
		}
	}
	return err
}

//----------------------------------------------------------------------------------------------------------------------
//	A u t h e n t i c a t i o n
//----------------------------------------------------------------------------------------------------------------------

// AuthSignIn try to log in and returns itemId, itemPayload, error
func (instance *Auth0) AuthSignIn(cleanUsername, cleanPassword string) *Auth0Response {
	response := new(Auth0Response)
	response.ItemId = ""
	response.ItemPayload = map[string]interface{}{}
	if nil != instance && instance.isAuthDbEnabled() {
		secretKey := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
		if len(secretKey) > 0 {
			// get user from database
			userId, payload, err := instance.getDecPayload(secretKey, cleanUsername, cleanPassword)
			if nil != err {
				response.Error = err.Error()
			} else {
				// confirmed????
				confirmed := qbc.Reflect.GetBool(payload, FLD_CONFIRMED)
				if !confirmed {
					response.Error = ErrorNotConfirmed.Error()
				} else {
					expired := isPasswordExpired(payload, instance.PasswordDurationDays)
					if expired {
						response.Error = ErrorPasswordExpired.Error()
					} else {
						response.ItemId = userId
						instance.setPayload(response.ItemPayload, payload) // user data

						adbKey, aToken, err := instance.generateToken(TAccess, userId, payload)
						if nil != err {
							response.Error = err.Error()
						} else {
							response.AccessToken = aToken
							// REFRESH token
							_, rToken, err := instance.generateToken(TRefresh, userId, map[string]interface{}{"refresh_uuid": adbKey})
							if nil != err {
								response.Error = err.Error()
							} else {
								response.RefreshToken = rToken
							}
						}
					}
				}
			}
		}
	} else {
		response.Error = ErrorUnauthorized.Error()
	}
	return response
}

func (instance *Auth0) AuthSignUp(cleanUsername, cleanPassword string, cleanPayload map[string]interface{}) *Auth0Response {
	response := new(Auth0Response)
	response.ItemPayload = map[string]interface{}{}
	if nil == cleanPayload {
		cleanPayload = map[string]interface{}{}
	}
	cleanPayload[FLD_CONFIRMED] = false // new user (need confirm with confirmToken)
	// cleanPayload[FLD_USERNAME] = cleanUsername
	// cleanPayload[FLD_USERPASSWORD] = cleanPassword
	instance.setUsername(cleanPayload, cleanUsername)
	instance.setPassword(cleanPayload, cleanPassword)
	cleanPayload[FLD_USERPASSWORD_TIMESTAMP] = time.Now().Unix()
	if nil != instance && instance.isAuthDbEnabled() {
		secretKey := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
		if len(secretKey) > 0 {
			var payload string
			payload, err := storage.EncryptPayload(secretKey, cleanPayload)
			if nil == err {
				userId := storage.BuildKey(cleanUsername, cleanPassword)
				err = instance.authDb.AuthRegister(userId, payload)
				isAlreadyRegistred := nil != err && err.Error() == "entity_already_registered"
				if nil != err {
					response.Error = err.Error()
				}
				if nil == err || isAlreadyRegistred {
					response.ItemId = userId
					instance.setPayload(response.ItemPayload, cleanPayload) // user data

					//-- ready to get tokens --//
					var accessToken, refreshToken, confirmToken, tokenError string
					accessToken, refreshToken, confirmToken, tokenError = instance.generateAllTokens(userId, cleanPayload)
					if len(tokenError) > 0 {
						response.Error = tokenError
					} else if isAlreadyRegistred {
						response.AccessToken = accessToken
						response.RefreshToken = refreshToken
						response.ConfirmToken = confirmToken
					} else {
						response.AccessToken = accessToken
						response.RefreshToken = refreshToken
						response.ConfirmToken = confirmToken
					}
				}
			} else {
				response.Error = err.Error()
			}
		} else {
			response.Error = ErrorMissingSecureKey.Error()
		}
	} else {
		response.Error = ErrorUnauthorized.Error()
	}
	return response
}

func (instance *Auth0) AuthConfirm(confirmToken string) *Auth0Response {
	response := new(Auth0Response)
	response.ItemId = ""
	response.ItemPayload = nil

	//-- validate ownerToken --//
	token, err := instance.parseToken(confirmToken)
	if nil != err {
		response.Error = err.Error()
	} else {
		//-- confirm and update user --//
		claims := token.GetMapClaims()
		userId := qbc.Reflect.GetString(claims, FLD_USERID)
		if payload, b := qbc.Reflect.Get(claims, FLD_PAYLOAD).(map[string]interface{}); b {
			payload[FLD_CONFIRMED] = true
			_, err = instance.AuthUpdate(userId, payload)
			if nil != err {
				response.Error = err.Error()
			} else {
				//-- login user to return auth info --//
				username := instance.getUsername(payload) // qbc.Reflect.GetString(payload, FLD_USERNAME)
				password := instance.getPassword(payload)

				response = instance.AuthSignIn(username, password)
			}
		} else {
			response.Error = ErrorMalformedAccountData.Error()
		}
	}
	return response
}

// AuthGetCredentials
// Expose user credentials in code may be Unsecure!!! Use this method only for debugging
func (instance *Auth0) AuthGetCredentials(currentId string) (username, password string, err error) {
	if nil != instance && instance.isAuthDbEnabled() {
		authSecret := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
		if len(authSecret) > 0 {
			// get existing
			var currentPayload map[string]interface{}
			_, currentPayload, err = instance.getDecPayloadById(authSecret, currentId)
			if nil == err {
				username = instance.getUsername(currentPayload)
				password = instance.getPassword(currentPayload)
			}
		} else {
			err = qbc.Errors.Prefix(qbc.ErrorSystem, "Unable to retrieve auth secret: ")
		}
	}
	return
}

// AuthChangeLogin change username or password for login.
// New entity is created
func (instance *Auth0) AuthChangeLogin(currentId, newCleanUsername, newCleanPassword string) *Auth0Response {
	response := new(Auth0Response)
	response.ItemId = ""
	response.ItemPayload = nil
	response.Error = ErrorUnauthorized.Error()
	if nil != instance && instance.isAuthDbEnabled() {
		authSecret := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
		if len(authSecret) > 0 {
			// get existing
			_, currentPayload, err := instance.getDecPayloadById(authSecret, currentId)
			if nil != err {
				response.Error = err.Error()
			}
			if nil == currentPayload {
				currentPayload = map[string]interface{}{}
			}
			newId := currentId
			if len(newCleanUsername) > 0 && len(newCleanPassword) > 0 {
				newId = storage.BuildKey(newCleanUsername, newCleanPassword)
				// currentPayload[FLD_USERNAME] = newCleanUsername
				// currentPayload[FLD_USERPASSWORD] = newCleanPassword
				instance.setUsername(currentPayload, newCleanUsername)
				instance.setPassword(currentPayload, newCleanPassword)
				currentPayload[FLD_USERPASSWORD_TIMESTAMP] = time.Now().Unix()
			}

			payload, err := storage.EncryptPayload(authSecret, currentPayload)
			if nil == err {
				err = instance.authDb.AuthOverwrite(newId, payload)
			}
			// remove old one
			if newId != currentId {
				err = instance.authDb.AuthRemove(currentId)
			}
			if nil != err {
				response.Error = err.Error()
			} else {
				response.ItemId = newId
				response.ItemPayload = currentPayload
				//-- ready to get tokens --//
				response.AccessToken, response.RefreshToken, response.ConfirmToken, response.Error = instance.generateAllTokens(newId, currentPayload)
			}
		}
	}
	return response
}

// AuthUpdate update existing item
// Parameter currentId is required.
// return itemId and error
func (instance *Auth0) AuthUpdate(currentId string, cleanPayload map[string]interface{}) (string, error) {
	if nil != instance && instance.isAuthDbEnabled() {
		authSecret := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
		if len(authSecret) > 0 {
			// get existing
			_, currentPayload, err := instance.getDecPayloadById(authSecret, currentId)
			if nil != err {
				return "", err
			}
			if nil == currentPayload {
				currentPayload = map[string]interface{}{}
			}
			if nil == cleanPayload {
				cleanPayload = map[string]interface{}{}
			}
			for k, v := range currentPayload {
				// add existing payload if not overwritten
				if _, b := cleanPayload[k]; !b {
					cleanPayload[k] = v
				}
			}

			payload, err := storage.EncryptPayload(authSecret, cleanPayload)
			if nil == err {
				err = instance.authDb.AuthOverwrite(currentId, payload)
			}
			return currentId, err
		}
	}
	return "", ErrorUnauthorized
}

// AuthRemove remove entity and associated tokens
func (instance *Auth0) AuthRemove(accessToken string) error {
	if nil != instance && instance.isAuthDbEnabled() {
		// validate ownerToken
		token, err := instance.parseToken(accessToken)
		if nil != err {
			return ErrorUnauthorized
		}
		if !token.Valid {
			return ErrorUnauthorized
		}

		claims := token.GetMapClaims()
		if cacheId, b := claims[CACHE_KEY].(string); b {
			_ = instance.removeTokenFromDatabase(cacheId)
		} else {
			return ErrorUnauthorized
		}

		if userId, b := claims[FLD_USERID].(string); b {
			err = instance.authDb.AuthRemove(userId)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

// AuthRemoveByUserId remove entity
func (instance *Auth0) AuthRemoveByUserId(userId string) error {
	if nil != instance && instance.isAuthDbEnabled() {
		err := instance.authDb.AuthRemove(userId)
		if nil != err {
			return err
		}
	}
	return nil
}

// AuthGrantDelegation create a delegation token that impersonate owner and can be used
func (instance *Auth0) AuthGrantDelegation(ownerToken string) *Auth0Response {
	response := new(Auth0Response)
	response.ItemId = ""
	response.ItemPayload = nil

	// validate ownerToken
	token, err := instance.parseToken(ownerToken)
	if nil != err {
		response.Error = ErrorUnauthorized.Error()
		return response
	}
	if !token.Valid {
		response.Error = ErrorUnauthorized.Error()
		return response
	}

	// get claims from original owner Token
	claims := token.GetMapClaims()
	userId := claims[FLD_USERID].(string)
	if len(userId) == 0 {
		response.Error = ErrorUnauthorized.Error()
		return response
	}

	// refresh payload
	encPayload, err := instance.authDb.AuthGet(userId)
	if nil != err {
		response.Error = ErrorUnauthorized.Error()
		return response
	}

	// decode payload
	secretKey := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
	payload, err := instance.decodePayload(secretKey, encPayload)
	if nil != err {
		response.Error = ErrorUnauthorized.Error()
		return response
	}

	// update payload
	payload[FLD_USERID] = userId
	payload["impersonate"] = true

	// assign response
	response.ItemId = userId
	response.ItemPayload = payload

	//-- ready to get tokens --//
	_, dToken, err := instance.generateToken(TDelegate, userId, payload)
	//accessToken, _, err := instance.generateTokens(true, userId, payload)
	if nil == err {
		response.AccessToken = dToken
		response.RefreshToken = "" // no need of refresh token
	} else {
		response.Error = err.Error()
	}

	return response
}

func (instance *Auth0) AuthRevokeDelegation(delegationToken string) error {

	// validate ownerToken
	token, err := instance.parseToken(delegationToken)
	if nil != err {
		return ErrorUnauthorized
	}
	if !token.Valid {
		return ErrorUnauthorized
	}

	claims := token.GetMapClaims()
	if cacheId, b := claims[CACHE_KEY].(string); b {
		_ = instance.removeTokenFromDatabase(cacheId)
	} else {
		return ErrorUnauthorized
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	T o k e n s    V a l i d a t i o n
//----------------------------------------------------------------------------------------------------------------------

func (instance *Auth0) TokenValidate(stringToken string) (bool, error) {
	// parse the token
	token, err := instance.parseToken(stringToken)
	if nil != err {
		return false, err
	}
	if !token.Valid {
		return false, ErrorUnauthorized
	}

	// get token from cache for a double check
	if instance.isCacheEnabled() {
		if key, b := token.GetMapClaims()[CACHE_KEY].(string); b {
			dbToken, err := instance.getTokenFromDatabase(key)
			if nil != err {
				return false, err
			}
			if !dbToken.Valid {
				return false, ErrorUnauthorized
			}
		} else {
			return false, ErrorUnauthorized
		}
	}

	return true, nil
}

func (instance *Auth0) TokenParse(stringToken string) (map[string]interface{}, error) {
	// parse the token
	token, err := instance.parseToken(stringToken)
	if nil != err {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrorUnauthorized
	}

	claims := token.GetMapClaims()
	response := map[string]interface{}{}
	instance.setPayload(response, claims)

	return response, nil
}

func (instance *Auth0) TokenClaims(stringToken string) (claims map[string]interface{}, err error) {
	claims, err = instance.parseClaims(stringToken)
	return
}

func (instance *Auth0) TokenClaimsNoValidate(stringToken string) (claims map[string]interface{}) {
	claims = instance.parseClaimsNoValidate(stringToken)
	return
}

func (instance *Auth0) TokenRefresh(stringRefreshToken string) *Auth0Response {
	response := emptyResponse()
	// parse the token
	refreshToken, err := instance.parseToken(stringRefreshToken)
	if nil != err {
		response.Error = err.Error() // Token is not a refresh token
	} else {
		if !refreshToken.Valid {
			response.Error = ErrorUnauthorized.Error()
		} else {
			// get payload from token
			if payload, b := refreshToken.GetMapClaims()["payload"].(map[string]interface{}); b {
				if refreshUuid, b := payload["refresh_uuid"].(string); b {
					// now we have id of token to refresh
					refreshed, claims, err := instance.refreshDBToken(refreshUuid)
					if nil != err {
						response.Error = err.Error()
					} else {
						if nil != claims {
							if userId, b := claims[FLD_USERID].(string); b {
								response.ItemId = userId // userId
							}
							if payload, b := claims[FLD_PAYLOAD].(map[string]interface{}); b {
								// response.ItemPayload = payload // user data
								instance.setPayload(response.ItemPayload, payload)
							}
						}
						response.RefreshToken = stringRefreshToken
						response.AccessToken = refreshed
					}
				} else {
					response.Error = ErrorUnauthorized.Error()
				}
			} else {
				response.Error = ErrorUnauthorized.Error()
			}
		}
	}
	return response
}

// TokenRefreshAccess utility method that do not check on db for token existance
func (instance *Auth0) TokenRefreshAccess(stringAccessToken, stringRefreshToken string) *Auth0Response {
	response := emptyResponse()
	// parse the token
	refreshToken, err := instance.parseToken(stringRefreshToken)
	if nil != err {
		response.Error = err.Error() // Token is not a refresh token
	} else {
		if !refreshToken.Valid {
			response.Error = ErrorUnauthorized.Error()
		} else {
			// get payload from token
			if payload, b := refreshToken.GetMapClaims()["payload"].(map[string]interface{}); b {
				if _, b := payload["refresh_uuid"].(string); b {
					// now we have id of token to refresh
					refreshed, claims, err := instance.refreshAccessToken(stringAccessToken)
					if nil != err {
						response.Error = err.Error()
					} else {
						if nil != claims {
							if userId, b := claims[FLD_USERID].(string); b {
								response.ItemId = userId // userId
							}
							if payload, b := claims[FLD_PAYLOAD].(map[string]interface{}); b {
								// response.ItemPayload = payload // user data
								instance.setPayload(response.ItemPayload, payload)
							}
						}
						response.RefreshToken = stringRefreshToken
						response.AccessToken = refreshed
					}
				} else {
					response.Error = ErrorUnauthorized.Error()
				}
			} else {
				response.Error = ErrorUnauthorized.Error()
			}
		}
	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Auth0) getDurations() map[string][]time.Duration {
	// contains token duration and cache duration
	return map[string][]time.Duration{
		AccessSecretName:  {instance.AccessTokenDuration, instance.AccessCacheDuration},
		RefreshSecretName: {instance.RefreshTokenDuration, instance.RefreshTokenDuration},
	}
}

func (instance *Auth0) isCacheEnabled() bool {
	return nil != instance && nil != instance.cacheDb && instance.cacheDb.Enabled()
}

func (instance *Auth0) isAuthDbEnabled() bool {
	return nil != instance && nil != instance.authDb && instance.authDb.Enabled()
}

func (instance *Auth0) getSecret(key string) string {
	if nil != instance && nil != instance.secrets {
		return instance.secrets.Get(key)
	}
	return ""
}

func (instance *Auth0) getEncPayload(cleanUsername, cleanPassword string) (id string, payload string, err error) {
	id = storage.BuildKey(cleanUsername, cleanPassword)
	return instance.getEncPayloadById(id)
}

func (instance *Auth0) getEncPayloadById(dbKey string) (id string, payload string, err error) {
	id = dbKey
	payload, err = instance.authDb.AuthGet(id)
	return
}

func (instance *Auth0) getDecPayload(key []byte, cleanUsername, cleanPassword string) (string, map[string]interface{}, error) {
	id, encPayload, err := instance.getEncPayload(cleanUsername, cleanPassword)
	if nil != err {
		return "", nil, err
	}
	payload, err := instance.decodePayload(key, encPayload)
	if nil != err {
		return "", nil, err
	}
	return id, payload, nil
}

func (instance *Auth0) getDecPayloadById(key []byte, dbKey string) (string, map[string]interface{}, error) {
	id, encPayload, err := instance.getEncPayloadById(dbKey)
	if nil != err {
		return "", nil, err
	}
	payload, err := instance.decodePayload(key, encPayload)
	if nil != err {
		return "", nil, err
	}
	return id, payload, nil
}

func (instance *Auth0) decodePayload(key []byte, encPayload string) (map[string]interface{}, error) {
	if len(encPayload) == 0 {
		return map[string]interface{}{}, nil
	}
	payload, err := storage.DecryptPayload(key, encPayload)
	if nil != err {
		return nil, err
	}
	return payload, nil
}

func (instance *Auth0) generateAllTokens(userId string, cleanPayload map[string]interface{}) (accessToken, refreshToken, confirmToken, retErr string) {
	//-- ready to get tokens --//
	// ACCESS token
	adbKey, aToken, err := instance.generateToken(TAccess, userId, cleanPayload)
	if nil != err {
		retErr = err.Error()
		return
	} else {
		accessToken = aToken
		// REFRESH token
		_, rToken, err := instance.generateToken(TRefresh, userId, map[string]interface{}{"refresh_uuid": adbKey})
		if nil != err {
			retErr = err.Error()
			return
		} else {
			refreshToken = rToken
			// CONFIRM token
			_, cToken, err := instance.generateToken(TConfirm, userId, cleanPayload)
			if nil != err {
				retErr = err.Error()
				return
			} else {
				confirmToken = cToken
			}
		}
	}
	return
}

func (instance *Auth0) generateToken(t int, userId string, payload map[string]interface{}) (string, string, error) {
	var cacheDuration, tokenDuration time.Duration
	var secretName string
	var isDelegation bool
	switch t {
	case TAccess:
		cacheDuration = instance.AccessCacheDuration
		tokenDuration = instance.AccessTokenDuration
		secretName = AccessSecretName
		isDelegation = false
	case TRefresh:
		cacheDuration = instance.RefreshTokenDuration
		tokenDuration = instance.RefreshTokenDuration
		secretName = RefreshSecretName
		isDelegation = false
		// payload := map[string]interface{}{"refresh_uuid": accessDbKey}
	case TConfirm:
		cacheDuration = instance.ConfirmTokenDuration
		tokenDuration = instance.ConfirmTokenDuration
		secretName = AccessSecretName
		isDelegation = false
	case TDelegate:
		cacheDuration = instance.DelegateTokenDuration
		tokenDuration = instance.DelegateTokenDuration
		secretName = AccessSecretName
		isDelegation = true
	default:
		return "", "", ErrorUnauthorized
	}
	secret := []byte(instance.secrets.GetNotEmpty(secretName))

	// ACCESS
	dbKey, token, err := instance.tokenString(isDelegation, secretName, secret, userId, payload, tokenDuration)
	if nil == err {
		// save tokens in cache
		if instance.isCacheEnabled() {
			err = instance.cacheDb.CacheAdd(dbKey, token, cacheDuration)
			if nil == err {
				return dbKey, token, nil
			}
		}
	}
	return "", "", err // dbKey, token, err
}

func (instance *Auth0) refreshDBToken(key string) (string, map[string]interface{}, error) {
	token, err := instance.getTokenFromDatabase(key)
	if nil != token && (nil == err || err != ErrorUnauthorized) {
		claims := token.GetMapClaims()
		if len(claims) > 0 {
			secretType := claims[FLD_SECRET_TYPE].(string)
			if len(secretType) > 0 {
				secretKey := []byte(instance.Secrets().Get(secretType))
				durations := instance.getDurations()
				if d, b := durations[secretType]; b {
					durationToken := d[0]
					durationCache := d[1]
					claims[FLD_EXP] = time.Now().Add(durationToken).Unix() // refresh duration
					newToken := jwt.NewWithClaims(signing.SigningMethodHS256, claims)
					signed, err := newToken.SignedString(secretKey)
					if nil == err {
						err := instance.cacheDb.CacheAdd(key, signed, durationCache)
						if nil == err {
							return signed, claims, nil
						} else {
							// error
							return "", nil, err
						}
					} else {
						// error
						return "", nil, err
					}
				}
			}
		}
	} else {
		// error
		return "", nil, err
	}
	// error
	return "", nil, ErrorUnauthorized
}

func (instance *Auth0) refreshAccessToken(accessTokenString string) (string, map[string]interface{}, error) {
	accessToken, err := instance.parseToken(accessTokenString)
	if nil != accessToken && (nil == err || err != ErrorUnauthorized) {
		claims := accessToken.GetMapClaims()
		if len(claims) > 0 {
			secretType := claims[FLD_SECRET_TYPE].(string)
			if len(secretType) > 0 {
				secretKey := []byte(instance.Secrets().Get(secretType))
				durations := instance.getDurations()
				if d, b := durations[secretType]; b {
					durationToken := d[0]
					claims[FLD_EXP] = time.Now().Add(durationToken).Unix() // refresh duration
					newToken := jwt.NewWithClaims(signing.SigningMethodHS256, claims)
					signed, err := newToken.SignedString(secretKey)
					if nil == err {
						return signed, claims, nil
					} else {
						// error
						return "", nil, err
					}
				}
			}
		}
	} else {
		// error
		return "", nil, err
	}
	// error
	return "", nil, ErrorUnauthorized
}

func (instance *Auth0) tokenString(isDelegation bool, secretName string, key []byte, userId string, payload map[string]interface{}, duration time.Duration) (string, string, error) {
	claims := new(Auth0Claims)
	claims.UserId = userId
	claims.Payload = payload
	claims.SecretType = secretName
	claims.Id = instance.tokenUid(isDelegation, secretName, key, userId)

	claims.ExpiresAt = time.Now().Add(duration).Unix()

	token := jwt.NewWithClaims(signing.SigningMethodHS256, claims)
	signed, err := token.SignedString(key)
	return claims.Id, signed, err
}

func (instance *Auth0) tokenUid(isDelegation bool, secretName string, key []byte, userId string) string {
	return qbc.Coding.MD5(fmt.Sprintf("%v-%v-%v-%v", isDelegation, secretName, key, userId))
}

func (instance *Auth0) parseToken(stringToken string) (*elements.Token, error) {
	return jwt.Parse(stringToken, func(token *elements.Token) (interface{}, error) {
		claims := token.GetMapClaims()
		if nil != claims {
			// secretType := token.GetMapClaims()[FLD_SECRET_TYPE].(string)
			secretType := qbc.Reflect.GetString(claims, FLD_SECRET_TYPE)
			if len(secretType) > 0 {
				accessKey := []byte(instance.Secrets().Get(secretType))
				if len(accessKey) == 0 {
					return nil, ErrorUnauthorized
				}
				return accessKey, nil
			}
			return nil, ErrorUnauthorized
		}
		return nil, ErrorMissingClaims
	})
}

func (instance *Auth0) parseClaims(stringToken string) (claims map[string]interface{}, err error) {
	_, err = jwt.Parse(stringToken, func(token *elements.Token) (interface{}, error) {
		claims = token.GetMapClaims()
		accessKey := []byte(instance.Secrets().Get("access"))
		return accessKey, nil
	})
	return
}

func (instance *Auth0) parseClaimsNoValidate(stringToken string) (claims map[string]interface{}) {
	token, _ := jwt.Parse(stringToken, func(token *elements.Token) (interface{}, error) {
		return nil, nil
	})
	if nil != token {
		claims = token.GetMapClaims()
	}
	return
}

func (instance *Auth0) getTokenFromDatabase(key string) (*elements.Token, error) {
	if len(key) > 0 {
		item, _ := instance.cacheDb.CacheGet(key)
		if len(item) > 0 {
			token, err := instance.parseToken(item)
			if nil != err {
				return token, err
			}
			return token, nil
		}
	}

	return nil, ErrorUnauthorized
}

func (instance *Auth0) removeTokenFromDatabase(key string) error {
	if len(key) > 0 {
		err := instance.cacheDb.CacheRemove(key)
		return err
	}
	return nil
}

func (instance *Auth0) setPayload(target map[string]interface{}, source map[string]interface{}) {
	if nil != target {
		for k, v := range source {
			if m, b := v.(map[string]interface{}); b {
				instance.setPayload(target, m)
			} else {
				if qbc.Arrays.IndexOf(k, PROTECTED_FIELDS) == -1 {
					target[k] = v
				}
			}
		}
	}
}

func (instance *Auth0) debug(name, t string) {
	token, err := instance.parseToken(t)
	if nil != err {
		fmt.Println(name, "DEBUG: PARSE ERROR", err)
	} else {
		if !token.Valid {
			fmt.Println(name, "DEBUG: INVALID TOKEN", err)
		}

		claims := token.GetMapClaims()
		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		fmt.Println(name, "DEBUG:\t", exp, qbc.JSON.Stringify(claims))
	}
}

func (instance *Auth0) getPassword(payload map[string]interface{}) string {
	return instance.decodeField(qbc.Reflect.GetString(payload, FLD_USERPASSWORD))
}

func (instance *Auth0) setPassword(payload map[string]interface{}, raw string) {
	psw := instance.encodeField(raw)
	payload[FLD_USERPASSWORD] = psw
}

func (instance *Auth0) getUsername(payload map[string]interface{}) string {
	return instance.decodeField(qbc.Reflect.GetString(payload, FLD_USERNAME))
}

func (instance *Auth0) setUsername(payload map[string]interface{}, raw string) {
	psw := instance.encodeField(raw)
	payload[FLD_USERNAME] = psw
}

func (instance *Auth0) decodeField(raw string) string {
	secretKey := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
	s, err := qbc.Coding.DecryptTextWithPrefix(raw, secretKey)
	if nil != err {
		return raw
	}
	return s
}

func (instance *Auth0) encodeField(raw string) string {
	secretKey := []byte(instance.secrets.GetNotEmpty(AuthSecretName))
	s, err := qbc.Coding.EncryptTextWithPrefix(raw, secretKey)
	if nil != err {
		return raw
	}
	return s
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func IsPasswordExpired(payload map[string]interface{}) bool {
	return isPasswordExpired(payload, 90)
}

func GetRefreshTokenDuration() time.Duration {
	return refreshTD
}

func SetRefreshTokenDuration(value time.Duration) {
	refreshTD = value
}

func GetAccessTokenDuration() time.Duration {
	return accessTD
}

func SetAccessTokenDuration(value time.Duration) {
	accessTD = value
}

func GetAccessCacheDuration() time.Duration {
	return accessCD
}

func SetAccessCacheDuration(value time.Duration) {
	accessCD = value
}

func isPasswordExpired(payload map[string]interface{}, durationDays int) bool {
	if durationDays > 0 {
		lastChangedTimestamp := int64(qbc.Reflect.GetInt(payload, FLD_USERPASSWORD_TIMESTAMP))
		if lastChangedTimestamp > 0 {
			lastChangedDate := time.Unix(lastChangedTimestamp, 0)
			diffTime := time.Now().Sub(lastChangedDate)
			diffDays := int(diffTime.Hours() / 24)
			return diffDays > durationDays
		}
	}
	return false
}

func emptyResponse() *Auth0Response {
	response := new(Auth0Response)
	response.ItemId = ""
	response.ItemPayload = map[string]interface{}{}

	return response
}
