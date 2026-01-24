package config

type Config struct {
	LogLevel string   `yaml:"log_level" required:"true"`
	Tracing  Tracing  `yaml:"tracing" required:"true"`
	Bot      Bot      `yaml:"bot" required:"true"`
	Postgres Postgres `yaml:"postgres" required:"true"`
	HTTPPort int      `yaml:"http_port" required:"true"`
}

func (c *Config) Level() string {
	return c.LogLevel
}

func (c *Config) ServiceName() string {
	return c.Tracing.ServiceName
}

func (c *Config) TracingEndpoint() string {
	return c.Tracing.Endpoint
}

type Bot struct {
	Token        string `yaml:"token" required:"true"`
	PaymentToken string `yaml:"payment_token" required:"true"`
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}
