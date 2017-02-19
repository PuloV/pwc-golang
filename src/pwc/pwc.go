package pwc

import (
	"fmt"
	"os"
	"pwchp"
	"time"
)

type Pwc struct {
	urlsToBeParsed []string
	parsedUrls     map[string]bool
}

func CreateCrawler() *Pwc {
	pwc := new(Pwc)
	pwc.parsedUrls = make(map[string]bool)

	pwc.addDefaultStartPoints()

	return pwc
}

func (pwc *Pwc) Start() {
	urlToBeCrawed := pwc.ShiftUrls()
	pwc.crawl(urlToBeCrawed)
	pwc.Start()
}

func (pwc *Pwc) ShiftUrls() string {
	url := ""
	url, pwc.urlsToBeParsed = pwc.urlsToBeParsed[0], pwc.urlsToBeParsed[1:]
	return url
}

func (pwc *Pwc) crawl(url string) {

	t := time.Now()
	p := pwchp.NewParser(url)

	logFoundData(p)

	pwc.addUrlsToBeParsed(p.GetAllPageUrls())
	fmt.Printf("\n - the url %s was parsed for %s \n", url, time.Since(t))
}

func (pwc *Pwc) addDefaultStartPoints() {
	defaultStartPoints := []string{"http://cars.bg"}
	pwc.addUrlsToBeParsed(defaultStartPoints)
}

func (pwc *Pwc) addUrlsToBeParsed(urls []string) {
	for _, url := range urls {
		pwc.addUrlToBeParsed(url)
	}
}

func (pwc *Pwc) addUrlToBeParsed(url string) {
	if urlValue, urlExists := pwc.parsedUrls[url]; !urlExists || !urlValue {
		pwc.urlsToBeParsed = append(pwc.urlsToBeParsed, url)
		pwc.parsedUrls[url] = true
	}
}

func logFoundData(p *pwchp.PwcHtmlParser) {

	f, err := os.OpenFile("./pages.log", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	header := fmt.Sprintf("\n [%s] %s\n", time.Now(), p.GetPageUrl())
	if _, err = f.WriteString(header); err != nil {
		panic(err)
	}

	for _, word := range p.GetValuableWords() {
		keyWord := fmt.Sprintf("\t - %s\n", word)
		if _, err = f.WriteString(keyWord); err != nil {
			panic(err)
		}
	}

	footer := "================================\n"
	if _, err = f.WriteString(footer); err != nil {
		panic(err)
	}
}
