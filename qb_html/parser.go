package qb_html

import (
	"bytes"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	qbc "github.com/rskvp/qb-core"
	"golang.org/x/net/html"
)

//----------------------------------------------------------------------------------------------------------------------
//	SemanticBlock
//----------------------------------------------------------------------------------------------------------------------

type SemanticBlock struct {
	Lang     string   `json:"lang"`  // detected language (maybe different from page lang)
	Level    int      `json:"level"` // title level. 0 is when a block is free text with no title
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Keywords []string `json:"keywords"`

	//-- private --//
	body bytes.Buffer
}

func (instance *SemanticBlock) Json() string {
	if nil != instance {
		instance.Body = instance.GetBody()
		return qbc.JSON.Stringify(instance)
	}
	return ""
}

func (instance *SemanticBlock) GetBody() string {
	if nil != instance {
		return instance.body.String()
	}
	return ""
}

func (instance *SemanticBlock) GetText() string {
	if nil != instance {
		var buf bytes.Buffer
		if len(instance.Title) > 0 {
			buf.WriteString(instance.Title)
			buf.WriteString("\n")
		}
		if instance.body.Len() > 0 {
			buf.Write(instance.body.Bytes())
		}
		return buf.String()
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	HtmlParser
//----------------------------------------------------------------------------------------------------------------------

type HtmlParser struct {
	html     *html.Node
	path     string
	fileName string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHtmlParser(input interface{}) (*HtmlParser, error) {
	var doc *html.Node
	var err error
	var path string
	var fileName string

	// parse input data
	if v, b := input.(string); b {
		if isURL(v) {
			doc, err = parseURL(v)
			path = v
			fileName = qbc.Paths.FileName(path, true)
		} else if b, err := qbc.Paths.IsFile(v); b && nil == err {
			doc, err = parseFile(v)
			path = qbc.Paths.Absolute(v)
			fileName = qbc.Paths.FileName(path, true)
		} else {
			doc, err = parseString(v)
		}
	} else if v, b := input.(io.Reader); b {
		doc, err = parse(v)
	}

	if nil != err {
		return nil, err
	} else {
		instance := new(HtmlParser)
		instance.html = doc
		instance.path = path
		instance.fileName = fileName

		return instance, nil
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HtmlParser) String() string {
	if nil != instance {
		return renderNode(instance.html)
	}
	return ""
}

func (instance *HtmlParser) Document() *html.Node {
	return instance.html
}

func (instance *HtmlParser) Path() string {
	if nil != instance {
		return instance.path
	}
	return ""
}

func (instance *HtmlParser) FileName() string {
	if nil != instance {
		if len(instance.fileName) == 0 && instance.BaseUrl() == instance.RootUrl() {
			return "index.html"
		}
		return instance.fileName
	}
	return ""
}

func (instance *HtmlParser) RootUrl() string {
	if nil != instance {
		if instance.IsURL() {
			uri, err := url.Parse(instance.path)
			if nil == err {
				path := uri.Scheme + "://" + uri.Host
				if len(uri.Port()) > 0 && (uri.Port() != "80" || uri.Port() != "443") {
					path += ":" + uri.Port()
				}
				return filepath.Join(path, "/")
			}
		}
		return filepath.Dir(instance.path)
	}
	return ""
}

func (instance *HtmlParser) BaseUrl() string {
	if nil != instance {
		if instance.IsURL() {
			uri, err := url.Parse(instance.path)
			if nil == err {
				path := uri.Scheme + "://" + uri.Host
				if len(uri.Port()) > 0 && (uri.Port() != "80" || uri.Port() != "443") {
					path += ":" + uri.Port()
				}
				ext := filepath.Ext(uri.Path)
				if len(ext) > 0 {
					path += filepath.Dir(uri.Path)
				} else {
					path += uri.Path
				}
				return filepath.Join(path, "/")
			}
		}
		return filepath.Dir(instance.path)
	}
	return ""
}

func (instance *HtmlParser) IsURL() bool {
	if nil != instance {
		return isURL(instance.path)
	}
	return false
}

func (instance *HtmlParser) Lang() string {
	lang := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			lang = getAttr(node, "lang")
			return len(lang) > 0 // next node?
		})
	}
	return lang
}

func (instance *HtmlParser) Title() string {
	title := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			if strings.ToLower(node.Data) == "title" {
				title = instance.InnerHtml(node)
			}
			return len(title) > 0 // next node?
		})
	}
	return title
}

