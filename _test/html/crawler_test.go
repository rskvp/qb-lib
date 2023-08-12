package _test

import (
	"fmt"
	"testing"

	ggx "bitbucket.org/digi-sense/gg-core-x"
	"github.com/rskvp/qb-lib/qb_html"
)

func TestBlacklist(t *testing.T) {
	urls := []string{"https://gianangelogeminiani.me",
		"https://www.facebook.com/angelo.geminiani/about?lst=1472675714%3A1472675714%3A1591518292",
		"http://facebook.com/angelo.geminiani/about?lst=1472675714%3A1472675714%3A1591518292",
	}
	for _, url := range urls {
		match := ggx.HTML.UrlMatch(url, qb_html.DefaultBlackList)
		fmt.Println(match, url)
	}
}

func TestCrawler(t *testing.T) {

	settings := new(qb_html.HtmlCrawlerSettings)
	settings.MaxThreads = 2
	settings.StartPoints = []string{"https://gianangelogeminiani.me"}
	settings.AllowExternals = true
	settings.WhiteList = []string{"https://gianangelogeminiani.me/*"}
	settings.BlackList = []string{"https://github.com/*"}
	crawler := ggx.HTML.NewCrawler(settings)

	crawler.OnContent(func(content *qb_html.HtmlCrawlerContend) {
		fmt.Println(content.Url)
		fmt.Println("\t", "error", content.Error)
		fmt.Println("\t", "links", content.Links)
		fmt.Println("\t", "blocks", len(content.Blocks))
	})

	// start and wait
	crawler.Start()
	crawler.Join()
}

func TestCrawlerLocal(t *testing.T) {

	settings := new(qb_html.HtmlCrawlerSettings)
	settings.MaxThreads = 2
	settings.StartPoints = []string{"./pages/index.html"}

	crawler := ggx.HTML.NewCrawler(settings)
	crawler.Start()

}
