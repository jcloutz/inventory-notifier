package inventory_notifier

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type DiscordNotifier struct {
	Webook     string
	Recipients []string
}

func (dn *DiscordNotifier) Notify(product *ProductNotification) {
	recipients := make([]string, len(dn.Recipients))
	for i, r := range dn.Recipients {
		recipients[i] = fmt.Sprintf("@%s", r)
	}

	log.WithFields(log.Fields{
		"card":      product.Name,
		"url":       product.Url,
		"price":     product.SalePrice,
		"maxPrice":  product.MaxPrice,
		"roi":       product.ROI,
		"recipient": recipients,
	}).Info("Discord notifying recipients")

	body := []byte(fmt.Sprintf(`{
		"content": "<%s>",
		"username": "Inventory Notifier",
		"embeds": [
			{ "title": "%s In Stock", "description": "Price: $%.2f\nROI: %.2f days\n%s"}
		]
	}`, strings.Join(recipients, ", "), product.Name, product.SalePrice, product.ROI, product.Url))

	_, err := http.Post(dn.Webook, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Error(err)
	}
}
