package bolt

import (
	"encoding/json"
	"strings"

	qbc "github.com/rskvp/qb-core"
)

const (
	ComparatorEqual        = "=="
	ComparatorNotEqual     = "!="
	ComparatorGreater      = ">"
	ComparatorLower        = "<"
	ComparatorLowerEqual   = "<="
	ComparatorGreaterEqual = ">="

	OperatorAnd = "&&"
	OperatorOr  = "||"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type BoltQuery struct {
	Conditions []*BoltQueryConditionGroup `json:"conditions"`
}

type BoltQueryConditionGroup struct {
	Operator string                `json:"operator"`
	Filters  []*BoltQueryCondition `json:"filters"`
}

type BoltQueryCondition struct {
	Field      interface{} `json:"field"`      // absolute value or field "doc.name"
	Comparator string      `json:"comparator"` // ==, !=, > ...
	Value      interface{} `json:"value"`      // absolute value or field "doc.surname", "Rossi"
}

//----------------------------------------------------------------------------------------------------------------------
//	BoltQuery
//----------------------------------------------------------------------------------------------------------------------

func NewQueryFromFile(path string) (*BoltQuery, error) {
	query := new(BoltQuery)
	text, err := qbc.IO.ReadTextFromFile(path)
	if nil != err {
		return nil, err
	}
	err = query.Parse(text)
	if nil != err {
		return nil, err
	}
	return query, nil
}

func (instance *BoltQuery) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *BoltQuery) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *BoltQuery) MatchFilter(entity interface{}) bool {
	response := false
	if nil != instance {
		if nil != instance.Conditions && nil != entity {
			conditions := instance.Conditions
			if len(conditions) > 0 {
				// OPTIMISTIC CONDITION
				response = true
				for _, condition := range conditions {
					if nil != condition {
						operator := condition.Operator
						filters := condition.Filters
						for _, filter := range filters {
							f1 := getValue(entity, filter.Field)
							f2 := getValue(entity, filter.Value)
							cp := filter.Comparator

							match := false
							switch cp {
							case ComparatorEqual:
								match = qbc.Compare.Equals(f1, f2)
							case ComparatorNotEqual:
								match = qbc.Compare.NotEquals(f1, f2)
							default:
								match = false
							}

							switch operator {
							case OperatorAnd:
								// AND
								response = response && match
							case OperatorOr:
								// OR
								response = response || match
							default:
								// invalid operator
								response = false
							}
							if !response {
								break
							}
						} // filters
						if !response {
							break
						}
					}
				} // conditions
			}
		}

	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func getValue(entity interface{}, propertyOrValue interface{}) interface{} {
	if b, property := qbc.Compare.IsString(propertyOrValue); b {
		if strings.Index(property, ".") > -1 {
			// found a property
			field := ""
			tokens := qbc.Strings.Split(property, ".")
			switch len(tokens) {
			case 1:
				field = tokens[0]
			default:
				field = tokens[1]
			}
			if len(field) > 0 {
				r := qbc.Reflect.Get(entity, field)
				if nil!=r{
					return r
				}
				return qbc.Reflect.Get(entity, qbc.Strings.CapitalizeAll(field))
			}
		}
	}

	return propertyOrValue
}
