package respondent

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Unknown command
type UnknownCommand struct {

}

func (c *UnknownCommand) CreateMessage(respondent *Respondent) tgbotapi.Chattable  {

	msgText := respondent.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "UnknownCommandMessage",
			Description: "Welcome message when chat firstly created",
			Other:       "I do not known this command. Please use /help for see *full list of commands*.",
		},
	})

	msg := tgbotapi.NewMessage(respondent.State.ID, msgText)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}

// Help command
type HelpCommand struct {

}

func (c *HelpCommand) CreateMessage(respondent *Respondent) tgbotapi.Chattable  {

	msgText := "help <b>command</b>"

	msg := tgbotapi.NewMessage(respondent.State.ID, msgText)
	msg.ParseMode = tgbotapi.ModeHTML

	return msg
}