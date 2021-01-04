package app

// Config ...
type Config struct {
	BotKey                string `toml:"tg_bot_public_key"`
	RabbitAddr            string `toml:"rabbit_addr"`
	RabbitMessageTimeout  int32  `toml:"rabbit_message_timeout"`
	RabbitMessageQueue    string `toml:"rabbit_message_queue"`
	RabbitSubscriberQueue string `toml:"rabbit_subscriber_queue"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
