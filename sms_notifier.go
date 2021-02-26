package inventory_notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SmsNotifier struct {
	Endpoint   string
	AccountID  string
	AuthToken  string
	Sender     string
	Recipients []string
}

func (sms *SmsNotifier) Notify(product *ProductNotification) {
	for _, number := range sms.Recipients {
		sms.sendNotification(product, number)
	}
}

func (sms *SmsNotifier) sendNotification(product *ProductNotification, recipient string) {
	urlStr := fmt.Sprintf("%s/Accounts/%s/Messages.json", sms.Endpoint, sms.AccountID)

	msgData := url.Values{}
	msgData.Set("To", recipient)
	msgData.Set("From", sms.Sender)
	msgData.Set("Body", fmt.Sprintf("%s (%.1f d): %s", product.Name, product.ROI, product.Url))
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(sms.AccountID, sms.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("sms error: %s\n", err.Error())
	}

	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		log.Errorf("sms error: %s\n", err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 300 {
		log.WithFields(log.Fields{
			"card":      product.Name,
			"url":       product.Url,
			"price":     product.SalePrice,
			"maxPrice":  product.MaxPrice,
			"roi":       product.ROI,
			"recipient": recipient,
		}).Info("SMS notifying recipient")
	} else {
		log.WithFields(log.Fields{
			"code":    data["code"],
			"status":  data["status"],
			"message": data["message"],
		}).Error("SMS Notifier error")
	}

}
