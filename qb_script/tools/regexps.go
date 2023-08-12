package tools

import (
	"github.com/dop251/goja"
	qbc "github.com/rskvp/qb-core"
	"github.com/rskvp/qb-lib/qb_script/commons"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const TOOL_REGEXPS = "$regexps"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScriptingToolRegExps struct {
	params  *ScriptingToolParams
	runtime *goja.Runtime
	context interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewToolRegExps(params *ScriptingToolParams) *ScriptingToolRegExps {
	result := new(ScriptingToolRegExps)
	result.params = params
	result.runtime = params.Runtime

	return result
}

//----------------------------------------------------------------------------------------------------------------------
//	i n t e r f a c e
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolRegExps) Init(params *ScriptingToolParams) {
	tool.params = params
	tool.runtime = params.Runtime
}

func (tool *ScriptingToolRegExps) SetContext(context interface{}) {
	tool.context = context
}

//----------------------------------------------------------------------------------------------------------------------
//	r e g e x p
//----------------------------------------------------------------------------------------------------------------------

func (tool *ScriptingToolRegExps) MatchNumbers(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		text := commons.GetArgsString(tool.context, args)
		if len(text) > 0 {

			result := qbc.Regex.Numbers(text)
			return tool.runtime.ToValue(result)
		}
	}

	return tool.runtime.ToValue([]string{})
}

func (tool *ScriptingToolRegExps) MatchPrices(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		text := commons.GetArgsString(tool.context, args)
		if len(text) > 0 {

			result := qbc.Regex.Prices(text)
			return tool.runtime.ToValue(result)
		}
	}

	return tool.runtime.ToValue([]string{})
}

//----------------------------------------------------------------------------------------------------------------------
//	r e g e x p
//----------------------------------------------------------------------------------------------------------------------

// Return true if passed pattern match with text
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return bool
func (tool *ScriptingToolRegExps) HasMatchExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.Match(text, pattern)
			v := len(result) > 0
			return tool.runtime.ToValue(v)
		}
	}

	return tool.runtime.ToValue(false)
}

// Return all matched values with passed pattern match with text
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchAllExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.Match(text, pattern)
			return tool.runtime.ToValue(result)
		}
	}

	return tool.runtime.ToValue([]string{})
}

// Return string  matched at index of matched array (works on array of all matched values) with passed pattern match with text
// If index passed is zero, response is same as $regexps.MatchFirst('hu?an*')
// usage: $regexps.MatchAt('hu?an*', 0)
// @param pattern string A standard regular expression
// @param index int
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchAtExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, index, text := tool.getArgsStringIntString(args)
		if len(pattern) > 0 && index > -1 && len(text) > 0 {

			result := qbc.Regex.Match(text, pattern)
			if len(result) > index {
				return tool.runtime.ToValue(result[index])
			}
		}
	}
	return tool.runtime.ToValue("")
}

// Return first matched value with passed pattern match with text
// usage: $regexps.MatchFirst('hu?an*')
// @param pattern string  A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return string
func (tool *ScriptingToolRegExps) MatchFirstExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.Match(text, pattern)
			if len(result) > 0 {
				return tool.runtime.ToValue(result[0])
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Return last matched value with passed pattern match with text
// usage: $regexps.MatchLast('hu?an*')
// @param pattern string  A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return string
func (tool *ScriptingToolRegExps) MatchLastExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.Match(text, pattern)
			if len(result) > 0 {
				return tool.runtime.ToValue(result[len(result)-1])
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Return an array of values detected between two expressions.
// @param offset int Index to start matching from. Pass 0 if you do not need any offset
// @param patternStart string A standard regular expression
// @param patternEnd string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchBetweenExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		offset, patternStart, patterEnd, text := tool.getArgsIntStringStringString(args)
		if len(patternStart) > 0 && len(patterEnd) > 0 && len(text) > 0 {

			result := qbc.Regex.MatchBetween(text, offset, patternStart, patterEnd, " ")
			if nil != result {
				return tool.runtime.ToValue(result)
			}
		}
	}
	return tool.runtime.ToValue([]string{})
}

// Return index array of matching expression in a text
// usage: $regexps.Index('hu?an*')
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int
func (tool *ScriptingToolRegExps) IndexExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.MatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

// Return index of first matching expression in a text
// usage: $regexps.IndexFirst('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return int
func (tool *ScriptingToolRegExps) IndexFirstExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.MatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				index := response[0][0]
				return tool.runtime.ToValue(index)
			}
		}
	}

	return tool.runtime.ToValue(-1)
}

// Return index of last matching expression in a text
// usage: $regexps.IndexLast('hu?an*')
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return int
func (tool *ScriptingToolRegExps) IndexLastExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.MatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				index := response[len(response)-1][0]
				return tool.runtime.ToValue(index)
			}
		}
	}

	return tool.runtime.ToValue(-1)
}

