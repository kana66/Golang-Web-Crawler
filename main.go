package main

import (
  "log"
  "net/http"
  "github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
    // ctx.HTML("<h1>Welcome</h1>")
    address := ctx.URLParam("address")
    log.Println(address)
    if address == "" {
      ctx.JSON(iris.Map{"error": "请携带参数请求当前接口"})
    } else {
      transactions := getTransactions(address, "1")
      ctx.JSON(iris.Map{"data": transactions})
    }
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

type Transaction struct {
  TxHash  string
  Block   string
  Age     string
  From    string
  Status  string
  To      string
  Value   string
  TxFee   string
}

//eth transactions
//https://etherscan.io/txs?a=0x0b3850d16e55be91ea816fcfa02a0d8905c4f469&p=1
//erc20 transactions
//https://etherscan.io/tokentxns?a=0x0b3850d16e55be91ea816fcfa02a0d8905c4f469&p=1

func getTransactions(address string, page string) []Transaction {
  // Request the HTML page.
  url := "https://etherscan.io/txs?a=" + address + "&p=" + page
  res, err := http.Get(url)
  if err != nil {
    return nil
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    return nil
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    return nil
  }

  result := []Transaction{}

  // Find the review items
  doc.Find("#ContentPlaceHolder1_mainrow tbody tr").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    t := Transaction{}
    t.TxHash = ""
    if s.Find(":nth-child(1) span a").Text() != "" {
      t.TxHash = string([]byte(s.Find(":nth-child(1) span a").Text())[:66])
    } else {
      return
    }
    t.Block = s.Find(":nth-child(2) .hidden-sm").Text()
    band, ok := s.Find(":nth-child(3) span").Attr("title")
    if ok {
      t.Age = band
    }
    t.From = s.Find(":nth-child(4)").Text()
    t.Status = s.Find(":nth-child(5)").Text()
    t.To = s.Find(":nth-child(6)").Text()
    t.Value = s.Find(":nth-child(7)").Text()
    t.TxFee = s.Find(":nth-child(8)").Text()
    result = append(result, t)
  })

  return result
}
