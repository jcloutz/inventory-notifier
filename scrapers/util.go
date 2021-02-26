package scrapers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
)

func AddScraper(name string, cfg ScraperConfig) {
	containerMutex.Lock()
	defer containerMutex.Unlock()

	scraperContainer[name] = cfg
}

func GetScraper(name string) (ScraperConfig, error) {
	containerMutex.RLock()
	defer containerMutex.RUnlock()

	if val, ok := scraperContainer[name]; ok {
		return val, nil
	}
	return ScraperConfig{}, errors.New("no scraper found")
}

type Scraper interface {
	Scrape()
}

type ScraperConfig struct {
	Selector string
	Handler  func(element *colly.HTMLElement) (result *ScraperResult, err error)
}

type ScraperResult struct {
	Title   string
	Url     string
	Price   float64
	InStock bool
}

func MatchProductAndNotify(title string, url string, site string, salePrice float64, matchers *inventory_notifier.MatcherContainer, notifiers *inventory_notifier.Notifiers) {
	cfg, err := matchers.Find(title)
	if err != nil {
		log.Errorf("no matching product config found for item %s\n", title)
		return
	}

	if salePrice <= cfg.Earns*95 {
		roi := salePrice / cfg.Earns
		log.Infof("%s in stock, roi=%.2f days", cfg.Name, roi)
		notifiers.Notify(&inventory_notifier.ProductNotification{
			Name:      cfg.Name,
			Url:       url,
			Site:      site,
			SalePrice: salePrice,
			MaxPrice:  cfg.MaxPrice,
			ROI:       roi,
		})
	} else {
		log.Infof("%s in stock, but exceeds price threshold", cfg.Name)
	}

}

func ConvertPrice(priceString string) (float64, error) {
	price := strings.ToLower(strings.TrimSpace(strings.TrimLeft(priceString, "$")))
	price = strings.Replace(price, ",", "", -1)
	if price == "" {
		return 0, errors.New(fmt.Sprintf("invalid price: '%s'", price))
	}

	curPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, errors.New("price not present")
	}

	return curPrice, nil
}
