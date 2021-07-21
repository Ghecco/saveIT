package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/Ghecco/saveIT/pkg/controllers"
	tb "gopkg.in/tucnak/telebot.v2"
)

const ModeMarkdown tb.ParseMode = "Markdown"

func Init(telegramToken string) {

	b, err := tb.NewBot(tb.Settings{
		URL:    "",
		Token:  telegramToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", func(m *tb.Message) {
		if m.Sender.Username == "" {
			b.Send(m.Sender, "You don't have a username, add it via the telegram settings to start using the bot.\nOnce added, use /start!")
			return
		}
		err, _ := controllers.GetIDByUsername(m.Sender.Username)
		if err != nil {
			fmt.Printf("[Telegram] User %s not found in database.", m.Sender.Username)
			b.Send(m.Sender, fmt.Sprintf("Hi *%s*\nYou're not present in the database", m.Sender.Username), ModeMarkdown)
			return
		}
		b.Send(m.Sender, fmt.Sprintf("Welcome *%s*\nYou're present in the database.", m.Sender.Username), ModeMarkdown)
	})

	b.Start()
}
