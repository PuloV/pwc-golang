// You can edit this code!
// Click here and start typing.
package main

import (
	"github.com/PuloV/pwc-golang/src/pwc"
)

func main() {
	crawler := pwc.CreateCrawler()
	crawler.Start()
	crawler.Wait()
}
