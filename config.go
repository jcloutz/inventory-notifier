package inventory_notifier

type Config struct {
	Interval  float32               `yaml:"interval"`
	Notifiers NotifierConfigs       `yaml:"notifiers"`
	Matchers  []ProductConfig       `yaml:"matchers"`
	Sites     map[string]SiteConfig `yaml:"sites"`
}

type ProductConfig struct {
	Name     string  `yaml:"name"`
	MaxPrice float64 `yaml:"max_price"`
}

type SiteConfig struct {
	Page string `yaml:"page"`
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
	Discord DiscordConfig `yaml:"discord"`
	Email   EmailConfig   `yaml:"email"`
	Sms     SmsConfig     `yaml:"sms"`
}