// Return index array of matching expression in a text starting search from offset position
// usage: $regexps.IndexStartAt(30, 'hu?an*')
// @param offset int number of characters to exclude from search
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int
func (tool *ScriptingToolRegExps) IndexStartAtExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		offset, pattern, text := commons.GetArgsIntStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {
			response := qbc.Regex.Index(text, pattern, offset)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

// Return array of pair index:word_len  of matching expression in a text
// usage: $regexps.IndexLenPair('hu?an*')
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return [][]int ex: [[12,3], [22,4]]
func (tool *ScriptingToolRegExps) IndexLenPairExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.IndexLenPair(text, pattern, 0)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([][]int{})
}

// Return last array of pair index:word_len  of matching expression in a text
// usage: $regexps.IndexLenPair('hu?an*')
// @param pattern string A standard regular expression
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int ex: [22,4]
func (tool *ScriptingToolRegExps) IndexLenPairLastExp(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.IndexLenPair(text, pattern, 0)
			if nil != response && len(response) > 0 {
				lastIndex := len(response) - 1
				return tool.runtime.ToValue(response[lastIndex])
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

//----------------------------------------------------------------------------------------------------------------------
//	w i l d c a r d
//----------------------------------------------------------------------------------------------------------------------

// Return true if passed pattern match with text
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return bool
func (tool *ScriptingToolRegExps) HasMatch(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.WildcardMatch(text, pattern)
			v := len(result) > 0
			return tool.runtime.ToValue(v)
		}
	}

	return tool.runtime.ToValue(false)
}

// Return all matched values with passed pattern match with text
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchAll(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.WildcardMatch(text, pattern)
			return tool.runtime.ToValue(result)
		}
	}

	return tool.runtime.ToValue([]string{})
}

// Return string  matched at index of matched array (works on array of all matched values) with passed pattern match with text
// If index passed is zero, response is same as $regexps.MatchFirst('hu?an*')
// usage: $regexps.MatchAt('hu?an*', 0)
// @param pattern string "hu?an*"
// @param index int
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, index, text := tool.getArgsStringIntString(args)
		if len(pattern) > 0 && index > -1 && len(text) > 0 {

			result := qbc.Regex.WildcardMatch(text, pattern)
			if len(result) > index {
				return tool.runtime.ToValue(result[index])
			}
		}
	}
	return tool.runtime.ToValue("")
}

// Return first matched value with passed pattern match with text
// usage: $regexps.MatchFirst('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return string
func (tool *ScriptingToolRegExps) MatchFirst(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.WildcardMatch(text, pattern)
			if len(result) > 0 {
				return tool.runtime.ToValue(result[0])
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Return last matched value with passed pattern match with text
// usage: $regexps.MatchLast('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return string
func (tool *ScriptingToolRegExps) MatchLast(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			result := qbc.Regex.WildcardMatch(text, pattern)
			if len(result) > 0 {
				return tool.runtime.ToValue(result[len(result)-1])
			}
		}
	}

	return tool.runtime.ToValue("")
}

// Return an array of values detected between two expressions.
// @param offset int Index to start matching from. Pass 0 if you do not need any offset
// @param patternStart string
// @param patternEnd string
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []string
func (tool *ScriptingToolRegExps) MatchBetween(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		offset, patternStart, patterEnd, text := tool.getArgsIntStringStringString(args)
		if len(patternStart) > 0 && len(patterEnd) > 0 && len(text) > 0 {

			result := qbc.Regex.WildcardMatchBetween(text, offset, patternStart, patterEnd, " ")
			if nil != result {
				return tool.runtime.ToValue(result)
			}
		}
	}
	return tool.runtime.ToValue([]string{})
}

//----------------------------------------------------------------------------------------------------------------------
//	t e x t     s e a r c h
//----------------------------------------------------------------------------------------------------------------------

// Return index array of matching expression in a text
// usage: $regexps.Index('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int
func (tool *ScriptingToolRegExps) Index(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.WildcardMatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

// Return index of first matching expression in a text
// usage: $regexps.IndexFirst('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return int
func (tool *ScriptingToolRegExps) IndexFirst(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.WildcardMatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				index := response[0][0]
				return tool.runtime.ToValue(index)
			}
		}
	}

	return tool.runtime.ToValue(-1)
}

// Return index of last matching expression in a text
// usage: $regexps.IndexLast('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return int
func (tool *ScriptingToolRegExps) IndexLast(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.WildcardMatchIndex(text, pattern)
			if nil != response && len(response) > 0 {
				index := response[len(response)-1][0]
				return tool.runtime.ToValue(index)
			}
		}
	}

	return tool.runtime.ToValue(-1)
}

