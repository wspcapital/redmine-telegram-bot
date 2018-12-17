package respondent

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type NewState struct {

}

func (c *NewState) CreateMessage(respondent *Respondent) tgbotapi.Chattable  {
	msg := tgbotapi.NewMessage(respondent.State.ID, "Hi :)")

	return msg
}