package helpers

import (
	"strings"

	qbc "github.com/rskvp/qb-core"
)

func BuildInsertCommands(tableName string, data []interface{}) []string {
	// INSERT INTO test.table1 (age, first_name)
	// VALUES (12, 'Giorgio');
	response := make([]string, 0)
	for _, raw := range data {
		if item, b := raw.(map[string]interface{}); b {
			if nil != item {
				names := make([]string, 0)
				values := make([]string, 0)
				for k, v := range item {
					names = append(names, k)
					if vv, b := v.(string); b {
						values = append(values, "'"+vv+"'")
					} else {
						values = append(values, qbc.Convert.ToString(v))
					}
				}

				var buf strings.Builder
				buf.WriteString("insert into ")
				buf.WriteString(tableName)
				// names
				buf.WriteString(" (")
				for i, v := range names {
					if i > 0 {
						buf.WriteString(",")
					}
					buf.WriteString(v)
				}
				buf.WriteString(") ")
				// values
				buf.WriteString("VALUES (")
				for i, v := range values {
					if i > 0 {
						buf.WriteString(",")
					}
					buf.WriteString(v)
				}
				buf.WriteString(");")
				response = append(response, buf.String())
			}
		}
	}

	return response
}

func BuildUpdateCommand(tableName string, keyName string, keyValue interface{}, data map[string]interface{}) string {
	// UPDATE test.table1 t
	// SET t.age = 23,
	//     t.first_name = 'Marcolino d'
	// WHERE t.id = 8;
	var buf strings.Builder
	buf.WriteString("UPDATE ")
	buf.WriteString(tableName)
	buf.WriteString(" t\n")
	buf.WriteString("SET ")

	// data
	count := 0
	for k, v := range data {
		if count > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("t." + k + " = " + toString(v))
		count++
	}

	// filter
	buf.WriteString(" WHERE ")
	buf.WriteString("t." + keyName + " = " + toString(keyValue))
	buf.WriteString(";")

	// response
	return buf.String()
}

func BuildDeleteCommand(tableName string, filter string) string {
	// DELETE
	// FROM test.table1
	// WHERE id = 18;

	var buf strings.Builder
	buf.WriteString("DELETE FROM ")
	buf.WriteString(tableName)

	// filter
	if len(filter)>0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(filter)
		buf.WriteString(";")
	}

	// response
	return buf.String()
}

func BuildCountCommand(tableName string, filter string) string {
	// SELECT COUNT(*) FROM tableName
	// WHERE id > 18;
	var buf strings.Builder
	buf.WriteString("SELECT COUNT(*) FROM ")
	buf.WriteString(tableName)

	// filter
	if len(filter)>0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(filter)
		buf.WriteString(";")
	}

	// response
	return buf.String()
}

func BuildCountDistinctCommand(tableName, fieldName, filter string) string {
	// SELECT COUNT(DISTINCT fieldName ) FROM tableName
	// WHERE id > 18;
	var buf strings.Builder
	buf.WriteString("SELECT COUNT(DISTINCT " + fieldName + ") FROM ")
	buf.WriteString(tableName)

	// filter
	if len(filter)>0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(filter)
		buf.WriteString(";")
	}

	// response
	return buf.String()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toString(val interface{}) string {
	str := qbc.Convert.ToString(val)
	if v, b := val.(string); b {
		return "'" + v + "'"
	}
	return str
}

func toFilter(keyName, operator string, keyValue interface{}) string {
	var buf strings.Builder
	// filter
	if nil != keyValue || strings.TrimSpace(strings.ToUpper(operator)) == "IS NOT NULL" {
		buf.WriteString("WHERE ")
		buf.WriteString(keyName + " " + operator)
		if nil != keyValue {
			buf.WriteString(" " + toString(keyValue))
		}
		buf.WriteString(";")
	}

	// response
	return buf.String()
}
