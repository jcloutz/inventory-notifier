package inventory_notifier

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailNotifier struct {
	Username   string
	Password   string
	Server     string
	Port       int
	Recipients []string
}

func (em *EmailNotifier) Notify(product *ProductNotification) {
	m := gomail.NewMessage()
	m.SetHeader("From", em.Username)
	m.SetHeader("To", em.Recipients...)
	m.SetHeader("Subject", fmt.Sprintf("%s In Stock", product.Name))
	m.SetBody("text/html", fmt.Sprintf(`
		<b>%s in stock:</b><br >
		<b>Price:</br> $%.2f<br />
		<b>ROI</b> %.1f days<br />
		<a href="%s">%s</a>`,
		product.Name, product.SalePrice, product.ROI, product.Url, product.Url))

	d := gomail.NewDialer(em.Server, em.Port, em.Username, em.Password)

	log.WithFields(log.Fields{
		"card":       product.Name,
		"url":        product.Url,
		"price":      product.SalePrice,
		"maxPrice":   product.MaxPrice,
		"roi":        product.ROI,
		"recipients": em.Recipients,
	}).Info("Email notifying recipients")

	if err := d.DialAndSend(m); err != nil {
		log.Error("email notifier: %s", err.Error())
	}
}
