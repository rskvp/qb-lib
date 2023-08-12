package qb_html

import qbc "github.com/rskvp/qb-core"

var (
	DefaultBlackList = []string{
		"http*//*facebook.*",
		"http*//*github.*",
		"http*//*linkedin.*",
		"http*//*bitbucket.*",
		"http*//*pinterest.*",
		"http*//*instagram.*",
		"http*//*twitter.*",
		"http*//*telegram.*",
		"http*//*google.*",
		"http*//*repubblica.*",
		"http*//*akismet.*",
		"http*//*jetpack.*",
	}
)

type HTMLHelper struct {
}

var HTML *HTMLHelper

func init() {
	HTML = new(HTMLHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HTMLHelper) NewParser(input interface{}) (*HtmlParser, error) {
	return NewHtmlParser(input)
}

func (instance *HTMLHelper) NewCrawler(settings *HtmlCrawlerSettings) *HtmlCrawler {
	return NewHtmlCrawler(settings)
}

func (instance *HTMLHelper) LoadCrawlerSettings(filename string) (*HtmlCrawlerSettings, error) {
	return LoadHtmlCrawlerSettings(filename)
}

func (instance *HTMLHelper) UrlMatch(url string, list []string) bool {
	return len(qbc.Regex.WildcardIndexArray(url, list, 0)) > 0
}

func (instance *HTMLHelper) Html2TextFromString(input string) (string, error) {
	parser, err := instance.NewParser(input)
	if nil != err {
		return "", err
	}
	return parser.TextAll(), nil
}
