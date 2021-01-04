package main

import (
	"flag"
	"fmt"

	"github.com/BurntSushi/toml"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/marktsoy/hb_sender/internal/app"
	"github.com/marktsoy/hb_sender/internal/sender"
)

var (
	configPath string
)

func init() {
	fmt.Println("Initializing config path")

	// Flag (cmd args) definitions
	flag.StringVar(&configPath, "config-path", "configs/senderconf.toml", "Path to config file")
}

func main() {

	config := app.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BotKey: %v \n", config.BotKey)
	fmt.Printf("Rabbit Addr: %v \n", config.RabbitAddr)
	fmt.Printf("Message Queue: %v \n", config.RabbitMessageQueue)
	fmt.Printf("Read message timout: %v \n", config.RabbitMessageTimeout)

	bot, err := tgbotapi.NewBotAPI(config.BotKey)
	if err != nil {
		panic(err)
	}

	s := sender.New(bot)
	go func() {
		app.ListenToSubscribers(config, s.ManageChat)
	}()
	app.ListenToMessages(config, s.Send)

}
