package auth0

import (
	"sync"

	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_auth0"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type JsAuth0 struct {
	runtime *goja.Runtime
	object  *goja.Object
	config  *qb_auth0.Auth0Config
	auth0   *qb_auth0.Auth0
	mux     sync.Mutex
}

//----------------------------------------------------------------------------------------------------------------------
//	JsHttpClient
//----------------------------------------------------------------------------------------------------------------------

func WrapAuth0Config(runtime *goja.Runtime, config interface{}) *JsAuth0 {
	instance := new(JsAuth0)
	instance.runtime = runtime
	instance.config = qb_auth0.Auth0ConfigParse(qbc.Convert.ToString(config))

	instance.object = instance.runtime.NewObject()
	instance.export()

	// add closable: all closable objects must be exposed to avoid
	commons.AddClosableObject(instance.runtime, instance.object)

	return instance
}

func WrapAuth0(runtime *goja.Runtime, auth0 *qb_auth0.Auth0) *JsAuth0 {
	instance := new(JsAuth0)
	instance.runtime = runtime
	instance.config = nil
	instance.auth0 = auth0

	instance.object = instance.runtime.NewObject()
	instance.export()

	// add closable: all closable objects must be exposed to avoid
	commons.AddClosableObject(instance.runtime, instance.object)

	return instance
}

func (instance *JsAuth0) Value() goja.Value {
	if nil != instance.object {
		return instance.object
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsAuth0) open(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		err := instance.auth0.Open()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) close(_ goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		err := instance.auth0.Close()
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userSignIn(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		user := commons.GetString(call, 0)
		password := commons.GetString(call, 1)
		if len(user) > 0 && len(password) > 0 {
			response := instance.auth0.AuthSignIn(user, password)
			return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
		} else {
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userUpdate(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var id string
		var payload interface{}
		switch len(call.Arguments) {
		case 2:
			id = commons.GetString(call, 0)
			payload = commons.GetExport(call, 1)
		case 4:
			id = commons.GetString(call, 0)
			payload = commons.GetExport(call, 3)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.auth0.AuthUpdate(id, qbc.Convert.ForceMap(payload))
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(response)
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userChangeLogin(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var id, username, password string
		switch len(call.Arguments) {
		case 3:
			id = commons.GetString(call, 0)
			username = commons.GetString(call, 1)
			password = commons.GetString(call, 2)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response := instance.auth0.AuthChangeLogin(id, username, password)
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userSignUp(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var user, password string
		var payload interface{}
		switch len(call.Arguments) {
		case 3:
			user = commons.GetString(call, 0)
			password = commons.GetString(call, 1)
			payload = commons.GetExport(call, 2)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response := instance.auth0.AuthSignUp(user, password, qbc.Convert.ForceMap(payload))
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userConfirm(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response := instance.auth0.AuthConfirm(token)
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userRemove(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		err := instance.auth0.AuthRemove(token)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userRemoveByAuthId(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var authId string
		switch len(call.Arguments) {
		case 1:
			authId = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		err := instance.auth0.AuthRemoveByUserId(authId)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userGrantDelegation(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response := instance.auth0.AuthGrantDelegation(token)
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

func (instance *JsAuth0) userRevokeDelegation(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		err := instance.auth0.AuthRevokeDelegation(token)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
	}
	return goja.Undefined()
}

func (instance *JsAuth0) tokenValidate(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		success, _ := instance.auth0.TokenValidate(token)
		return instance.runtime.ToValue(success)
	}
	return goja.Undefined()
}

func (instance *JsAuth0) tokenParse(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response, err := instance.auth0.TokenParse(token)
		if nil != err {
			panic(instance.runtime.NewTypeError(err.Error()))
		}
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

func (instance *JsAuth0) tokenRefresh(call goja.FunctionCall) goja.Value {
	if nil != instance {
		instance.init()
		var token string
		switch len(call.Arguments) {
		case 1:
			token = commons.GetString(call, 0)
		default:
			panic(instance.runtime.NewTypeError(commons.ErrorMissingParam))
		}
		response := instance.auth0.TokenRefresh(token)
		return instance.runtime.ToValue(qbc.Convert.ForceMap(response))
	}
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *JsAuth0) init() {
	if nil == instance.auth0 {
		instance.auth0 = qb_auth0.NewAuth0(instance.config)
	}
}

func (instance *JsAuth0) export() {
	o := instance.object

	_ = o.Set("open", instance.open)
	_ = o.Set("close", instance.close)

	// login
	_ = o.Set("userSignIn", instance.userSignIn)
	_ = o.Set("userSignUp", instance.userSignUp)
	_ = o.Set("userUpdate", instance.userUpdate)
	_ = o.Set("userChangeLogin", instance.userChangeLogin)
	_ = o.Set("userConfirm", instance.userConfirm)
	_ = o.Set("userRemove", instance.userRemove)
	_ = o.Set("userRemoveByAuthId", instance.userRemoveByAuthId)
	_ = o.Set("userGrantDelegation", instance.userGrantDelegation)
	_ = o.Set("userRevokeDelegation", instance.userRevokeDelegation)

	_ = o.Set("tokenValidate", instance.tokenValidate)
	_ = o.Set("tokenRefresh", instance.tokenRefresh)
	_ = o.Set("tokenParse", instance.tokenParse)
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------
