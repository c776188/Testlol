package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var visited = make(map[string]bool)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	url_prefix := "https://lol.moa.tw/Ajax/recentgames_more2/" + os.Getenv("LOL_ID") + "/page:"
	url_suffix := "/sort:GameMatch.createDate/direction:desc"
	queue := make(chan int, 1)
	go func() {
		queue <- 1
	}()
	for i := 1; i <= 10; i++ {
		download(url_prefix, i, url_suffix, make(chan int, 1))
	}
}

func download(url_prefix string, page int, url_suffix string, queue chan int) {
	visited[url_prefix+strconv.Itoa(page)+url_suffix] = true
	timeout := time.Duration(10 * time.Second)

	client := &http.Client{
		Timeout: timeout,
	}

	fmt.Println("--------------------------------page = ", page)

	req, _ := http.NewRequest("GET", url_prefix+strconv.Itoa(page)+url_suffix, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	doc.Find("th>div[class=pull-right]").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})

	// link := url_prefix + queue + url_suffix
	// go func() {
	// 	queue <- queue + 1
	// }()
}
