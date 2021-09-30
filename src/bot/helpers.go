package bot

import (
	"event-bot-wpp/src/event"
	"fmt"
	"regexp"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

type command struct {
	isValid bool
	resStr  string
	name    string
	value   string
}

var usageMap = map[string]string{
	"!novo":   "nome do arquivo",
	"!abrir":  "nome do arquivo",
	"!nome":   "título do evento",
	"!data":   "data do evento",
	"!local":  "local do evento",
	"!entrar": "seu nome",
}

/*
	NO event loaded required // ADMIN only
		→ !novo [nome do arquivo]
		→ !abrir [nome do arquivo]
		→ !eventos

	event loaded REQUIRED // ADMIN only
		→ !nome [nome do evento]
		→ !data [data/hora do evento]
		→ !local [local do evento]
		→ !convidar
		→ !desconvidar
		→ !enviar

	event loaded REQUIRED // ME only
		→ !grupo
		→ !admin

	event loaded REQUIRED // anyone INVITED
		→ !status // !lista
		→ !sim
		→ !nao // !não

	event loaded REQUIRED // anyone from INVITED group
		→ !status // !lista
		→ !entrar [nome]

*/
func parseText(text string) (cmd command) {
	space := regexp.MustCompile(`\s+`)
	cleanText := strings.TrimSpace(space.ReplaceAllString(text, " "))
	lowerCaseText := strings.ToLower(cleanText)

	cmd = command{isValid: true}

	switch lowerCaseText {
	case "!eventos":
		cmd.name = "LIST_EVENTS"
		return
	case "!convidar":
		cmd.name = "INVITE_USERS"
		return
	case "!desconvidar":
		cmd.name = "UNCONFIRM_ALL_USERS"
		return
	case "!enviar":
		cmd.name = "SEND_INVITES"
		return
	case "!grupo":
		cmd.name = "INVITE_GROUP"
		return
	case "!admin":
		cmd.name = "MAKE_ADMIN"
		return
	case "!lista", "!status":
		cmd.name = "GET_STATUS"
		return
	case "!não", "!nao":
		cmd.name = "NOT_GOING"
		return
	case "!sim":
		cmd.name = "IS_GOING"
		return

	// FOR INCORECT COMMANDS USAGE
	case "!novo", "!abrir", "!nome", "!data", "!local", "!entrar":
		cmd.isValid = false
		cmd.resStr = fmt.Sprint("Utilização do comando:\n\n```%s [%s]```", lowerCaseText, usageMap[lowerCaseText])
		return

	}
	cmdList := strings.SplitN(lowerCaseText, " ", 2)
	cmdListCapitalized := strings.SplitN(cleanText, " ", 2)
	if len(cmdList) != 2 {
		cmd.isValid = false
		return
	}

	switch 0 {
	case strings.Index(lowerCaseText, "!novo "):
		cmd.name = "NEW_EVENT"
		cmd.value = cmdList[1]
		return
	case strings.Index(lowerCaseText, "!abrir "):
		cmd.name = "LOAD_EVENT"
		cmd.value = cmdList[1]
		return
	case strings.Index(lowerCaseText, "!nome "):
		cmd.name = "SET_EVENT_ACTIVITY"
		cmd.value = cmdListCapitalized[1]
		return
	case strings.Index(lowerCaseText, "!local "):
		cmd.name = "SET_EVENT_VENUE"
		cmd.value = cmdListCapitalized[1]
		return
	case strings.Index(lowerCaseText, "!data "):
		cmd.name = "SET_EVENT_DATE"
		cmd.value = cmdListCapitalized[1]
		return
	case strings.Index(lowerCaseText, "!entrar "):
		cmd.name = "JOIN_FROM_GROUP"
		cmd.value = cmdListCapitalized[1]
		return
	}

	cmd.isValid = false
	return

}

func (wa *WaHandler) invalidUsage(command string) (bool, bool) {
	isLoaded := wa.Event.IsEventLoaded()

	switch command {
	case "NEW_EVENT", "LIST_EVENTS", "LOAD_EVENT":
		return false, false
	case
		"SET_EVENT_ACTIVITY",
		"SET_EVENT_DATE",
		"SET_EVENT_VENUE",
		"INVITE_USERS",
		"UNCONFIRM_ALL_USERS",
		"SEND_INVITES",
		"INVITE_GROUP",
		"MAKE_ADMIN":
		return !isLoaded, !isLoaded

	default:
		return !wa.Event.IsEventLoaded(), false
	}
}

func (wa *WaHandler) isAdmin(info whatsapp.MessageInfo) bool {
	id := event.Waid(info.RemoteJid)
	if info.FromMe || wa.Event.IsAdmin(id) {
		return true
	}
	return false
}
