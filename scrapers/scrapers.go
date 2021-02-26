package scrapers

import (
	"strings"
	"sync"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

var containerMutex = sync.RWMutex{}
var scraperContainer = map[string]ScraperConfig{
	"newegg": {
		Selector: ".items-view div.item-container",
		Handler: func(e *colly.HTMLElement) (*ScraperResult, error) {
			title := e.ChildText(".item-title")
			productUrl := e.ChildAttr(".item-title", "href")
			unitPrice := e.ChildAttr("input[name='ItemUnitPrice']", "value")
			buttonType := e.ChildAttr(".same-td-elastic button", "title")
			inStock := false

			if buttonType == "ADD TO CART" {
				inStock = true
			}

			price, err := ConvertPrice(unitPrice)
			if err != nil {
				return &ScraperResult{}, err
			}

			return &ScraperResult{
				Title:   title,
				Url:     productUrl,
				Price:   price,
				InStock: inStock,
			}, nil
		},
	},

	"gamestop": {
		Selector: ".row.infinitescroll-results-grid .product-grid-tile-wrapper",
		Handler: func(e *colly.HTMLElement) (*ScraperResult, error) {
			title := e.ChildText(".pd-name")
			productUrl := "https://gamestop.com" + e.ChildAttr(".link-name", "href")
			unitPrice := e.ChildText(".actual-price")
			buttonText := e.ChildText(".add-to-cart")
			inStock := false

			if strings.TrimSpace(strings.ToLower(buttonText)) == "add to cart" {
				inStock = true
			}

			price, err := ConvertPrice(unitPrice)
			if err != nil {
				log.Errorf("price: %s, product: %s", unitPrice, title)
				return &ScraperResult{}, err
			}

			return &ScraperResult{
				Title:   title,
				Url:     productUrl,
				Price:   price,
				InStock: inStock,
			}, nil
		},
	},

	"bestbuy": {
		Selector: ".sku-item-list .sku-item",
		Handler: func(e *colly.HTMLElement) (result *ScraperResult, err error) {
			title := e.ChildText(".sku-header a")
			productUrl := "https://bestbuy.com" + e.ChildAttr(".sku-header a", "href")
			unitPrice := e.ChildText(".priceView-hero-price span[aria-hidden=true]")
			buttonText := e.ChildText(".sku-item-list .sku-item .fulfillment-add-to-cart-button button")
			inStock := false

			if strings.TrimSpace(strings.ToLower(buttonText)) == "add to cart" {
				inStock = true
			}

			price, err := ConvertPrice(unitPrice)
			if err != nil {
				log.Errorf("price: %s, product: %s", unitPrice, title)
				return &ScraperResult{}, err
			}

			return &ScraperResult{
				Title:   title,
				Url:     productUrl,
				Price:   price,
				InStock: inStock,
			}, nil
		},
	},

	"officedepot": {
		Selector: ".sku_item",
		Handler: func(e *colly.HTMLElement) (result *ScraperResult, err error) {
			title := e.ChildText(".desc_text a")
			productUrl := "https://officedepot.com" + e.ChildAttr(".desc_text a", "href")
			unitPrice := e.ChildText(".price_column.right")

			btnDisabled := e.ChildAttr("li.cart input", "disabled")
			inStock := false

			if strings.TrimSpace(strings.ToLower(btnDisabled)) != "disabled" {
				inStock = true
			}

			price, err := ConvertPrice(unitPrice)
			if err != nil {
				log.Errorf("price: %s, product: %s", unitPrice, title)
				return &ScraperResult{}, err
			}

			return &ScraperResult{
				Title:   title,
				Url:     productUrl,
				Price:   price,
				InStock: inStock,
			}, nil
		},
	},
}
