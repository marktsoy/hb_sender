package sender

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/marktsoy/hb_sender/internal/models"
)

type Sender struct {
	bot   *tgbotapi.BotAPI
	chats []int64
	mutex sync.Mutex
}

func New(bot *tgbotapi.BotAPI) *Sender {
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	return &Sender{
		bot:   bot,
		chats: make([]int64, 0),
	}
}

func (sender *Sender) ManageChat(s models.Subscription) {
	if s.IsSubscribed {
		sender.AddChat(s.ChatID)
	} else {
		sender.RemoveChat(s.ChatID)
	}
}

func (sender *Sender) AddChat(chatID int64) {
	sender.mutex.Lock()
	sender.chats = append(sender.chats, chatID)
	sender.mutex.Unlock()
}

func (sender *Sender) RemoveChat(chatID int64) {
	sender.mutex.Lock()
	for i, v := range sender.chats {
		if v == chatID {
			sender.chats = append(sender.chats[:i], sender.chats[i+1:]...) // Spread operator ???
			break
		}
	}
	sender.mutex.Unlock()
}

func (sender *Sender) Send(m models.Message) {
	sender.sendText(m.Content)
}

func (sender *Sender) sendText(text string) {
	for _, chatID := range sender.chats {
		msg := tgbotapi.NewMessage(chatID, text)
		sender.bot.Send(msg)
	}
}
