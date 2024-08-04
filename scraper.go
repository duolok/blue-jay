package main

import (
    "fmt"
    "github.com/gocolly/colly"
)

func main() {
    fmt.Println("hello world")
    c := colly.NewCollector(colly.AllowedDomains("www.allkeyshop.com"))
    fmt.Println(c.Visit("www.allkeyshop.com"))

}

