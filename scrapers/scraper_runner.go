package scrapers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
)

type ScraperRunner struct {
	Config   inventory_notifier.SiteConfig
	Notifier *inventory_notifier.Notifiers
	Matchers *inventory_notifier.MatcherContainer
	ticker   *time.Ticker
	scraper  ScraperConfig
	running  chan bool
}

func (ns ScraperRunner) Scrape() {
	ns.ticker = time.NewTicker(ns.Config.Interval * time.Second)

	log.Infof("initializing %s", ns.Config.Name)
	ns.run()

	go func() {
		for {
			select {
			case <-ns.ticker.C:
				ns.run()
			case <-ns.running:
				ns.ticker.Stop()
			}
		}
	}()
}

func (ns ScraperRunner) run() {
	log.Infof("Checking %s", ns.Config.Name)
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	extensions.RandomUserAgent(c)

	c.Limit(&colly.LimitRule{
		DomainGlob:  fmt.Sprintf("*%s.*", ns.Config.Domain),
		Delay:       5 * time.Second,
		RandomDelay: 1 * time.Second,
		Parallelism: 0,
	})

	ns.RunQueue(c)
}

func (ns ScraperRunner) RunQueue(collector *colly.Collector) {
	scraper, err := GetScraper(ns.Config.Name)
	if err != nil {
		log.Error(err)

		return
	}

	wg := sync.WaitGroup{}

	collector.OnRequest(func(r *colly.Request) {
		log.Infof("Requesting: %s\n", r.URL)

		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	collector.OnResponse(func(response *colly.Response) {
		log.Infof("Received: %s", response.Request.URL.String())
		//log.Info(string(response.Body))
	})

	collector.OnError(func(response *colly.Response, err error) {
		log.Error(err)
	})

	collector.OnHTML(scraper.Selector, func(e *colly.HTMLElement) {
		result, err := scraper.Handler(e)
		if err != nil {
			log.Error(err)
		}

		if result.InStock {
			go func() {
				wg.Add(1)
				MatchProductAndNotify(result.Title, result.Url, ns.Config.Name, result.Price, ns.Matchers, ns.Notifier)
				wg.Done()
			}()
		}
	})

	for _, page := range ns.Config.Pages {
		collector.Visit(page)
	}

	wg.Wait()
}
