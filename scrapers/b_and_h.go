package scrapers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
)

type BandH struct {
	Url      string
	Notifier *inventory_notifier.Notifiers
	Matchers *inventory_notifier.MatcherContainer
	ticker   time.Ticker
}

func (ns BandH) Scrape(collector *colly.Collector) {
	counter := 0
	wg := sync.WaitGroup{}
	collector.OnHTML("[class^=\"product_\"]", func(e *colly.HTMLElement) {
		counter++
		title := e.ChildText("[data-selenium='miniProductPageProductName']")
		fmt.Println("title", title)
		//title := e.ChildText(".item-title")
		//productUrl := e.ChildAttr(".item-title", "href")
		//unitPrice := e.ChildAttr("input[name='ItemUnitPrice']", "value")
		//buttonType := e.ChildAttr(".same-td-elastic button", "title")
		//fmt.Printf("[Newegg][%d] Checking: %s\n", counter, title)
		//if buttonType == "ADD TO CART" {
		//	price, err := ConvertPrice(unitPrice)
		//	if err != nil {
		//		log.WithFields(log.Fields{
		//			"title":      title,
		//			"productUrl": productUrl,
		//			"unitPrice":  unitPrice,
		//		}).Errorf("[%d] unable to parse price", counter)
		//		return
		//	}
		//
		//	go func() {
		//		wg.Add(1)
		//		MatchProductAndNotify(title, productUrl, "Newegg", price, ns.Matchers, ns.Notifier)
		//		wg.Done()
		//	}()
		//	log.WithFields(log.Fields{
		//		"title":      title,
		//		"productUrl": productUrl,
		//		"unitPrice":  unitPrice,
		//	}).Infof("[%d] Newegg in stock", counter)
		//} else {
		//	log.WithFields(log.Fields{
		//		"title":      title,
		//		"productUrl": productUrl,
		//		"unitPrice":  unitPrice,
		//	}).Warnf("[%d] Newegg out of stock", counter)
		//}
	})

	collector.OnResponse(func(response *colly.Response) {
		fmt.Println(response.StatusCode)
	})
	collector.OnError(func(response *colly.Response, err error) {
		log.Error(err)
	})
	collector.OnRequest(func(r *colly.Request) {
		log.Info("Visiting", r.URL)

		//“Accept”: “text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,
		//
		//	image/apng,*/*;q=0.8,application/signed-exchange;v=b3″,
		//
		//“Accept-Encoding”: “gzip”,
		//
		//“Accept-Language”: “en-US,en;q=0.9,es;q=0.8”,
		//
		//“Upgrade-Insecure-Requests”: “1”,
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	err := collector.Visit(ns.Url)
	if err != nil {
		log.Errorf("B&H Scraper: %s", err)
	}

	//collector.Wait()
	wg.Wait()
	fmt.Println("Done")
}
