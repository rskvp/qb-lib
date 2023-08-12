package semantic_search

import (
	"fmt"
	"sort"
	"strings"

	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_dbal/commons"
	"github.com/rskvp/qb-lib/qb_dbal/drivers"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	COLLECTION = "elastic_search_indexed"

	FLD_DBKEY  = "_key"
	FLD_KEY    = "key"
	FLD_GROUP  = "group"
	FLD_ENTITY = "entity"
	FLD_TAGS   = "tags"

	ANALYZER_LOWER = "semantic_lowercase"
)

type SemanticEngineData struct {
	Score  int                    `json:"score"`
	Key    string                 `json:"key"`
	Group  string                 `json:"group"`
	Entity map[string]interface{} `json:"entity"`
}

//----------------------------------------------------------------------------------------------------------------------
//	SemanticEngine
//----------------------------------------------------------------------------------------------------------------------

type SemanticEngine struct {
	config     *commons.SemanticConfig
	dbInternal drivers.IDatabase
	dbExternal drivers.IDatabase
}

func NewSemanticEngine(config *commons.SemanticConfig) (*SemanticEngine, error) {
	instance := new(SemanticEngine)
	instance.config = config

	// internal
	if nil != config.DbInternal && config.DbInternal.IsValid() {
		db, err := drivers.OpenDatabase(config.DbInternal.Driver, config.DbInternal.Dsn)
		if nil != err {
			return nil, err
		}
		instance.dbInternal = db
		instance.initInternal()
	} else {
		// internal db is required but not properly configured
		return nil, commons.ErrorMismatchConfiguration
	}

	// external with data indexed by elastic search
	if nil != config.DbExternal && config.DbExternal.IsValid() {
		db, err := drivers.OpenDatabase(config.DbExternal.Driver, config.DbExternal.Dsn)
		if nil != err {
			return nil, err
		}
		instance.dbExternal = db
	}

	return instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *SemanticEngine) IsCaseSensitive() bool {
	if nil != instance && nil != instance.config {
		return instance.config.CaseSensitive
	}
	return false
}

func (instance *SemanticEngine) Put(group, key, text string) error {
	if nil != instance && nil != instance.dbInternal {
		_, err := instance.dbInternal.Upsert(COLLECTION, instance.prepare(group, key, text))
		return err
	}
	return commons.ErrorEngineNotReady
}

func (instance *SemanticEngine) Get(group, text string, offset, count int) ([]*SemanticEngineData, error) {
	if nil != instance && nil != instance.dbInternal {
		// build filter
		var err error
		if len(text) > 0 {
			bindVars := map[string]interface{}{"@collection": COLLECTION}
			query, tokens := buildQuery(group, text, offset, count, bindVars, instance.IsCaseSensitive())

			data, err := instance.dbInternal.ExecNative(query, bindVars)
			if nil != err {
				return nil, err
			}
			if v, b := data.([]interface{}); b {
				response := make([]*SemanticEngineData, 0)
				for _, item := range v {
					if entity, b := item.(map[string]interface{}); b {
						data := new(SemanticEngineData)
						data.Entity = entity
						data.Key = entity[FLD_KEY].(string)
						data.Group = entity[FLD_GROUP].(string)
						data.Score = getScore(tokens, qbc.Reflect.GetArrayOfString(entity, FLD_TAGS), instance.IsCaseSensitive())

						// need recover an entity from external?
						if nil != instance.dbExternal {
							item, err := instance.dbExternal.Get(data.Group, data.Key)
							if nil != err {
								return nil, err
							}
							if nil != item {
								data.Entity = item
								response = append(response, data)
							} else {
								// ENTITY NOT FOUND: remove indexed key
								_ = instance.dbInternal.Remove(COLLECTION, entity[FLD_DBKEY].(string))
							}
						} else {
							response = append(response, data)
						}
					}
				}
				// sort by score
				sort.Slice(response, func(i, j int) bool {
					return response[i].Score > response[j].Score
				})
				return response, nil
			}
		}
		return nil, err
	}
	return nil, commons.ErrorEngineNotReady
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *SemanticEngine) initInternal() {
	if nil != instance.dbInternal {
		_, _ = instance.dbInternal.EnsureIndex(COLLECTION, "", []string{FLD_GROUP, FLD_TAGS}, false)
		if arango, b := instance.dbInternal.(*drivers.DriverArango); b {
			// https://www.arangodb.com/docs/stable/arangosearch-case-sensitivity-and-diacritics.html
			// analyzers.save("norm_en", "norm", { locale: "en.utf-8", accent: false, case: "lower" }, []);
			_ = arango.EnsureAnalyzer(map[string]interface{}{
				"name": ANALYZER_LOWER,
				"type": "norm",
				"properties": map[string]interface{}{
					"locale": "en.utf-8",
					"accent": false,
					"case":   "lower",
				},
			})
		}
	}
}

func (instance *SemanticEngine) prepare(group, key, text string) map[string]interface{} {
	response := map[string]interface{}{
		FLD_DBKEY: qbc.Coding.MD5(group + key),
		FLD_GROUP: group,
		FLD_KEY:   key,
		FLD_TAGS:  ToKeywords(text),
	}

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func ToKeywords(text string) []string {
	response := make([]string, 0)
	if len(text) > 0 {
		tokens := qbc.Strings.Split(text, " ,.!?")
		for _, token := range tokens {
			if len(token) > 2 {
				if len(token) > 4 {
					for i := 0; i < len(token)-2; i++ {
						response = append(response, token[:len(token)-i])
					}
				} else {
					response = append(response, token)
				}
			}
		}
	}
	return response
}

func buildQuery(group, text string, offset, count int, bindVars map[string]interface{}, caseSensitive bool) (string, []string) {
	if !caseSensitive {
		text = strings.ToLower(text)
	}
	tokens := qbc.Strings.Split(text, " ,.!?") // ToKeywords(text)
	bindVars["tokens"] = tokens
	query := "FOR doc IN @@collection"
	if len(group) > 0 {
		bindVars["group"] = group
		query += fmt.Sprintf(" FILTER doc.%s==@group AND", FLD_GROUP)
	} else {
		query += " FILTER"
	}
	if !caseSensitive {
		query += fmt.Sprintf(" @tokens ANY IN FLATTEN(TOKENS(doc.%s,  \"%s\"))", FLD_TAGS, ANALYZER_LOWER)
	} else {
		query += fmt.Sprintf(" doc.%s[*] ANY IN @tokens", FLD_TAGS)
	}

	if count > 0 {
		query += fmt.Sprintf(" LIMIT %v, %v", offset, count)
	}
	query += " RETURN doc"

	return query, tokens
}

func getScore(keywords, tags []string, caseSensitive bool) int {
	count := 0
	for _, k := range keywords {
		for _, t := range tags {
			if caseSensitive {
				if k == t {
					count++
				}
			} else {
				if strings.ToLower(k) == strings.ToLower(t) {
					count++
				}
			}
		}
	}
	return count
}
