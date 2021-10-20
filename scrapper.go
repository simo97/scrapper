package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {

	fName := "documents.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"title", "url", "cover"})

	c := colly.NewCollector()

	c.OnHTML(".thumbnail", func(e *colly.HTMLElement) {
		// Visit all the category list
		c.Visit(e.Request.AbsoluteURL(e.ChildAttr("a[href]", "href")))
	})

	c.OnHTML(".span5", func(e *colly.HTMLElement) {
		// visit sub categories
		c.Visit(e.Request.AbsoluteURL(e.ChildAttr("a[href]", "href")))
	})

	c.OnHTML("table[style]", func(h *colly.HTMLElement) {
		// Scrap the summary info of books in one page
		fmt.Println("Une table de livre ici ")
		h.ForEach("tr", func(i int, sub *colly.HTMLElement) {
			if sub.ChildText("a") != "" {
				writer.Write([]string{
					sub.ChildText("a"),
					h.Request.AbsoluteURL(sub.ChildAttr("a", "href")),
					sub.ChildAttr("img", "src"),
				})
			}

		})
	})

	c.OnHTML(" div > strong", func(h *colly.HTMLElement) {
		// Go through pagination to extract books
		h.ForEach("a", func(i int, sub *colly.HTMLElement) {
			_link := sub.Attr("href")
			c.Visit(h.Request.AbsoluteURL(_link))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://www.dastudy.net/docs")

}
