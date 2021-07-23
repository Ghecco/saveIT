package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/Ghecco/saveIT/pkg/controllers"
	tb "gopkg.in/tucnak/telebot.v2"
)

var Sessions []Session

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
	// Handles

	// Command /start
	b.Handle("/start", func(m *tb.Message) {
		if m.Sender.Username == "" {
			b.Send(m.Sender, "You don't have a username, add it via the telegram settings to start using the bot.\nOnce added, use /start!")
			return
		}
		err, _ := controllers.GetIDByUsername(m.Sender.Username)
		if err != nil {
			fmt.Printf("[Telegram] User %s not found in database.", m.Sender.Username)
			b.Send(m.Sender, fmt.Sprintf("Hi *%s*\nYou're not present in the database", m.Sender.Username), ModeMarkdown)
			b.Send(m.Sender, "Please register!\nEnter a password of up to 24 characters, you can also leave it blank (not recommended)")
			Sessions = append(Sessions, Session{ChatID: m.Sender.ID, Status: STATUS_REGISTER_PASSWORD})
			return
		}
		b.Send(m.Sender, fmt.Sprintf("Welcome *%s*\nYou're present in the database.", m.Sender.Username), ModeMarkdown)
		Sessions = append(Sessions, Session{ChatID: m.Sender.ID, Status: STATUS_LOGIN_PASSWORD})
	})
	// When user send text, it is managed by the Session status
	b.Handle(tb.OnText, func(m *tb.Message) {
		var sessionID, status int
		sessionID = GetSessionID(m.Sender.ID)
		if sessionID == -1 {
			fmt.Printf("chat id %d does not have a session started", m.Sender.ID)
			b.Delete(m)
			return
		}
		status = Sessions[sessionID].Status
		switch status {
		case STATUS_NONE:
			b.Delete(m)
		case STATUS_REGISTER_PASSWORD:
			if !controllers.AddUser(m.Sender.Username, m.Text) {
				fmt.Printf("%s has not entered a valid password.\n", m.Sender.Username)
				b.Send(m.Sender, "Invalid password, maximum 24 characters.\nTry again.")
				return
			}

			Sessions[sessionID].IsLogged = true
			b.Send(m.Sender, "Registration successful!\nYou will be logged in automatically")
			return // temp

		case STATUS_LOGIN_PASSWORD:
			if !controllers.LoginUser(m.Sender.Username, m.Text) {
				fmt.Printf("%s has not entered the valid password.\n", m.Sender.Username)
				b.Send(m.Sender, "Incorrect password, try again")
			}
			if DestroyOtherSession(m.Sender.Username) {
				fmt.Printf("Destroyed another session ID for the user %s", m.Sender.Username)
			}

			userSessionID := GetSessionID(m.Sender.ID)

			Sessions[userSessionID].Username = m.Sender.Username
			Sessions[userSessionID].IsLogged = true

			fmt.Printf("User %s logged", Sessions[userSessionID].Username)
			b.Send(m.Sender, "You are logged in successfully !")
			return
		}
		b.Start()
	})
}
