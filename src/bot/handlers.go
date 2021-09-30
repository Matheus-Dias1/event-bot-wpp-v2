package bot

import (
	"log"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
)

func (wa *WaHandler) HandleTextMessage(message whatsapp.TextMessage) {

	// fmt.Printf("Timestamp: %v\nMessage ID: %v\nRemoteJid: %v\nContextInfo: %v\n\tText: %v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo, message.Text)

	if message.Info.Timestamp < wa.StartTime {
		return
	}

	quotedMessage := proto.Message{
		Conversation: &message.Text,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: message.Info.Id,
		Participant:     "", //Who sent the original message
	}

	cmd := parseText(message.Text)
	if !cmd.isValid {
		if cmd.resStr != "" {
			wa.sendMessage(message.Info.RemoteJid, cmd.resStr, ContextInfo)
		}
		return
	}

	// checks if event is loaded when using some commands
	if invalid, shouldAlert := wa.invalidUsage(cmd.name); invalid {
		if shouldAlert {
			res := "Nenhum evento carregado.\n\nDiga ```!abrir [filename]``` para carregar um evento existente.\nDiga ```!novo [filename]``` para criar um novo evento\n"
			wa.sendMessage(message.Info.RemoteJid, res, ContextInfo)
		}
		return
	}

	// message sent from admins
	if wa.isAdmin(message.Info) {
		switch cmd.name {
		case "LIST_EVENTS":
			wa.listEvents(message, ContextInfo)
		case "INVITE_USERS":
			wa.inviteUser(message, ContextInfo)
		case "UNCONFIRM_ALL_USERS":
			wa.unconfirmAllUsers(message, ContextInfo)
		case "SEND_INVITES":
			wa.sendInvites(message, ContextInfo)
		case "NEW_EVENT":
			wa.newEvent(message, ContextInfo, cmd.value)
		case "LOAD_EVENT":
			wa.loadEvent(message, ContextInfo, cmd.value)
		case "SET_EVENT_ACTIVITY":
			wa.setName(message, ContextInfo, cmd.value)
		case "SET_EVENT_VENUE":
			wa.setVenue(message, ContextInfo, cmd.value)
		case "SET_EVENT_DATE":
			wa.setDate(message, ContextInfo, cmd.value)
		}

	}

	// message sent from me
	if message.Info.FromMe {
		switch cmd.name {
		case "INVITE_GROUP":
			wa.InviteGroup(message, ContextInfo)
		case "MAKE_ADMIN":
			wa.makeAdmin(message, ContextInfo)
		}
	}

	// message sent from someone who was invited for the event
	if invited, waid := wa.Event.IsInvited(message.Info.RemoteJid); invited {
		switch cmd.name {
		case "IS_GOING":
			wa.isGoing(message, ContextInfo, waid)
		case "NOT_GOING":
			wa.notGoing(message, ContextInfo, waid)
		case "GET_STATUS":
			wa.getStatus(message, ContextInfo)
		}
	}

}

func (wa *WaHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		for {
			log.Printf("Connection failed, underlying error: %v", e.Err)
			log.Println("Waiting 30sec...")
			<-time.After(30 * time.Second)
			log.Println("Reconnecting...")
			err := wa.C.Restore()
			if err == nil {
				wa.StartTime = uint64(time.Now().Unix())
				break
			}
		}

	} else {
		log.Printf("error occoured: %v\n", err)
	}
}
