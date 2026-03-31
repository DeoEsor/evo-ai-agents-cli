package config

type IAMConfig struct {
	BffHost      string `env:"BFF_HOST" envDefault:""`
	ClientID     string `env:"CLIENT_ID" envDefault:""`
	ClientSecret string `env:"CLIENT_SECRET" envDefault:""`
}
