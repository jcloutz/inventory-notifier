package scrapers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
)

type NeweggScraper struct {
	Url      string
	Notifier *inventory_notifier.Notifiers
	Matchers *inventory_notifier.MatcherContainer
	ticker   time.Ticker
}

func (ns NeweggScraper) Scrape(collector *colly.Collector) {
	counter := 0
	wg := sync.WaitGroup{}
	collector.OnHTML(".items-view div.item-container", func(e *colly.HTMLElement) {
		counter++
		title := e.ChildText(".item-title")
		productUrl := e.ChildAttr(".item-title", "href")
		unitPrice := e.ChildAttr("input[name='ItemUnitPrice']", "value")
		buttonType := e.ChildAttr(".same-td-elastic button", "title")
		fmt.Printf("[Newegg][%d] Checking: %s\n", counter, title)
		if buttonType == "ADD TO CART" {
			price, err := ConvertPrice(unitPrice)
			if err != nil {
				log.WithFields(log.Fields{
					"title":      title,
					"productUrl": productUrl,
					"unitPrice":  unitPrice,
				}).Errorf("[%d] unable to parse price", counter)
				return
			}

			go func() {
				wg.Add(1)
				MatchProductAndNotify(title, productUrl, "Newegg", price, ns.Matchers, ns.Notifier)
				wg.Done()
			}()
			log.WithFields(log.Fields{
				"title":      title,
				"productUrl": productUrl,
				"unitPrice":  unitPrice,
			}).Infof("[%d] Newegg in stock", counter)
		} else {
			log.WithFields(log.Fields{
				"title":      title,
				"productUrl": productUrl,
				"unitPrice":  unitPrice,
			}).Warnf("[%d] Newegg out of stock", counter)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := collector.Visit(ns.Url)

	collector.Wait()
	wg.Wait()
	fmt.Println("Done")
	if err != nil {
		fmt.Println(err)
	}
}
