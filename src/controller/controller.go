package controller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/tucnak/telebot.v2"
)

var statusCheckFunc func() (string, string, []Button, []File)

type Button struct {
	Text          string
	RMMsgsOnClick bool
	OnClick       func()
}

type File struct {
	Name    string
	Content string
}

type Chat struct {
	Id string
}

func (c *Chat) Recipient() string {
	return c.Id
}

var chat = &Chat{os.Getenv("WRAPPER_TELEGRAM_CHAT")}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getChatInt64() int64 {
	i, err := strconv.ParseInt(chat.Id, 10, 64)
	check(err)
	return i
}

var chatIdInt64 = getChatInt64()

func txtToHtml(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, "&", "&amp;"), "<", "&lt;"), ">", "&gt;")
}

func fileToDoc(file File) *telebot.Document {
	return &telebot.Document{
		File:     telebot.FromReader(strings.NewReader(file.Content)),
		MIME:     "text/plain",
		FileName: file.Name,
	}
}

func filesToTelebotAlbum(files []File, caption ...string) telebot.Album {
	album := telebot.Album{}
	captionPresent := len(caption) > 0
	for i := 0; i < len(files); i++ {
		doc := fileToDoc(files[i])
		if captionPresent && i == len(files) - 1 {
			doc.Caption = caption[0]
		}
		album = append(album, doc)
	}
	return album
}

func initBot() *telebot.Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("WRAPPER_TELEGRAM_TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	check(err)
	bot.Handle("/status", func(m *telebot.Message) {
		if m.Chat.ID != chatIdInt64 {
			return
		}
		messagesToDelete := []*telebot.Message{}
		state, startTime, buttons, files := statusCheckFunc()
		messageText := fmt.Sprintf("<b>Статус сервиса:</b> %s\n", txtToHtml(state))
		if startTime != "" {
			messageText += fmt.Sprintf("<b>Работает с %s</b>\n", startTime)
		}
		keyboard := buttonsToTelebotKeyboard(buttons, func() []*telebot.Message {
			return messagesToDelete
		}, bot)
		if len(files) == 0 {
			messageText += "<b>Сервис не выводил ничего в stdout/stderr</b>"
		}
		msg, _ := bot.Send(chat, messageText, &telebot.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: keyboard,
		})
		messagesToDelete = append(messagesToDelete, msg)
		if len(files) != 0 {
			msgs, _ := bot.SendAlbum(chat, filesToTelebotAlbum(files, messageText), &telebot.SendOptions{
				ReplyTo: msg,
			})
			for i := 0; i < len(msgs); i++ {
				messagesToDelete = append(messagesToDelete, &msgs[i])
			}
		}
	})
	go bot.Start()
	return bot
}

var bot = initBot()

func buttonsToTelebotKeyboard(buttons []Button, getMessagesToDelete func() []*telebot.Message, bot *telebot.Bot) *telebot.ReplyMarkup {
	menu := &telebot.ReplyMarkup{}
	rows := []telebot.Row{}
	for i := 0; i < len(buttons); i++ {
		button := buttons[i]
		inlineBtn := menu.Data(button.Text, uuid.NewString())
		bot.Handle(&inlineBtn, func(c *telebot.Callback) {
			bot.Respond(c)
			if button.RMMsgsOnClick {
				messages := getMessagesToDelete()
				for i := 0; i < len(messages); i++ {
					bot.Delete(messages[i])
				}
			}
			button.OnClick()
		})
		rows = append(rows, menu.Row(inlineBtn))
	}
	menu.Inline(rows...)
	return menu
}

func Send(text string, buttons []Button, files ...File) {
	messagesToDelete := []*telebot.Message{}
	keyboard := buttonsToTelebotKeyboard(buttons, func() []*telebot.Message {
		return messagesToDelete
	}, bot)
	msg, _ := bot.Send(chat, text, &telebot.SendOptions{
		ParseMode:   "HTML",
		ReplyMarkup: keyboard,
	})
	messagesToDelete = append(messagesToDelete, msg)
	msgs, _ := bot.SendAlbum(chat, filesToTelebotAlbum(files), &telebot.SendOptions{
		ReplyTo: msg,
	})
	for i := 0; i < len(msgs); i++ {
		messagesToDelete = append(messagesToDelete, &msgs[i])
	}
}

func OnStatusCheck(check func() (string, string, []Button, []File)) {
	statusCheckFunc = check
}
