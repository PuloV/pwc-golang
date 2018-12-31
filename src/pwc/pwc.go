package pwc

import (
	"fmt"
	// "os"
	"sync"
	"time"

	"github.com/PuloV/pwc-golang/src/model"
	"github.com/PuloV/pwc-golang/src/pwchp"
)

type Pwc struct {
	urlsToBeParsed []string
	parsedUrls     map[string]bool

	incomingUrls chan string
	outgoingUrls chan string
	bufferedChan chan string

	maxUrslPingedPerBatch     int
	currentUrslPingedPerBatch int

	parsingUrls *sync.WaitGroup
}

var mutex sync.RWMutex

func CreateCrawler() *Pwc {
	pwc := new(Pwc)
	pwc.parsedUrls = make(map[string]bool)

	pwc.incomingUrls = make(chan string)
	pwc.outgoingUrls = make(chan string)
	pwc.bufferedChan = make(chan string)

	pwc.parsingUrls = new(sync.WaitGroup)

	return pwc
}

func (pwc *Pwc) Start() {
	go func() {
		for {
			// check if there are urls that can be parsed
			if len(pwc.urlsToBeParsed) > 0 {
				select {

				case url := <-pwc.bufferedChan: // a url is comming from the buffer channel => save it in the array of urls that needs to be parsed
					// mark url in the wait group as not parsed
					pwc.parsingUrls.Add(1)
					pwc.urlsToBeParsed = append(pwc.urlsToBeParsed, url)

				case pwc.outgoingUrls <- pwc.urlsToBeParsed[0]: // a url is taken from the url array and send to channel
					pwc.urlsToBeParsed = pwc.urlsToBeParsed[1:]
				}
			} else {
				select {

				case url := <-pwc.bufferedChan: // there are no urls that can be parsed => waiting for a url from the buffer channel
					// mark url in the wait group as not parsed
					pwc.parsingUrls.Add(1)
					pwc.urlsToBeParsed = append(pwc.urlsToBeParsed, url)

				case <-time.After(5 * time.Second):
					continue
				}
			}
		}
	}()

	pwc.addDefaultStartPoints()

	// for i := 0; i < 4; i++ {
	go func(urls chan string) {
		// endless loop for getting the urls
		for {
			url := <-urls
			go func(url string) {
				pwc.crawl(url)
			}(url)
		}
	}(pwc.outgoingUrls)
	// }
}

func (pwc *Pwc) crawl(url string) {
	// mark calendar in the wait group as  parsed
	defer pwc.parsingUrls.Done()
	p := pwchp.NewParser(url)

	logFoundData(p)
	pwc.addUrlsToBeParsed(p.GetAllPageUrls())
}

func (pwc *Pwc) addDefaultStartPoints() {
	defaultStartPoints := []string{"http://mobile.bg", "http://cars.bg"}
	pwc.addUrlsToBeParsed(defaultStartPoints)
}

func (pwc *Pwc) addUrlsToBeParsed(urls []string) {
	for _, url := range urls {
		pwc.addUrlToBeParsed(url)
	}
}

func (pwc *Pwc) addUrlToBeParsed(url string) {
	mutex.Lock()
	if _, urlExists := pwc.parsedUrls[url]; !urlExists {
		pwc.bufferedChan <- url
		pwc.parsedUrls[url] = true
	}
	mutex.Unlock()
}

func logFoundData(p *pwchp.PwcHtmlParser) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("\n - the url %s [%d] was parsed for %s \n", p.GetPageUrl(), p.GetStatusCode(), time.Since(p.GetStartedParsingTime()))
	page := model.Page{
		URL:          model.ConvertToUTF8(p.GetPageUrl()),
		Domain:       model.ConvertToUTF8(p.GetDomain()),
		ResponseCode: uint(p.GetStatusCode()),
		LoadTime:     uint(time.Since(p.GetStartedParsingTime()).Seconds())}
	pageError := page.Save()

	if pageError != nil {
		fmt.Printf("ERROR = \n", pageError)
		return
	}

	go func(urls []string, pageId uint) {
		for _, word := range p.GetValuableWords() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("ERROR = \n", err)
				}
			}()

			keyword := model.KeyWord{
				PageID:  pageId,
				Weight:  0,
				KeyWord: model.ConvertToUTF8(word)}
			keyword.Save()
		}
	}(p.GetValuableWords(), page.ID)
}

func (p *Pwc) Wait() {
	p.parsingUrls.Wait()
}
