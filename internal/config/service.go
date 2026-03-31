package config

type AppEnvironment string

const (
	Dev   AppEnvironment = "dev"
	Stage AppEnvironment = "stage"
	Prod  AppEnvironment = "prod"
)

type ServiceConfig struct {
	AppEnvironment AppEnvironment `env:"APP_ENVIRONMENT" envDefault:"dev"`
	ServiceName    string         `env:"NAME" envDefault:"ai-agents-cli"`
	LogLevel       string         `env:"LOG_LEVEL" envDefault:"debug"`
	Version        string         `env:"VERSION" envDefault:"1.0.0"`
}
