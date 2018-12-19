package respondent

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type NewState struct {

}

func (c *NewState) CreateMessage(respondent *Respondent) tgbotapi.Chattable  {

	msgText := respondent.Localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "WelcomeMessage",
			Description: "Welcome message when chat firstly created",
			Other:       "Hi! Nice to meet you.\nFor more information, use the command. /help",
		},
	})

	msg := tgbotapi.NewMessage(respondent.State.ID, msgText)

	respondent.State.CurrentQuestion = QuestionUnset
	respondent.State.UpdateState()

	return msg
}