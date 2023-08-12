package qb_auth0

import qbc "github.com/rskvp/qb-core"

//----------------------------------------------------------------------------------------------------------------------
//	Auth0Response
//----------------------------------------------------------------------------------------------------------------------

type Auth0Response struct {
	Error        string                 `json:"error"`
	ItemId       string                 `json:"item_id"`
	ItemPayload  map[string]interface{} `json:"item_payload"`
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	ConfirmToken string                 `json:"confirm_token"`
}

func (instance *Auth0Response) GoString() string {
	return qbc.JSON.Stringify(instance)
}

func (instance *Auth0Response) String() string {
	return qbc.JSON.Stringify(instance)
}