func (instance *HtmlParser) MetaTitle() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("title")
	}
	return value
}

func (instance *HtmlParser) MetaDescription() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("description")
	}
	return value
}

func (instance *HtmlParser) MetaAuthor() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("author")
	}
	return value
}

func (instance *HtmlParser) MetaKeywords() []string {
	if nil != instance && nil != instance.html {
		return qbc.Strings.SplitTrimSpace(instance.GetMetaContent("keywords"), ",")
	}
	return []string{}
}

func (instance *HtmlParser) GetText() []string {
	if nil != instance && nil != instance.html {
		return qbc.Strings.SplitTrimSpace(instance.GetMetaContent("keywords"), ",")
	}
	return []string{}
}

func (instance *HtmlParser) GetMetaContent(name string) string {
	value := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			if strings.ToLower(node.Data) == "meta" {
				attrName := getAttr(node, "name")
				if attrName == name {
					value = getAttr(node, "content")
				}
			}
			return len(value) > 0 // next node?
		})
	}
	return value
}

func (instance *HtmlParser) TextAll() string {
	if nil != instance {
		return text(instance.html)
	}
	return ""
}

func (instance *HtmlParser) Text(node *html.Node) string {
	if nil != instance {
		return text(node)
	}
	return ""
}

func (instance *HtmlParser) SemanticBlocksAll() []*SemanticBlock {
	if nil != instance {
		return semantic(instance.Lang(), instance.html)
	}
	return make([]*SemanticBlock, 0)
}

func (instance *HtmlParser) SemanticBlocks(node *html.Node) []*SemanticBlock {
	if nil != instance {
		return semantic(instance.Lang(), node)
	}
	return make([]*SemanticBlock, 0)
}

func (instance *HtmlParser) OuterHtml(n *html.Node) string {
	if nil != instance {
		return renderNode(n)
	}
	return ""
}

func (instance *HtmlParser) InnerHtml(n *html.Node) string {
	if nil != instance {
		return renderNode(n.FirstChild)
	}
	return ""
}

func (instance *HtmlParser) ForEach(callback func(node *html.Node) bool) {
	if nil != instance && nil != instance.html && nil != callback {
		forEach(instance.html, callback)
	}
}

func (instance *HtmlParser) GelLinks() []*html.Node {
	if nil != instance && nil != instance.html {
		return queryNodes(instance.html, "a")
	}
	return []*html.Node{}
}

func (instance *HtmlParser) GetLinkURLs() []string {
	response := make([]string, 0)
	if nil != instance && nil != instance.html {
		links := queryNodes(instance.html, "a")
		for _, link := range links {
			href := getAttr(link, "href")
			if len(href) > 0 && strings.Index(href, "#") != 0 {
				response = qbc.Arrays.AppendUnique(response, href).([]string)
			}
		}
	}
	return response
}

func (instance *HtmlParser) GeNodeAttributes(nodes []*html.Node) []map[string]string {
	response := make([]map[string]string, 0)
	if nil != instance && nil != instance.html {
		for _, node := range nodes {
			m := map[string]string{
				"tag": node.Data,
			}
			for _, attr := range node.Attr {
				m[attr.Key] = attr.Val
			}
			response = append(response, m)
		}
	}
	return response
}

