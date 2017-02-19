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
	"strings"
	"time"
)

type PwcHtmlParser struct {
	domain             string
	url                *url.URL
	content            string
	startedParsingTime time.Time
}

func generateUrl(link string) *url.URL {

	u, err := url.Parse(link)
	if err != nil {
		return new(url.URL)
		log.Fatal(err)
	}
	return u
}

func NewParser(link string) *PwcHtmlParser {
	p := new(PwcHtmlParser)
	p.startedParsingTime = time.Now()
	p.url = generateUrl(link)
	p.GetContent()
	return p
}

func (p *PwcHtmlParser) GetPageUrl() string {
	pageUrl := fmt.Sprintf("%s", p.url)
	return pageUrl
}

func (p *PwcHtmlParser) GetContent() {

	resp, err := http.Get(p.GetPageUrl())

	if err != nil {
		return
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
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
	re := regexp.MustCompile(`.*? href=('|")(.*?)('|") .*`)

	parsedHrefs := []string{}

	for _, rawLink := range rawLinks {
		originalHref := re.ReplaceAllString(rawLink, "$2")
		parsedHref := p.parseHref(originalHref)
		parsedHrefs = append(parsedHrefs, parsedHref)
	}

	uniqueValues := map[string]bool{}
	uniqueParsedHrefs := []string{}

	checkBadHref := regexp.MustCompile(`(mailto:.*|.*javascript.*|.*\.jpg$|.*\.git$)`)

	for _, href := range parsedHrefs {
		isBadHref := checkBadHref.MatchString(href)

		if urlValue, urlExists := uniqueValues[href]; !isBadHref && (!urlExists || !urlValue) {
			uniqueValues[href] = true
		}
	}

	for k := range uniqueValues {
		if len(k) >= 3 {
			uniqueParsedHrefs = append(uniqueParsedHrefs, k)
		}
	}

	return uniqueParsedHrefs
}

func (p *PwcHtmlParser) parseHref(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
		log.Fatal(err)
	}

	if u.Host == "" {
		u.Host = p.url.Host
	}

	if u.Scheme == "" {
		u.Scheme = p.url.Scheme
	}

	newHref := fmt.Sprintf("%s", u)

	re := regexp.MustCompile(`\s*`)
	newHref = re.ReplaceAllString(newHref, "")

	return newHref
}

func (p *PwcHtmlParser) GetValuableWords() []string {
	tagTypes := []string{"strong", "b", "i", "a", "td"}

	valuableWords := []string{}

	// collect the valuable words from the page
	for _, tagType := range tagTypes {
		tags := p.GetAllTags(tagType)

		for _, tagHtml := range tags {
			removeMainTag := regexp.MustCompile(fmt.Sprintf(`<%s.*?>(.*?)</%s>`, tagType, tagType))
			removeHtmlTags := regexp.MustCompile(`<[^>]*>`)
			// clear the words from the main tag
			tagInnerHtml := removeMainTag.ReplaceAllString(tagHtml, "$1")
			// clear the words any html tag
			tagInnerHtml = removeHtmlTags.ReplaceAllString(tagInnerHtml, "")

			// make the words with better view (lower and trim from spaces)
			tagInnerHtml = strings.Trim(tagInnerHtml, " ")
			tagInnerHtml = strings.ToLower(tagInnerHtml)

			valuableWords = append(valuableWords, tagInnerHtml)
		}
	}

	uniqueValuableWordsMap := map[string]bool{}
	uniqueValuableWords := []string{}

	for _, word := range valuableWords {
		if urlValue, urlExists := uniqueValuableWordsMap[word]; !urlExists || !urlValue {
			uniqueValuableWordsMap[word] = true
		}
	}

	for k := range uniqueValuableWordsMap {
		if len(k) >= 3 {
			uniqueValuableWords = append(uniqueValuableWords, k)
		}
	}

	return uniqueValuableWords
}

func (p *PwcHtmlParser) GetStartedParsingTime() time.Time {
	return p.startedParsingTime
}
