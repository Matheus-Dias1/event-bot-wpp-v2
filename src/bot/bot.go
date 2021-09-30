package bot

import (
	"event-bot-wpp/src/event"
	"fmt"
	"os"

	"github.com/Rhymen/go-whatsapp"
)

type WaHandler struct {
	C           *whatsapp.Conn
	Event       event.Event
	EventLoaded bool
	StartTime   uint64
}

func (wa *WaHandler) sendMessage(RemoteJid string, text string, ContextInfo whatsapp.ContextInfo) {

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: RemoteJid,
		},
		ContextInfo: ContextInfo,
		Text:        text,
	}

	_, err := wa.c.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		os.Exit(1)
	}
}
