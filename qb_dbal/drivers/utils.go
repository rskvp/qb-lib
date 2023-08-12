package drivers

import (
	"strings"

	qbc "github.com/rskvp/qb-core"
)

//----------------------------------------------------------------------------------------------------------------------
//	u t i l s
//----------------------------------------------------------------------------------------------------------------------

// QueryGetParamNames return unique param names
func QueryGetParamNames(query string) []string {
	response := make([]string, 0)
	query = strings.ReplaceAll(query, ";", " ;")
	query += " "
	params := qbc.Regex.TextBetweenStrings(query, "@", " ")
	for _, param := range params {
		if qbc.Arrays.IndexOf(param, response) == -1 {
			response = append(response, param)
		}
	}
	return response
}

func QuerySelectParams(query string, allParams map[string]interface{}) map[string]interface{} {
	names := QueryGetParamNames(query)
	params := map[string]interface{}{}
	for _, v := range names {
		params[v] = allParams[v]
	}
	return params
}

func IsRecordNotFoundError(err error) bool {
	return nil != err && err.Error() == "record not found"
}

func mergeParams(query string, bindVars map[string]interface{}) (string, []interface{}) {
	if nil == bindVars || len(bindVars) == 0 {
		return query, []interface{}{}
	}
	paramNames := queryParamNames(query)
	if len(paramNames) > 0 {
		for _, name := range paramNames {
			query = strings.ReplaceAll(query, "@"+name, qbc.Reflect.GetString(bindVars, name))
		}
	}
	return query, []interface{}{}
}

func queryParamNames(query string) []string {
	query = query + " "
	query = strings.ReplaceAll(query, ";", " ;")
	return qbc.Regex.TextBetweenStrings(query, "@", " ")
}
