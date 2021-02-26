package inventory_notifier

import (
	"time"
)

type Config struct {
	Interval  float32         `yaml:"interval"`
	Notifiers NotifierConfigs `yaml:"notifiers"`
	Matchers  []ProductConfig `yaml:"matchers"`
	Sites     []SiteConfig    `yaml:"sites"`
}

type ProductConfig struct {
	Name     string  `yaml:"name"`
	MaxPrice float64 `yaml:"max_price"`
	Earns    float64 `yaml:"earns"`
}

type SiteConfig struct {
	Name     string        `yaml:"name"`
	Domain   string        `yaml:"domain"`
	Interval time.Duration `yaml:"interval"`
	Pages    []string      `yaml:"pages"`
}

type DiscordConfig struct {
	Callback   string   `yaml:"callback"`
	Recipients []string `yaml:"recipients"`
}

type EmailConfig struct {
	Sender     string   `yaml:"sender"`
	Recipients []string `yaml:"recipients"`
	Server     string   `yaml:"server"`
	Port       int      `yaml:"port"`
	Password   string   `yaml:"password"`
}

type SmsConfig struct {
	ApiEndpoint string   `yaml:"api_endpoint"`
	AccountID   string   `yaml:"account_id"`
	AuthToken   string   `yaml:"auth_token"`
	Sender      string   `yaml:"sender"`
	Recipients  []string `yaml:"recipients"`
}

type NotifierConfigs struct {
	Discord *DiscordConfig `yaml:"discord"`
	Email   *EmailConfig   `yaml:"email"`
	Sms     *SmsConfig     `yaml:"sms"`
}
