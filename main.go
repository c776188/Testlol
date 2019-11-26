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
	for i := 1; i <= 20; i++ {
		isContinue := crawler_lol_self_info(url_prefix, i, url_suffix, make(chan int, 1))
		if !isContinue {
			break
		}
	}
}

// 抓個人的下一頁資料
func crawler_lol_self_info(url_prefix string, page int, url_suffix string, queue chan int) bool {
	// 睡避免爬太快
	time.Sleep(5 * time.Second)

	// 連線逾時
	visited[url_prefix+strconv.Itoa(page)+url_suffix] = true
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	// 顯示頁數 不然不知道自己在哪
	fmt.Println("--------------------------------page = ", page)

	req, _ := http.NewRequest("GET", url_prefix+strconv.Itoa(page)+url_suffix, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return false
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	isResult := true
	doc.Find("a[class=text-white]").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		href, ok := selection.Attr("href")
		if !ok {
			fmt.Println("error")
		}
		// fmt.Println("https://lol.moa.tw" + href)
		isDetailSuccess := crawler_lol_detail_info("https://lol.moa.tw"+href, make(chan int, 1))
		if isDetailSuccess == false {
			isResult = false
			return false
		}

		return true
	})

	return isResult
}

// 抓該場戰績資料
func crawler_lol_detail_info(url string, queue chan int) bool {
	// 睡避免爬太快
	time.Sleep(2 * time.Second)

	// 連線逾時
	visited[url] = true
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	req, _ := http.NewRequest("GET", url, nil)
	// 自定义Header
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http get error", err)
		return false
	}
	//函数结束后关闭相关链接
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// 抓人名有無等於XXXXX
	isSearch := false
	doc.Find("span[class=sumtip]>a").Each(func(i int, selection *goquery.Selection) {
		// fmt.Println(selection.Text())
		if selection.Text() == os.Getenv("LOL_NAME") {
			isSearch = true
		}
	})

	// 沒抓到顯示該場網址並且中斷crawler
	if !isSearch {
		fmt.Println("fail url: ", url)
		return false
	}
	return true
}
