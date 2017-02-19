package pwc

import (
	"fmt"
	"os"
	"pwchp"
	"sync"
	"time"
)

type Pwc struct {
	urlsToBeParsed []string
	parsedUrls     map[string]bool

	incomingUrls chan string
	outgoingUrls chan string
	bufferedChan chan string

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

				case pwc.outgoingUrls <- pwc.urlsToBeParsed[0]: // a url is taken from the url array and send to channel
					pwc.urlsToBeParsed = pwc.urlsToBeParsed[1:]

				case url := <-pwc.bufferedChan: // a url is comming from the buffer channel => save it in the array of urls that needs to be parsed
					// mark url in the wait group as not parsed
					pwc.parsingUrls.Add(1)
					pwc.urlsToBeParsed = append(pwc.urlsToBeParsed, url)
				}
			} else {
				// there are no urls that can be parsed => waiting for a url from the buffer channel
				url := <-pwc.bufferedChan
				// mark url in the wait group as not parsed
				pwc.parsingUrls.Add(1)
				pwc.urlsToBeParsed = append(pwc.urlsToBeParsed, url)
			}
		}
	}()

	pwc.addDefaultStartPoints()

	go func(urls chan string) {
		// endless loop for getting the ics urls
		for {
			url := <-urls

			go func(url string) {
				pwc.crawl(url)
			}(url)
		}
	}(pwc.outgoingUrls)
}

func (pwc *Pwc) crawl(url string) {
	// mark calendar in the wait group as  parsed
	defer pwc.parsingUrls.Done()

	p := pwchp.NewParser(url)

	logFoundData(p)

	pwc.addUrlsToBeParsed(p.GetAllPageUrls())
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
	mutex.Lock()
	if urlValue, urlExists := pwc.parsedUrls[url]; !urlExists || !urlValue {
		pwc.bufferedChan <- url
		pwc.parsedUrls[url] = true
	}
	mutex.Unlock()
}

func logFoundData(p *pwchp.PwcHtmlParser) {
	mutex.Lock()

	fmt.Printf("\n - the url %s was parsed for %s \n", p.GetPageUrl(), time.Since(p.GetStartedParsingTime()))

	f, err := os.OpenFile("./pages.log", os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	header := fmt.Sprintf("\n [%s] %s\n", time.Now(), p.GetPageUrl())
	if _, err = f.WriteString(string(header)); err != nil {
		panic(err)
	}

	for _, word := range p.GetValuableWords() {
		keyWord := fmt.Sprintf("\t - %s\n", word)
		fmt.Println(keyWord)
		if _, err = f.WriteString(string(keyWord)); err != nil {
			panic(err)
		}
	}

	footer := "================================\n"
	if _, err = f.WriteString(string(footer)); err != nil {
		panic(err)
	}

	mutex.Unlock()
}

func (p *Pwc) Wait() {
	p.parsingUrls.Wait()
}
