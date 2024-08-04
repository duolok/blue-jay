package main

import (
    "fmt"
    "log"
    "github.com/gocolly/colly"
)

func main() {
    fmt.Println("hello world")
    c := colly.NewCollector(colly.AllowedDomains("www.allkeyshop.com"))
    fmt.Println(c.Visit("www.allkeyshop.com"))

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        log.Println("Something went wrong: ", err)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Visited", r.Request.URL)
    })

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        e.Request.Visit(e.Attr("href"))
    })

    c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
        fmt.Println("First column of a table row:", e.Text)
    })

    c.OnXML("//h1", func(e *colly.XMLElement) {
        fmt.Println(e.Text)
    })

    c.OnScraped(func(r *colly.Response) {
        fmt.Println("Finished", r.Request.URL)
    })
}

