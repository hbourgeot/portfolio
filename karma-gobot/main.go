package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"karmagot/internal/karma"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var botUsername string
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

	bot, err := tgbotapi.NewBotAPI("5917051686:AAEf5hPR7stgvKb2Ig38IXfydEk88vpnUnI")
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	db, err := openDB("root:Wini.h16b.@/karmabot?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
		return
	}
	infoLog.Println("Starting bot...")

	karmas := karma.KarmaModel{DB: db}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		for _, user := range update.Message.NewChatMembers {
			if user.UserName == "karmagobot" {
				botUsername = user.UserName
				break
			}
		}

		if update.Message.Chat.IsPrivate() || update.Message.Chat.IsChannel() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This bot can't run on private conversations and channels. Use it in a group")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				errorLog.Fatal(err)
				return
			}
			continue
		}

		chat := strings.ToLower(strings.ReplaceAll(update.Message.Chat.Title, " ", "_"))
		if botUsername == "karmagobot" {
			err = karmas.CreateTable(chat)
			if err != nil {
				errorLog.Println(err)
			}
			continue
		}

		lastUpdated, noRows := karmas.GetLastUpdated(update.Message.From.UserName, chat)

		if noRows {
			err := karmas.InsertUsers(update.Message.From.UserName, chat)
			if err != nil {
				errorLog.Fatal(err)
			}
			continue
		}

		// For +1 or -1
		if strings.Contains(update.Message.Text, "+1") || strings.Contains(update.Message.Text, "-1") {
			if update.Message.From.UserName == update.Message.ReplyToMessage.From.UserName {
				msgError := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR, You cannot add or subtract karma yourself.")
				if _, err := bot.Send(msgError); err != nil {
					errorLog.Fatal(err)
					return
				}
				continue
			}

			if checkGiveKarma(lastUpdated) {
				msgError := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR, You must to have to wait 1 minute to give karma.")
				if _, err := bot.Send(msgError); err != nil {
					errorLog.Fatal(err)
					return
				}
				continue
			} else if strings.Contains(update.Message.Text, "+1") {
				err = karmas.AddKarma(update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName, chat)
				if err != nil {
					errorLog.Println(err)
					err := karmas.InsertUsers(update.Message.ReplyToMessage.From.UserName, chat)
					if err != nil {
						errorLog.Println(err)
					}
					continue
				}

			} else if strings.Contains(update.Message.Text, "-1") {
				fmt.Println("me diste -1", update.Message.Text)
				err = karmas.SubstractKarma(update.Message.From.UserName, update.Message.ReplyToMessage.From.UserName, chat)
				if err != nil {
					errorLog.Println(err, ", creating new user...")
					err := karmas.InsertUsers(update.Message.ReplyToMessage.From.UserName, chat)
					if err != nil {
						errorLog.Println(err)
					}
					continue
				}
			}
		} else {
			continue
		}

		userKarma, _, err := karmas.GetActualKarma(update.Message.ReplyToMessage.From.UserName, chat)
		if err != nil {
			errorLog.Println(err)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.ReplyToMessage.From.UserName+" has now "+strconv.Itoa(userKarma)+" of karma")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(msg); err != nil {
			errorLog.Fatal(err)
			return
		}
		// For bot commands
		if update.Message.IsCommand() {
			cmdText := update.Message.Command()
			switch cmdText {
			case "/karma":
				userKarma, _, err := karmas.GetActualKarma(update.Message.From.UserName, chat)
				if err != nil {
					errorLog.Fatal(err)
					return
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "@"+update.Message.From.UserName+" has "+strconv.Itoa(userKarma)+" of karma.")
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := bot.Send(msg); err != nil {
					errorLog.Fatal(err)
					return
				}
				break
			case "/karmalove":
				users, err := karmas.GetKarmas(chat, true)
				if err != nil {
					errorLog.Println(err)
					continue
				}

				usersString := "Most loved users of " + chat + "\n"
				for i, user := range users {
					usersString += fmt.Sprintf("%d. %s, %d karma.\n", i, update.Message.From.UserName, user.Count)
				}
				break
			case "/karmahate":
				users, err := karmas.GetKarmas(chat, false)
				if err != nil {
					errorLog.Println(err)
					continue
				}

				usersString := "Most hated users of " + chat + "\n"
				for i, user := range users {
					usersString += fmt.Sprintf("%d. %s, %d karma.\n", i, update.Message.From.UserName, user.Count)
				}
			}
		}
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func checkGiveKarma(lastKarmaGived time.Time) bool {
	log.Println(lastKarmaGived, time.Now().UTC())
	return lastKarmaGived.Day() == time.Now().UTC().Day() &&
		lastKarmaGived.Hour() == time.Now().UTC().Hour() &&
		lastKarmaGived.Minute() == time.Now().UTC().Minute()
}
