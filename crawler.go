package spider

import (
	"fmt"
	"strings"
	"sync"

	. "github.com/PuerkitoBio/goquery"
)

type ICrawler interface {
	GetName() string
	GetDoc() *Document
	GetType() string
	Valid() bool
	Visit(chan interface{}, *sync.Mutex)
	Run()
}

type CrawlerConfig struct {
	Url, Type, Host string
	IsGBK           bool
}

type CrawlerBase struct {
	Name string
	Cfg  *CrawlerConfig
	doc  *Document
}

func (this *CrawlerBase) GetName() string {
	return this.Name
}

func (this *CrawlerBase) GetDoc() *Document {
	return this.doc
}

func (this *CrawlerBase) GetType() string {
	return this.Cfg.Type
}

func (this *CrawlerBase) Valid() bool {
	doc := this.GetDoc()
	if doc == nil {
		return false
	}

	return true
}

func (this *CrawlerBase) Visit(Q chan interface{}, L *sync.Mutex) {
	//do sth!
	fmt.Printf("visit: ========%s (%s)===========", this.GetName(), this.Cfg.Url)
	if this.doc == nil {
		fmt.Printf("[-] http request faild, skip (%s)!", this.Cfg.Url)
		return
	}

	con, _ := this.GetDoc().Find("meta").Eq(0).Attr("content")
	con = strings.ToUpper(con)

	if strings.Contains(con, "GBK") || strings.Contains(con, "GB2312") {
		this.Cfg.IsGBK = true
	}
}

func (this *CrawlerBase) Run() {
	var doc *Document
	var err error
	if doc, err = NewDocument(this.Cfg.Url); err != nil {
		this.doc = nil
		return
	}

	this.doc = doc
}
