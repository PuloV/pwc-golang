// You can edit this code!
// Click here and start typing.
package pwchp

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type PwcHtmlParser struct {
	domain  string
	url     *url.URL
	content string
}

func generateUrl(link string) *url.URL {
	u, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func NewParser(link string) *PwcHtmlParser {
	p := new(PwcHtmlParser)
	p.url = generateUrl(link)
	p.GetContent()
	p.GetAllPageUrls()
	return p
}

func (p *PwcHtmlParser) getPageUrl() string {
	pageUrl := fmt.Sprintf("%s", p.url)
	return pageUrl
}

func (p *PwcHtmlParser) GetContent() {

	resp, err := http.Get(p.getPageUrl())

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	p.content = string(body)
}

func (p *PwcHtmlParser) GetAllTags(tag string) []string {
	re := regexp.MustCompile(fmt.Sprintf(`<%s.*?</%s>`, tag, tag))
	return re.FindAllString(p.content, -1)
}

func (p *PwcHtmlParser) GetAllPageUrls() []string {
	rawLinks := p.GetAllTags("a")
	re := regexp.MustCompile(`.*?href=('|")(.*?)('|").*`)
	for _, rawLink := range rawLinks {
		originalHref := re.ReplaceAllString(rawLink, "$2")
		parsedHref := p.parseHref(originalHref)
		fmt.Println(parsedHref)
	}
	return rawLinks
}

func (p *PwcHtmlParser) parseHref(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		log.Fatal(err)
	}

	if u.Host == "" {
		u.Host = p.url.Host
	}

	if u.Scheme == "" {
		u.Scheme = p.url.Scheme
	}
	newHref := fmt.Sprintf("%s", u)
	return newHref
}
