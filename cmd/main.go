package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	inventory_notifier "github.com/jcloutz/inventory-notifier"
	"github.com/jcloutz/inventory-notifier/scrapers"
)

func main() {

	var config inventory_notifier.Config
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(source, &config)

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("Parsed config")

	matcherContainer := inventory_notifier.NewMatchContainer()
	for _, match := range config.Matchers {
		matcherContainer.Add(&match)
	}

	notifiers := inventory_notifier.Notifiers{}

	smsConfig := config.Notifiers.Sms
	if smsConfig != nil {
		sms := inventory_notifier.SmsNotifier{
			Endpoint:   smsConfig.ApiEndpoint,
			AccountID:  smsConfig.AccountID,
			AuthToken:  smsConfig.AuthToken,
			Sender:     smsConfig.Sender,
			Recipients: smsConfig.Recipients,
		}
		notifiers.Add(&sms)
	}

	emailConfig := config.Notifiers.Email
	if emailConfig != nil {
		email := inventory_notifier.EmailNotifier{
			Username:   emailConfig.Sender,
			Password:   emailConfig.Password,
			Server:     emailConfig.Server,
			Port:       emailConfig.Port,
			Recipients: emailConfig.Recipients,
		}
		notifiers.Add(&email)
	}

	discordConfig := config.Notifiers.Discord
	if discordConfig != nil {
		discord := inventory_notifier.DiscordNotifier{
			Webook:     discordConfig.Callback,
			Recipients: discordConfig.Recipients,
		}
		notifiers.Add(&discord)
	}

	var siteScrapers []scrapers.Scraper
	for _, siteConfig := range config.Sites {
		siteScrapers = append(siteScrapers, scrapers.ScraperRunner{
			Config:   siteConfig,
			Notifier: &notifiers,
			Matchers: matcherContainer,
		})
	}

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	for _, scraper := range siteScrapers {
		scraper.Scrape()
	}

	for {
		select {
		case <-done:
			fmt.Println("Shutting down")
			return
		}
	}

}
