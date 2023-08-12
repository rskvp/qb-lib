package showcase_search

import (
	"fmt"
	"strings"
	"time"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-core/qb_utils"
	"github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	COLL_SHOWCASE_CONTENTS       = "dbal_showcase_contents"
	COLL_SHOWCASE_SESSIONS       = "dbal_showcase_sessions"
	COLL_SHOWCASE_SESSIONS_CACHE = "dbal_showcase_sessions_cache" // params used in session

	FLD_DBKEY               = "_key"
	FLD_SESSION_ID          = "session_id" // each user has a session id
	FLD_TIMESTAMP           = "timestamp"
	FLD_PAYLOAD             = "payload"
	FLD_KEY                 = "key"
	FLD_CATEGORY            = "category"
	FLD_CATEGORY_WEIGHT_IN  = "category_weight_in"
	FLD_CATEGORY_WEIGHT_OUT = "category_weight_out"
)

//----------------------------------------------------------------------------------------------------------------------
//	SemanticEngine
//----------------------------------------------------------------------------------------------------------------------

type ShowcaseEngine struct {
	root             string // used to store internal files
	config           *commons.SemanticConfigDb
	db               drivers.IDatabase
	autoResetSession bool

	initialized bool
	categories  *ShowcaseCategories
}

func NewShowcaseEngine(config *commons.SemanticConfigDb) (*ShowcaseEngine, error) {
	instance := new(ShowcaseEngine)
	instance.config = config
	instance.root = qbc.Paths.Concat(qbc.Paths.GetWorkspacePath(), "dbal", "showcase")

	// internal
	if nil != config && config.IsValid() {
		db, err := drivers.OpenDatabase(config.Driver, config.Dsn)
		if nil != err {
			return nil, err
		}
		instance.db = db
		instance.initInternal()
	} else {
		// internal db is required but not properly configured
		return nil, commons.ErrorMismatchConfiguration
	}

	instance.autoResetSession = true // never ending contents, when finished the system reset cache

	return instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// SetRoot set custom working root (default is "./_workspace/dbal/showcase/")
func (instance *ShowcaseEngine) SetRoot(value string) {
	if nil != instance {
		instance.Reset()
		instance.root = value
	}
}

// Reset reset Categories
func (instance *ShowcaseEngine) Reset() {
	if nil != instance {
		instance.initialized = false
		instance.categories.Clear()
		instance.root = qbc.Paths.Concat(qbc.Paths.GetWorkspacePath(), "dbal", "showcase")
	}
}

// SetCategoryWeight set a custom global Category Weight
func (instance *ShowcaseEngine) SetCategoryWeight(category string, inTime bool, weight int) map[string]ShowcaseCategoryWeight {
	if nil != instance {
		instance.init()
		instance.categories.SetWeight(category, inTime, weight)
	}
	return nil
}

// SetAutoResetSession set custom autoResetSession param
func (instance *ShowcaseEngine) SetAutoResetSession(value bool) {
	if nil != instance {
		instance.autoResetSession = value
	}
}

func (instance *ShowcaseEngine) Categories() *ShowcaseCategories {
	if nil != instance {
		instance.init()
		return instance.categories
	}
	return nil
}

// Put put data into storage
func (instance *ShowcaseEngine) Put(payload interface{}, timestamp int64, category string) (interface{}, error) {
	if nil != instance {
		instance.init()
		if timestamp == 0 {
			timestamp = time.Now().Unix()
		}
		key := qbc.Coding.MD5(qbc.Convert.ToString(payload))
		entity := map[string]interface{}{
			FLD_DBKEY:               key,
			FLD_TIMESTAMP:           timestamp,
			FLD_PAYLOAD:             payload,
			FLD_CATEGORY:            category,
			FLD_CATEGORY_WEIGHT_IN:  instance.categories.Get(category).WeightInDate,
			FLD_CATEGORY_WEIGHT_OUT: instance.categories.Get(category).WeightOutDate,
		}
		response, err := instance.db.Upsert(COLL_SHOWCASE_CONTENTS, entity)
		return response, err
	}
	return nil, nil
}

// Get return single entity by key
func (instance *ShowcaseEngine) Get(key string) (interface{}, error) {
	entity, err := instance.db.Get(COLL_SHOWCASE_CONTENTS, key)
	if nil != err {
		return nil, err
	}
	return entity, nil
}

// Update update a payload
func (instance *ShowcaseEngine) Update(key string, payload interface{}) (interface{}, error) {
	entity, err := instance.db.Get(COLL_SHOWCASE_CONTENTS, key)
	if nil != err {
		return nil, err
	}
	entity[FLD_PAYLOAD] = payload
	entity, err = instance.db.Upsert(COLL_SHOWCASE_CONTENTS, entity)
	if nil != err {
		return nil, err
	}
	return entity, nil
}

// Delete remove a payload
func (instance *ShowcaseEngine) Delete(key string) (interface{}, error) {
	entity, err := instance.db.Get(COLL_SHOWCASE_CONTENTS, key)
	if nil != err {
		return nil, err
	}
	err = instance.db.Remove(COLL_SHOWCASE_CONTENTS, key)
	if nil != err {
		return nil, err
	}
	return entity, nil
}

// Query get results as array of db documents with a payload
func (instance *ShowcaseEngine) Query(sessionId string, limit int) []interface{} {
	if nil != instance {
		instance.init()
		return instance.query(sessionId, limit, instance.autoResetSession)
	}
	return nil
}

// SetSessionCategoryWeight change weight just for one user session
func (instance *ShowcaseEngine) SetSessionCategoryWeight(sessionId, category string, inTime bool, weight int) map[string]ShowcaseCategoryWeight {
	session := instance.getSessionCache(sessionId)
	rawCategories := qbc.Reflect.Get(session, "categories")
	if nil != rawCategories {
		var categories map[string]ShowcaseCategoryWeight
		_ = qbc.JSON.Read(qbc.JSON.Stringify(rawCategories), &categories)
		if v, b := categories[category]; b {
			if inTime {
				v.WeightInDate = weight
			} else {
				v.WeightOutDate = weight
			}
			categories[category] = v
		}
		session["categories"] = categories
		_, _ = instance.db.Upsert(COLL_SHOWCASE_SESSIONS_CACHE, session)
	}
	return instance.categories.data
}

// ResetSession reset user session cache
func (instance *ShowcaseEngine) ResetSession(sessionId string) {
	_ = instance.db.Remove(COLL_SHOWCASE_SESSIONS_CACHE, sessionId)
	bindVars := map[string]interface{}{
		"@collection": COLL_SHOWCASE_SESSIONS,
		"session_id":  sessionId,
	}
	query := "FOR doc IN @@collection\n" +
		"  FILTER doc.`session_id` == @session_id\n" +
		"  REMOVE doc IN @@collection"
	_, _ = instance.db.ExecNative(query, bindVars)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseEngine) initInternal() {
	if nil != instance.db {
		// contents
		_, _ = instance.db.EnsureIndex(COLL_SHOWCASE_CONTENTS, "", []string{FLD_CATEGORY, FLD_CATEGORY_WEIGHT_IN, FLD_CATEGORY_WEIGHT_OUT}, false)

		// session
		_, _ = instance.db.EnsureIndex(COLL_SHOWCASE_SESSIONS, "", []string{FLD_SESSION_ID}, false)

		// session categories
		_, _ = instance.db.EnsureIndex(COLL_SHOWCASE_SESSIONS_CACHE, "", []string{FLD_CATEGORY}, false)
	}
}

func (instance *ShowcaseEngine) init() {
	if nil != instance && !instance.initialized {
		instance.initialized = true
		_ = qbc.Paths.Mkdir(instance.root + qb_utils.OS_PATH_SEPARATOR)
		instance.categories = NewShowcaseCategories(instance.root)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	sessions
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseEngine) query(sessionId string, limit int, autoReset bool) []interface{} {
	response := make([]interface{}, 0)
	now := time.Now().Unix()

	for i := 0; i < limit; i++ {
		// get already passed contents
		history := instance.getSessionHistory(sessionId)
		// get contents not in history key array and not in used categories
		usedCategories := instance.getSessionCategoriesAbsolute(sessionId)
		contents := instance.getContents(now, history, usedCategories)
		// add contents to session history
		contentCategories := make([]string, 0)
		for _, content := range contents {
			key := qbc.Reflect.GetString(content, FLD_DBKEY) // key of entity in showcase content collection
			category := qbc.Reflect.GetString(content, FLD_CATEGORY)
			inTime := qbc.Reflect.GetBool(content, "in_time")
			if inTime {
				category += ":in"
			} else {
				category += ":out"
			}
			contentCategories = append(contentCategories, category)
			entity := map[string]interface{}{
				FLD_DBKEY:      qbc.Coding.MD5(sessionId + key),
				FLD_SESSION_ID: sessionId,
				FLD_KEY:        key,
				FLD_TIMESTAMP:  time.Now().Unix(), // added at...
			}
			_, _ = instance.db.Upsert(COLL_SHOWCASE_SESSIONS, entity)
		}
		// update session categories for next usage
		instance.updateSessionCache(sessionId, contentCategories)

		response = append(response, contents...)
	}

	if len(response) == 0 && autoReset {
		instance.ResetSession(sessionId)
		return instance.query(sessionId, limit, false)
	}

	return response
}

func (instance *ShowcaseEngine) getSessionHistory(sessionId string) []interface{} {
	bindVars := map[string]interface{}{
		"@collection": COLL_SHOWCASE_SESSIONS,
		"session_id":  sessionId,
	}
	query := "FOR doc IN @@collection\n" +
		"  FILTER doc.`session_id` == @session_id\n" +
		"  RETURN doc.key"
	resp, _ := instance.db.ExecNative(query, bindVars)
	if v, b := resp.([]interface{}); b {
		return v
	}
	return []interface{}{}
}

//----------------------------------------------------------------------------------------------------------------------
//	sessions persistent params
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseEngine) getSessionCache(sessionId string) map[string]interface{} {
	resp, err := instance.db.Get(COLL_SHOWCASE_SESSIONS_CACHE, sessionId)
	if nil == resp && nil == err {
		entity := map[string]interface{}{
			"_key":       sessionId,
			"absolute":   make([]interface{}, 0),
			"relative":   make([]interface{}, 0),
			"categories": instance.categories.data,
		}
		resp, _ = instance.db.Upsert(COLL_SHOWCASE_SESSIONS_CACHE, entity)
	}
	return resp
}

func (instance *ShowcaseEngine) getSessionCategories(sessionId string) map[string]ShowcaseCategoryWeight {
	session := instance.getSessionCache(sessionId)
	categories := qbc.Reflect.Get(session, "categories")
	if nil != categories {
		var response map[string]ShowcaseCategoryWeight
		_ = qbc.JSON.Read(qbc.JSON.Stringify(categories), &response)
		return response
	}
	return instance.categories.data
}

func (instance *ShowcaseEngine) getSessionCategoriesAbsolute(sessionId string) []interface{} {
	entity := instance.getSessionCache(sessionId)
	if nil != entity {
		return qbc.Reflect.GetArray(entity, "absolute")
	}
	return []interface{}{}
}

func (instance *ShowcaseEngine) updateSessionCache(sessionId string, usedCategories []string) {
	cache := instance.getSessionCache(sessionId)
	absolute := qbc.Reflect.GetArray(cache, "absolute")
	relative := qbc.Reflect.GetArray(cache, "relative")
	sessionCategories := instance.getSessionCategories(sessionId)
	for _, uc := range usedCategories {
		category, inTime := splitCategoryName(uc)
		count := qbc.Arrays.Count(category, relative) // how many times
		cat := GetCategoryWeight(sessionCategories, category)
		limit := cat.WeightInDate
		if !inTime {
			limit = cat.WeightOutDate
		}
		relative = append(relative, category)
		if count+1 >= limit && qbc.Arrays.IndexOf(category, absolute) == -1 {
			absolute = append(absolute, category)
		}
	}

	if len(commons.CATEGORIES) == len(absolute) {
		// RESET
		absolute = make([]interface{}, 0)
		relative = make([]interface{}, 0)
	}
	// update
	entity := map[string]interface{}{}
	entity[FLD_DBKEY] = sessionId
	entity["absolute"] = absolute
	entity["relative"] = relative
	_, _ = instance.db.Upsert(COLL_SHOWCASE_SESSIONS_CACHE, entity)
}

//----------------------------------------------------------------------------------------------------------------------
//	content
//----------------------------------------------------------------------------------------------------------------------

func (instance *ShowcaseEngine) getContents(now int64, excludeKeys []interface{}, excludeCategories []interface{}) []interface{} {
	result := make([]interface{}, 0)
	if nil != instance {
		// collect data
		data, _ := instance.getContentsWithinTime(now, 1, excludeKeys, excludeCategories, true)
		result = append(result, data...)

		if len(result) == 0 {
			data, _ = instance.getContentsWithinTime(now, 1, excludeKeys, excludeCategories, false)
			result = append(result, data...)
		}
	}

	return result
}

func (instance *ShowcaseEngine) getContentsWithinTime(now int64, limit int, excludeKeys []interface{}, excludeCategories []interface{}, inTime bool) ([]interface{}, error) {
	if nil != instance {
		bindVars := map[string]interface{}{
			"@collection":       COLL_SHOWCASE_CONTENTS,
			"limit":             limit,
			"excludeKeys":       excludeKeys,
			"excludeCategories": excludeCategories,
			"now":               now,
		}
		query := "FOR doc IN @@collection\n"
		query += " LET t = doc.timestamp - @now\n"
		query += " FILTER doc._key NOT IN @excludeKeys AND doc.category NOT IN @excludeCategories"
		if inTime {
			query += " AND t>=0\n"
			query += " SORT doc.category_weight_in DESC, doc.category\n"
		} else {
			query += " AND t<0\n"
			query += " SORT doc.category_weight_out DESC, doc.category\n"
		}
		query += " LIMIT @limit\n"
		query += " RETURN {\n" +
			"\"_key\":doc._key,\n" +
			"\"timestamp\":doc.timestamp,\n" +
			"\"payload\":doc.payload,\n" +
			"\"category\":doc.category,\n" +
			"\"in_time\":t>=0,\n" +
			"\"seconds_ago\":ROUND(t*-1),\n" +
			"\"minutes_ago\":ROUND(t/60*-1),\n" +
			"\"hours_ago\":ROUND(t/60/60*-1),\n" +
			"\"days_ago\":ROUND(t/60/60/24*-1)\n" +
			"}"

		resp, err := instance.db.ExecNative(query, bindVars)
		if nil != err {
			return nil, err
		}
		if v, b := resp.([]interface{}); b {
			return v, nil
		}
		return []interface{}{}, nil
	}
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func splitCategoryName(text interface{}) (name string, inTime bool) {
	catTokens := strings.Split(fmt.Sprintf("%v", text), ":")
	switch len(catTokens) {
	case 1:
		name = catTokens[0]
	case 2:
		name = catTokens[0]
		inTime = catTokens[1] == "in"
	}
	return
}
