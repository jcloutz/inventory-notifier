package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
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
	src := []byte(source)

	yaml.Unmarshal(src, &config)

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("Parsed config")

	matcherContainer := inventory_notifier.NewMatchContainer()
	for _, match := range config.Matchers {
		matcherContainer.Add(&match)
	}

	notifiers := inventory_notifier.Notifiers{}

	smsConfig := config.Notifiers.Sms
	sms := inventory_notifier.SmsNotifier{
		Endpoint:   smsConfig.ApiEndpoint,
		AccountID:  smsConfig.AccountID,
		AuthToken:  smsConfig.AuthToken,
		Sender:     smsConfig.Sender,
		Recipients: smsConfig.Recipients,
	}
	notifiers.Add(&sms)

	emailConfig := config.Notifiers.Email
	email := inventory_notifier.EmailNotifier{
		Username:   emailConfig.Sender,
		Password:   emailConfig.Password,
		Server:     emailConfig.Server,
		Port:       emailConfig.Port,
		Recipients: emailConfig.Recipients,
	}
	notifiers.Add(&email)

	discordConfig := config.Notifiers.Discord
	discord := inventory_notifier.DiscordNotifier{
		Webook:     discordConfig.Callback,
		Recipients: discordConfig.Recipients,
	}
	notifiers.Add(&discord)

	var siteScrapers []scrapers.Scraper
	for site, siteConfig := range config.Sites {
		switch site {
		case "newegg":
			siteScrapers = append(siteScrapers, scrapers.NeweggScraper{
				Url:      siteConfig.Page,
				Notifier: &notifiers,
				Matchers: matcherContainer,
			})
		case "b_and_h":
			siteScrapers = append(siteScrapers, scrapers.BandH{
				Url:      siteConfig.Page,
				Notifier: &notifiers,
				Matchers: matcherContainer,
			})
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	tickerChan := ticker.C
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	// configure colly
	c := colly.NewCollector(
		//colly.Async(true),
		colly.AllowURLRevisit(),
	)
	extensions.RandomUserAgent(c)

	//siteScrapers[0].Scrape(c)

	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-tickerChan:
			fmt.Println("Firing")

			for _, scraper := range siteScrapers {
				go func() {
					col := c.Clone()
					extensions.RandomUserAgent(col)
					scraper.Scrape(col)
				}()
			}

			v := config.Interval + rand.Float32()*6
			ticker.Reset(time.Duration(v) * time.Second)
		case <-done:
			fmt.Println("Shutting down")
			return
		}
	}

}
