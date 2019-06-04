package servers

import (
  "fmt"
  "log"
  "net/http"
  "github.com/PuerkitoBio/goquery"
)

func GetDomainHead(domain string) (string, string) {
  response, err := http.Get(fmt.Sprintf("http://%s", domain))

  if err != nil {
    log.Fatal(err)
  }

  defer response.Body.Close()

  document, err := goquery.NewDocumentFromReader(response.Body)

  if err != nil {
    log.Fatal("Error parsing Response body. ", err)
  }

  title := document.Find("title").Text()
  logo, _  := document.Find("head [rel*='icon']").Attr("href")
  return title, logo
}