func (instance *HtmlParser) Select(selector string) []*html.Node {
	if nil != instance && nil != instance.html {
		return queryNodes(instance.html, selector)
	}
	return []*html.Node{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func parse(input io.Reader) (*html.Node, error) {
	return html.Parse(input)
}

func parseString(input string) (*html.Node, error) {
	if strings.Index(input, "http") == 0 {
		return parseURL(input)
	} else if b, err := qbc.Paths.IsFile(input); b && nil == err {
		return parseFile(input)
	}
	return html.Parse(strings.NewReader(input))
}

func parseBytes(input []byte) (*html.Node, error) {
	return html.Parse(bytes.NewReader(input))
}

func parseFile(input string) (*html.Node, error) {
	text, err := qbc.IO.ReadTextFromFile(input)
	if nil != err {
		return nil, err
	}
	return html.Parse(strings.NewReader(text))
}

func parseURL(url string) (*html.Node, error) {
	data, err := qbc.IO.Download(url)
	if nil != err {
		return nil, err
	}
	return parseBytes(data)
}

func renderNode(n *html.Node) string {
	if nil != n {
		var buf bytes.Buffer
		w := io.Writer(&buf)
		_ = html.Render(w, n)
		return buf.String()
	}
	return ""
}

func isURL(path string) bool {
	return strings.Index(path, "http") == 0
}

func isTitle(node *html.Node) (bool, int) {
	tag := strings.ToLower(node.Data)
	if len(tag) == 2 && strings.Index(tag, "h") == 0 {
		level := tag[1]
		return true, qbc.Convert.ToInt(string(level))
	}
	return false, -1
}

func text(node *html.Node) string {
	var buf bytes.Buffer
	forEach(node, func(node *html.Node) bool {
		if node.Data != "title" && nil != node.FirstChild && node.FirstChild.Type == html.TextNode {
			text := qbc.Strings.Clear(renderNode(node.FirstChild))
			if len(text) > 0 {
				buf.WriteString(text + "\n")
			}
		}
		return false // next
	})
	return buf.String()
}

func unescape(text string) string {
	return html.UnescapeString(qbc.Strings.Clear(text))
}

func semantic(lang string, node *html.Node) []*SemanticBlock {
	response := make([]*SemanticBlock, 0)
	tmpBlock := new(SemanticBlock)
	tmpBlock.Lang = lang
	forEach(node, func(node *html.Node) bool {
		if node.Data != "title" && nil != node.FirstChild && node.FirstChild.Type == html.TextNode {
			if b, level := isTitle(node); b {
				// new block
				if tmpBlock.body.Len() > 0 || len(tmpBlock.Title) > 0 {
					response = append(response, tmpBlock)
				}
				tmpBlock = new(SemanticBlock)
				tmpBlock.Lang = lang
				tmpBlock.Title = unescape(renderNode(node.FirstChild))
				tmpBlock.Level = level
			} else {
				text := unescape(renderNode(node.FirstChild))
				if len(text) > 0 {
					tmpBlock.body.WriteString(text + "\n")
				}
			}
		}
		return false // next
	})

	if tmpBlock.body.Len() > 0 {
		response = append(response, tmpBlock)
	}

	// semantic check
	for _, block := range response {
		// detect language
		//lang := detectLanguage(block.GetText())
		if len(lang) > 0 {
			block.Lang = lang
		}
		// detect keywords

	}

	return response
}

func forEach(node *html.Node, callback func(node *html.Node) bool) {
	if nil != node {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode || nil != child.FirstChild {
				exit := callback(child)
				if exit {
					break
				}
				forEach(child, callback)
			}
		}
	}
}

func getAttr(node *html.Node, name string) string {
	if nil != node && len(node.Attr) > 0 {
		for _, attr := range node.Attr {
			if strings.ToLower(attr.Key) == strings.ToLower(name) {
				return attr.Val
			}
		}
	}
	return ""
}

func queryNodes(root *html.Node, selector string) []*html.Node {
	response := make([]*html.Node, 0)
	forEach(root, func(node *html.Node) bool {
		if matches(node, selector) {
			response = append(response, node)
		}
		return false // next node
	})
	return response
}

func matches(node *html.Node, selector string) bool {
	if strings.Index(selector, ".") == 0 {
		// class matching
		className := strings.Replace(selector, ".", "", 1)
		classes := strings.Split(getAttr(node, "class"), " ")
		return qbc.Arrays.IndexOf(className, classes) > -1
	} else {
		// tag name matching
		return node.Data == selector
	}
}
