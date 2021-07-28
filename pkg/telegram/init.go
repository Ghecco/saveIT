package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	var (
		//home
		menuHome = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		// reply buttons
		btnIdea  = menuHome.Text("🖊 Idea")
		btnMedia = menuHome.Text("💿 Media")

		// Archive
		menuHomeArchive = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		// Reply Buttons
		btnAddMedia   = menuHomeArchive.Text("➕ Add Media")
		btnShowMedias = menuHomeArchive.Text("👁‍🗨 Show my Medias")

		// Ideas

		menuHomeIdea = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		menuBack     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		// Reply buttons.
		btnAddIdea    = menuHomeIdea.Text("➕ Add Idea")
		btnRemoveIdea = menuHomeIdea.Text("➖ Remove Idea")
		btnShowIdeas  = menuHomeIdea.Text("👁‍🗨 Show my Ideas")

		btnBack = menuBack.Text("❌ Back")

		btnHome = menuBack.Text("🔝 Home")
	)
	menuHome.Reply(
		menuHome.Row(btnIdea),
		menuHome.Row(btnMedia),
	)

	menuHomeArchive.Reply(
		menuHomeArchive.Row(btnAddMedia),
		menuHomeArchive.Row(btnShowMedias),
		menuHomeArchive.Row(btnHome),
	)

	menuHomeIdea.Reply(
		menuHomeIdea.Row(btnAddIdea),
		menuHomeIdea.Row(btnRemoveIdea),
		menuHomeIdea.Row(btnShowIdeas),
		menuHomeIdea.Row(btnHome),
	)
	menuBack.Reply(
		menuBack.Row(btnBack),
	)
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
		b.Send(m.Sender, "Please login!\nEnter your password !")
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
			return
		case STATUS_REGISTER_PASSWORD:
			if !controllers.AddUser(m.Sender.Username, m.Text) {
				fmt.Printf("%s has not entered a valid password.\n", m.Sender.Username)
				b.Send(m.Sender, "Invalid password, maximum 24 characters.\nTry again.")
				return
			}
			Sessions[sessionID].Username = m.Sender.Username
			Sessions[sessionID].IsLogged = true
			Sessions[sessionID].Status = STATUS_NONE
			b.Send(m.Sender, "Registration successful!\nYou will be logged in automatically", menuHomeIdea)
			return // temp

		case STATUS_LOGIN_PASSWORD:
			if !controllers.LoginUser(m.Sender.Username, m.Text) {
				fmt.Printf("%s has not entered the valid password.\n", m.Sender.Username)
				b.Send(m.Sender, "Incorrect password, try again")
				return
			}
			if DestroyOtherSession(m.Sender.Username) {
				fmt.Printf("Destroyed another session ID for the user %s", m.Sender.Username)
			}

			userSessionID := GetSessionID(m.Sender.ID)

			Sessions[userSessionID].Username = m.Sender.Username
			Sessions[userSessionID].IsLogged = true
			Sessions[userSessionID].Status = STATUS_NONE

			fmt.Printf("User %s logged", Sessions[userSessionID].Username)
			b.Send(m.Sender, "You are logged in successfully !")
			b.Send(m.Sender, "What do you want to save?", menuHome)
			return
		case STATUS_ADD_IDEA:
			var title, content string
			indexEndTitle := strings.Index(m.Text, ",")
			if indexEndTitle != -1 {
				title = m.Text[0:indexEndTitle]
				content = m.Text[indexEndTitle+1 : len(m.Text)]
			} else {
				title = m.Text
			}
			if len(title) < 4 || len(title) > 50 {
				b.Send(m.Sender, "Please, Insert title with 4-50 characters:")
				return
			}
			sessionID := GetSessionID(m.Sender.ID)
			_, userID := controllers.GetIDByUsername(Sessions[sessionID].Username)
			ok := controllers.AddIdea(userID, title, content)
			fmt.Printf("Title:%s Content:%s", title, content)
			if ok {
				b.Send(m.Sender, fmt.Sprintf("Title:*%s*\nContent:%s", title, content), ModeMarkdown, menuHomeIdea)
				Sessions[sessionID].Status = STATUS_NONE
			}
			return
		case STATUS_REMOVE_IDEA:
			val, err := strconv.Atoi(m.Text)
			if err != nil {
				b.Send(m.Sender, "Invalid text, insert one number of the list for remove an idea!")
				return
			}
			sessionID := GetSessionID(m.Sender.ID)
			_, userID := controllers.GetIDByUsername(Sessions[sessionID].Username)
			_, ideas := controllers.GetUserIdeas(userID)

			if val < 1 || val > len(ideas) {
				b.Send(m.Sender, fmt.Sprintf("Invalid number, insert one number 1-%d for remove an idea:", len(ideas)))
				return
			}
			ideaID := uint64(ideas[val-1].ID)
			bool := controllers.RemoveIdea(ideaID)
			if bool {
				b.Send(m.Sender, fmt.Sprintf("Removed idea %d.", val), menuHomeIdea)
				Sessions[sessionID].Status = STATUS_NONE
			}
			return
		case STATUS_ADD_MEDIA:
			b.Delete(m)
			return
		}

	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		if Sessions[GetSessionID(m.Sender.ID)].Status != STATUS_ADD_MEDIA {
			b.Delete(m)
			return
		}
		ok := controllers.AddPicture(m.Sender.Username, m.Photo.MediaFile(), b) // temp
		if ok {
			b.Send(m.Sender, "✅ Added new picture.")
		} else {
			b.Send(m.Sender, "Error: Picture not added")
		}
	})
	b.Handle(tb.OnVideo, func(m *tb.Message) {
		if Sessions[GetSessionID(m.Sender.ID)].Status != STATUS_ADD_MEDIA {
			b.Delete(m)
			return
		}
		ok := controllers.AddVideo(m.Sender.Username, m.Video.MediaFile(), b) // temp
		if ok {
			b.Send(m.Sender, "✅ Added new video.")
		} else {
			b.Send(m.Sender, "Error: Video not added")
		}
	})
	//buttons
	b.Handle(&btnAddIdea, func(m *tb.Message) {
		b.Send(m.Sender, "Insert the new idea Title(4-50)\nAnd down write it content\nExample:\nBuy new car,Buy the new lamborghini", menuBack)
		Sessions[GetSessionID(m.Sender.ID)].Status = STATUS_ADD_IDEA
	})
	b.Handle(&btnRemoveIdea, func(m *tb.Message) {

		sessionID := GetSessionID(m.Sender.ID)
		err, uID := controllers.GetIDByUsername(Sessions[sessionID].Username)
		if err != nil {
			fmt.Printf("User %s not present in database. (GetIDByUsername func)", m.Sender.Username)
			return
		}

		err, ideas := controllers.GetUserIdeas(uID)
		if err != nil {
			b.Send(m.Sender, "You have no saved ideas to remove", menuHomeIdea)
			Sessions[GetSessionID(m.Sender.ID)].Status = STATUS_NONE
			return
		}
		var message, localMessage string
		message = " "
		for i, v := range ideas {
			localMessage = fmt.Sprintf("\n ▪️ %d %s", i+1, v.Title)
			message += localMessage
		}
		b.Send(m.Sender, message, menuBack)
		Sessions[sessionID].Status = STATUS_REMOVE_IDEA

	})
	b.Handle(&btnShowIdeas, func(m *tb.Message) {
		err, uID := controllers.GetIDByUsername(Sessions[GetSessionID(m.Sender.ID)].Username)
		if err != nil {
			fmt.Printf("User %s not present in database. (GetIDByUsername func)", m.Sender.Username)
			return
		}

		err, ideas := controllers.GetUserIdeas(uID)
		if err != nil {
			b.Send(m.Sender, "You have no saved ideas", menuHomeIdea)
			Sessions[GetSessionID(m.Sender.ID)].Status = STATUS_NONE
			return
		}
		var message, localMessage, localContent string

		for i, v := range ideas {
			localMessage = fmt.Sprintf("\n ▪️ %d *%s*", i+1, v.Title)
			if len(v.Content) > 0 {
				localContent = fmt.Sprintf(":\n_%s_", v.Content)
				localMessage += localContent
			}
			message += localMessage
		}
		b.Send(m.Sender, message, ModeMarkdown)
	})
	b.Handle(&btnBack, func(m *tb.Message) {
		sessionID := GetSessionID(m.Sender.ID)
		if Sessions[sessionID].Status == STATUS_ADD_MEDIA {
			b.Send(m.Sender, "Stop media addition.", menuHomeArchive)
		} else {
			b.Send(m.Sender, "Home", menuHomeIdea)
		}
		Sessions[sessionID].Status = STATUS_NONE

	})

	b.Handle(&btnIdea, func(m *tb.Message) {
		b.Send(m.Sender, "Idea !", menuHomeIdea)
	})
	b.Handle(&btnMedia, func(m *tb.Message) {
		b.Send(m.Sender, "Media!", menuHomeArchive)
	})
	b.Handle(&btnHome, func(m *tb.Message) {
		b.Send(m.Sender, "Home !", menuHome)
		sessionID := GetSessionID(m.Sender.ID)
		Sessions[sessionID].Status = STATUS_NONE
	})
	b.Handle(&btnAddMedia, func(m *tb.Message) {
		b.Send(m.Sender, "Send any pictures or videos for save it !", menuBack)
		sessionID := GetSessionID(m.Sender.ID)
		Sessions[sessionID].Status = STATUS_ADD_MEDIA
	})
	b.Handle(&btnShowMedias, func(m *tb.Message) {
		controllers.SendFilesToUser(m.Sender.Username, b, m)
	})
	b.Start()
}
