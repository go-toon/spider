package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	. "github.com/PuerkitoBio/goquery"
	"github.com/go-toon/spider"
	"gopkg.in/mgo.v2"
)

var CACHEMAP map[string]bool = make(map[string]bool)

type Joke struct {
	Title, Content, Src, ImgSrc, ImgLocal, ThumbSrc, ThumbLocal   string
	Id, Date, Likes, ImgWidth, ImgHeight, ThumbWidth, ThumbHeight int64
}

/////////////////////////////////////百思不得其解（文字）////////////////////////////////////////////
// http://www.budejie.com/text/

type budejieCrawler struct {
	spider.CrawlerBase
}

func (this *budejieCrawler) Visit(Q chan interface{}, L *sync.Mutex) {
	this.CrawlerBase.Visit(Q, L)

	this.GetDoc().Find(".j-r-list").Each(func(i int, s *Selection) {
		s.Find("li").Each(func(j int, ls *Selection) {
			src, ok := ls.Find(".j-list-comment").Attr("href")
			if !ok {
				src = this.GetDoc().Url.String()
			} else {
				src = this.Cfg.Host + src
			}
			fmt.Println(src)
			cont := strings.TrimSpace(ls.Find(".j-r-list-c-desc").Text())
			tm := strings.TrimSpace(ls.Find("u-time").Text())
			tmUnix, _ := time.Parse("2006-01-02 15:04:05", tm)

			if cont == "" {
				return
			}
			_, ok = CACHEMAP[src]
			fmt.Println(ok)
			if ok {
				fmt.Println("[!] 已经导入过")
				return
			} else {
				L.Lock()
				CACHEMAP[src] = true
				L.Unlock()

				fmt.Println(cont)
				joke := Joke{
					Content: cont,
					Src:     src,
					Date:    tmUnix.Unix(),
				}

				Q <- joke
			}

		})
	})
}

type mongoBeat struct {
	session *mgo.Session
}

func (this *mongoBeat) Insert(item interface{}) bool {
	switch v := interface{}(item).(type) {
	case Joke:
		fmt.Println(v.Content)
		err := this.session.DB("mytoon").C("myjokes").Insert(v)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("sth wrong!")
	}

	return true
}

func main() {
	fmt.Println("开始抓取笑话数据...")

	myScheduler := &spider.Scheduler{
		Crawlers: make(map[string]spider.ICrawler),
		Locker:   new(sync.Mutex),
		Queue:    make(chan interface{}, 10),
	}

	myScheduler.Add(&budejieCrawler{
		spider.CrawlerBase{
			Name: "百思不得姐－内涵段子",
			Cfg: &spider.CrawlerConfig{
				Host: "http://www.budejie.com",
				Url:  "http://www.budejie.com/text",
			},
		},
	})

	go func() {
		fmt.Println("ccccc")
		for {
			fmt.Println("dfsdfsafasf")
			myScheduler.Run()
			fmt.Println("[+] Zsss")
			time.Sleep(time.Second * 5)
		}
	}()

	mgoDialInfo := &mgo.DialInfo{
		Addrs:    []string{"112.126.77.185:27017"},
		Timeout:  60 * time.Second,
		Database: "mytoon",
		Username: "mytoon",
		Password: "123456",
	}

	//	mgoSession, err := mgo.Dial("172.28.50.179:28010")
	mgoSession, err := mgo.DialWithInfo(mgoDialInfo)
	if err != nil {
		panic(err)
	}
	mgoSession.SetMode(mgo.Monotonic, true)

	defer mgoSession.Close()

	myScheduler.Beat(&mongoBeat{
		session: mgoSession,
	})

}
