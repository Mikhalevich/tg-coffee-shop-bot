package config

import "time"

type Config struct {
	LogLevel            string      `yaml:"log_level" required:"true"`
	Tracing             Tracing     `yaml:"tracing" required:"true"`
	Postgres            Postgres    `yaml:"postgres" required:"true"`
	Bot                 Bot         `yaml:"bot" required:"true"`
	ButtonRedis         ButtonRedis `yaml:"button_redis" required:"true"`
	MessageWorker       Worker      `yaml:"message_worker" required:"true"`
	AnswerPaymentWorker Worker      `yaml:"answer_payment_worker" required:"true"`
	InvoiceWorker       Worker      `yaml:"invoice_worker" required:"true"`
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

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type Bot struct {
	Token        string `yaml:"token" required:"true"`
	PaymentToken string `yaml:"payment_token" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr" required:"true"`
	Pwd  string        `yaml:"pwd" required:"true"`
	DB   int           `yaml:"db" required:"true"`
	TTL  time.Duration `yaml:"ttl" required:"true"`
}

type Worker struct {
	Count     int           `yaml:"count" required:"true"`
	Interval  time.Duration `yaml:"interval" required:"true"`
	BatchSize int           `yaml:"batch_size" required:"true"`
}
