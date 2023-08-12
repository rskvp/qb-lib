package rewrite

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	qbc "github.com/rskvp/qb-core"
)

/// fork of https://github.com/gofiber/rewrite

const actionIgnore = "#IGNORE:" // allow to declare a match path to ignore. This solve a Go regexp problem with  lookarounds
const actionRoutes = "#ROUTES:" // only routes. Exclude files

var actions = []string{actionIgnore, actionRoutes}

// Config ...
type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// Rules defines the URL path rewrite rules. The values captured in asterisk can be
	// retrieved by index e.g. $1, $2 and so on.
	// Required. Example:
	// "/old":              "/new",
	// "/api/*":            "/$1",
	// "/js/*":             "/public/javascripts/$1",
	// "/users/*/orders/*": "/user/$1/order/$2",
	Rules []*Rule

	rulesRegex []*RegexWrapper // keep rules order
}

type Rule struct {
	Key   string
	Value string
}

type RegexWrapper struct {
	matcher *regexp.Regexp
	replace string
	action  string
}

// New fiber handler
// app := fiber.New()
//
//	app.Use(rewrite.New(rewrite.Config{
//	  Rules: map[string]string{
//	    "/old":   "/new",
//	    "/old/*": "/new/$1",
//	  },
//	}))
func New(config ...interface{}) fiber.Handler {
	if len(config) == 1 {

		// build configuration
		param := config[0]
		var cfg Config
		if c, ok := param.(Config); ok {
			cfg = c
		} else if c, ok := param.(*Config); ok {
			cfg = *c
		} else if m, ok := param.(map[string]interface{}); ok {
			cfg = Config{
				Rules: make([]*Rule, 0),
			}
			for k, v := range m {
				cfg.Rules = append(cfg.Rules, &Rule{k, qbc.Convert.ToString(v)})
			}
		} else if m, ok := param.(map[string]string); ok {
			cfg = Config{
				Rules: make([]*Rule, 0),
			}
			for k, v := range m {
				cfg.Rules = append(cfg.Rules, &Rule{k, v})
			}
		} else if a, ok := param.([]*Rule); ok {
			cfg = Config{
				Rules: a,
			}
		} else if a, ok := param.([]Rule); ok {
			cfg = Config{
				Rules: make([]*Rule, 0),
			}
			for _, rule := range a {
				cfg.Rules = append(cfg.Rules, &Rule{rule.Key, rule.Value})
			}
		} else {
			cfg = Config{}
		}

		cfg.rulesRegex = make([]*RegexWrapper, 0)

		// Initialize
		for _, rule := range cfg.Rules {
			k := rule.Key
			v := rule.Value
			var action string
			// replace actions
			for _, a := range actions {
				if strings.Index(k, a) > -1 {
					k = strings.Replace(k, a, "", 1)
					action = a
					break
				}
			}

			k = strings.Replace(k, "*", "(.*)", -1)
			k = k + "$"
			rx := regexp.MustCompile(k)
			cfg.rulesRegex = append(cfg.rulesRegex, &RegexWrapper{rx, v, action})
		}

		// Middleware function
		return func(c *fiber.Ctx) error {
			// Filter request to skip middleware
			if cfg.Filter != nil && cfg.Filter(c) {
				return c.Next()
			}
			// Rewrite
			for _, r := range cfg.rulesRegex {
				rx := r.matcher
				action := r.action
				v := r.replace
				replacer := captureTokens(rx, c.Path())
				if replacer != nil {
					switch action {
					case actionIgnore:
						// do not rewrite due exclusion
						goto exit
					case actionRoutes:
						// rewrite only if path is not a file
						ext := qbc.Paths.Extension(c.Path())
						if len(ext) == 0 {
							if rewrite(c, replacer, v) {
								goto exit
							}
						}
					default:
						if rewrite(c, replacer, v) {
							goto exit
						}
					}
				}
			}
		exit:
			return c.Next()
		}
	} // Init config

	// skip middleware
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}

func rewrite(c *fiber.Ctx, replacer *strings.Replacer, v string) bool {
	path := c.Path()
	newPath := replacer.Replace(v)
	if path != newPath {
		// fmt.Println(path, newPath)
		c.Path(newPath)
		return true
	}
	return false
}

// https://github.com/labstack/echo/blob/master/middleware/rewrite.go
func captureTokens(pattern *regexp.Regexp, input string) *strings.Replacer {
	groups := pattern.FindAllStringSubmatch(input, -1)
	if groups == nil {
		return nil
	}
	values := groups[0][1:]
	replace := make([]string, 2*len(values))
	for i, v := range values {
		j := 2 * i
		replace[j] = "$" + strconv.Itoa(i+1)
		replace[j+1] = v
	}
	return strings.NewReplacer(replace...)
}
