package spider

import (
	"fmt"
	"sync"
)

/////////////////////////////////////////Scheduler////////////////////////////////////////////////

type IStore interface {
	Insert(interface{}) bool
}

type Scheduler struct {
	Crawlers map[string]ICrawler
	Locker   *sync.Mutex
	Queue    chan interface{}
}

func (this *Scheduler) Add(crawler ICrawler) {
	name := crawler.GetName()
	if _, ok := this.Crawlers[name]; ok != true {
		this.Crawlers[name] = crawler
	}
}

func (this *Scheduler) Run() {
	for _, crawler := range this.Crawlers {
		go func(c ICrawler) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			c.Run()
			if c.Valid() {
				fmt.Println("Visit()")
				c.Visit(this.Queue, this.Locker)
			}
		}(crawler)
	}
}

func (this *Scheduler) Beat(store IStore) {
	for {
		select {
		case item := <-this.Queue:
			{
				store.Insert(item)
			}
		}
	}
}