// Return index array of matching expression in a text starting search from offset position
// usage: $regexps.IndexStartAt(30, 'hu?an*')
// @param offset int number of characters to exclude from search
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int
func (tool *ScriptingToolRegExps) IndexStartAt(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		offset, pattern, text := commons.GetArgsIntStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {
			response := qbc.Regex.WildcardIndex(text, pattern, offset)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

// Return array of pair index:word_len  of matching expression in a text
// usage: $regexps.IndexLenPair('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return [][]int ex: [[12,3], [22,4]]
func (tool *ScriptingToolRegExps) IndexLenPair(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.WildcardIndexLenPair(text, pattern, 0)
			if nil != response && len(response) > 0 {
				return tool.runtime.ToValue(response)
			}
		}
	}

	return tool.runtime.ToValue([][]int{})
}

// Return last array of pair index:word_len  of matching expression in a text
// usage: $regexps.IndexLenPair('hu?an*')
// @param pattern string "hu?an*"
// @param text string (Optional) CONTEXT is used if not found. "hello humanity!!"
// @return []int ex: [22,4]
func (tool *ScriptingToolRegExps) IndexLenPairLast(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		pattern, text := commons.GetArgsStringString(tool.context, args)
		if len(pattern) > 0 && len(text) > 0 {

			response := qbc.Regex.WildcardIndexLenPair(text, pattern, 0)
			if nil != response && len(response) > 0 {
				lastIndex := len(response) - 1
				return tool.runtime.ToValue(response[lastIndex])
			}
		}
	}

	return tool.runtime.ToValue([]int{})
}

//----------------------------------------------------------------------------------------------------------------------
//	n a t u r a l   l a n g u a g e    p r o c e s s i n g
//----------------------------------------------------------------------------------------------------------------------

// Calculate a matching score between a phrase and a check test using expressions.
// ALL expressions are evaluated.
// Score result is different if there's more than one expression (separated by "|" symbol).
// If mode equals "all": Failed expressions add negative score to result
// If mode equals "any": Failed expressions do not add negative score to result
// If mode equals "best": Failed expressions do not add negative score to result and best score is returned
// @param [string] phrase. "hello humanity!! I'm Mario rossi"
// @param [string] expressions. All expressions to match separated by | (pipe) hel??0 h* | I* * ros*"
// @param [string] mode. "all", "any", "best"
func (tool *ScriptingToolRegExps) Score(call goja.FunctionCall) goja.Value {
	args := call.Arguments
	if len(args) > 0 {
		phrase, expressionsTxt, mode := commons.GetArgsStringStringString(tool.context, args)
		if len(phrase) > 0 && len(expressionsTxt) > 0 {
			expressions := qbc.Strings.SplitTrimSpace(expressionsTxt, "|")
			if len(expressions) > 0 {
				switch mode {
				case "all":
					return tool.runtime.ToValue(qbc.Regex.WildcardScoreAll(phrase, expressions))
				case "any":
					return tool.runtime.ToValue(qbc.Regex.WildcardScoreAny(phrase, expressions))
				default:
					return tool.runtime.ToValue(qbc.Regex.WildcardScoreBest(phrase, expressions))
				}
			}
		}
	}
	return tool.runtime.ToValue(0.0)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------


func (tool *ScriptingToolRegExps) getArgsStringIntString(args []goja.Value) (string, int, string) {
	var arg1 string
	var arg2 int
	var argCtx string

	if len(args) > 0 {
		arg1 = qbc.Convert.ToString(args[0].Export())
		if len(arg1) > 0 {

			if len(args) > 1 {
				arg2 = qbc.Convert.ToInt(args[1].Export())
			}

			if len(args) > 2 {
				argCtx = qbc.Convert.ToString(args[2].Export())
			}
		}

		// fallback on context for latest arg
		if len(argCtx) == 0 {
			if nil != tool.context {
				argCtx = qbc.Convert.ToString(tool.context)
			}
		}
	}

	return arg1, arg2, argCtx
}

func (tool *ScriptingToolRegExps) getArgsIntStringStringString(args []goja.Value) (int, string, string, string) {
	var arg1 int
	var arg2 string
	var arg3 string
	var argCtx string

	if len(args) > 0 {
		arg1 = qbc.Convert.ToInt(args[0].Export())

		if len(args) > 1 {
			arg2 = qbc.Convert.ToString(args[1].Export())
		}

		if len(args) > 2 {
			arg3 = qbc.Convert.ToString(args[2].Export())
		}

		if len(args) > 3 {
			argCtx = qbc.Convert.ToString(args[3].Export())
		}

		// fallback on context for latest arg
		if len(argCtx) == 0 {
			if nil != tool.context {
				argCtx = qbc.Convert.ToString(tool.context)
			}
		}
	}

	return arg1, arg2, arg3, argCtx
}
