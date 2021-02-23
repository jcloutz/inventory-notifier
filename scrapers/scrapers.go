package scrapers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
)

type Scraper interface {
	Scrape(collector *colly.Collector)
}

func MatchProductAndNotify(title string, url string, site string, salePrice float64, matchers *inventory_notifier.MatcherContainer, notifiers *inventory_notifier.Notifiers) {
	cfg, err := matchers.Find(title)
	if err != nil {
		fmt.Printf("no matching product config found for item %s\n", title)
		return
	}

	notifiers.Notify(&inventory_notifier.ProductNotification{
		Name:      cfg.Name,
		Url:       url,
		Site:      site,
		SalePrice: salePrice,
		MaxPrice:  cfg.MaxPrice,
	})
}

func ConvertPrice(priceString string) (float64, error) {
	price := strings.ToLower(strings.TrimSpace(strings.TrimLeft(priceString, "$")))
	if price == "" {
		return 0, errors.New("invalid price")
	}

	curPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return 0, errors.New("price not present")
	}

	return curPrice, nil
}
