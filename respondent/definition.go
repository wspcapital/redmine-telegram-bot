package respondent

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"redmine-telegram-bot/converse_state"
	state_store "redmine-telegram-bot/converse_state/dynamodb"
)

const QuestionUnset  = 0

func New(sender MessageSender, receivedMessage *tgbotapi.Message) (*Respondent, error) {
	r := new(Respondent)

	r.Sender = sender
	r.ReceivedMessage = receivedMessage

	state, err := converse_state.FindOrCreateState(state_store.New(), r.ReceivedMessage.Chat.ID)
	if err != nil {
		return nil, err
	}

	r.State = state

	if os.Getenv("BOT_DEBUG") == "true" {
		log.Println("Current state:")
		log.Println(state)
	}

	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	//bundle.MustLoadMessageFile("assets/lang/response.en.yaml")
	//bundle.MustLoadMessageFile("assets/lang/response.ru.yaml")
	//bundle.MustLoadMessageFile("assets/lang/response.ua.yaml")

	r.Localizer = i18n.NewLocalizer(bundle, "en")

	return r, nil
}

// MessageSender interface define methods what use respondent.
type MessageSender interface {
	Send(chattable tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Controller interface {
	CreateMessage(respondent *Respondent) tgbotapi.Chattable
}

// Main structure of respondent.
type Respondent struct {
	Sender          MessageSender
	ReceivedMessage *tgbotapi.Message
	State           *converse_state.State
	Localizer       *i18n.Localizer
}

// Reply to a message.
// Here is all the business logic of response.
func (r *Respondent) Reply() {

	var ctrl Controller

	// Routing for messages

	if r.ReceivedMessage.IsCommand() {
		switch r.ReceivedMessage.Command() {
		case "help":
			ctrl = new(HelpCommand)
			break
		default:
			ctrl = new(UnknownCommand)
			break
		}
	} else if r.State.IsJustCreated() {
		ctrl = new(NewState)
	} else {

	}

	// Send message
	if ctrl != nil {
		r.sendMessage(ctrl.CreateMessage(r))
	}
}

// Send message.
func (r *Respondent) sendMessage(msg tgbotapi.Chattable) {
	_, err := r.Sender.Send(msg)
	if err != nil {
		log.Println("Can't send response:")
		log.Println(err)
	}
}
