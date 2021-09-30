package bot

import (
	"event-bot-wpp/src/event"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

func (wa *WaHandler) listEvents(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	files, err := ioutil.ReadDir("../events/")
	if err != nil {
		log.Fatalf("error listing events: %v", err)
	}
	strFiles := ""
	for _, f := range files {
		strFiles = strFiles + "\n→ " + f.Name()[:len(f.Name())-5]
	}
	res := "𝐄𝐕𝐄𝐍𝐓𝐎𝐒 𝐂𝐑𝐈𝐀𝐃𝐎𝐒" + strFiles + "\nPara abrir um evento, use o comando\n```!abrir [nome do evento]```"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) inviteUser(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	singleContact := message.ContextInfo.QuotedMessage.ContactMessage

	// when used without quoting a contact
	if message.ContextInfo.QuotedMessage == nil || (message.ContextInfo.QuotedMessage.ContactsArrayMessage == nil && singleContact == nil) {
		invalidUsageStr := "Utilização do comando:\n\n```!convidar``` como *RESPOSTA* a um *CONTATO* ou uma *LISTA DE CONTATOS*"
		wa.sendMessage(message.Info.RemoteJid, invalidUsageStr, contextInfo)
		return
	}

	inviteStr := fmt.Sprintf("     ------ 📩 𝐂𝐎𝐍𝐕𝐈𝐓𝐄 📩 ------\n\n🎈🥂 *Você foi convidado para um evento!* 🥳🎉\n\nResponda ```!sim``` para confirmar sua presença.\n\nResponda ```!lista``` para ver a lista de convidados\n\n🎫 *Atividade*: %v\n🏡 *Local*: %v\n📅 *Horário*: %v\n",
		wa.Event.Activity,
		wa.Event.Venue,
		wa.Event.Date,
	)

	if singleContact != nil {
		vcard := *(singleContact.Vcard)
		displayName := *(singleContact.DisplayName)
		indexWaid := strings.Index(vcard, "waid=")
		if indexWaid != -1 {
			i := 5
			for ; vcard[indexWaid+i] != ':'; i++ {
			}
			waid := event.Waid(vcard[indexWaid+5 : indexWaid+i])
			wa.Event.Invite(waid, displayName)
			invitedNoticeStr := displayName + " foi adicionado à lista de convidados com sucesso!"
			wa.sendMessage(message.Info.RemoteJid, invitedNoticeStr, whatsapp.ContextInfo{})
			if wa.Event.InvitesSent {
				composedWaid := (fmt.Sprintf("%v@s.whatsapp.net", waid))
				wa.sendMessage(composedWaid, inviteStr, whatsapp.ContextInfo{})
			}
		}
		return
	}
	contacts := message.ContextInfo.QuotedMessage.ContactsArrayMessage.Contacts
	for i := range contacts {
		vcard := *(contacts[i].Vcard)
		displayName := *(contacts[i].DisplayName)
		indexWaid := strings.Index(vcard, "waid=")
		if indexWaid != -1 {
			i := 5
			for ; vcard[indexWaid+i] != ':'; i++ {
			}
			waid := event.Waid(vcard[indexWaid+5 : indexWaid+i])
			wa.Event.Invite(waid, displayName)
			invitedNoticeStr := displayName + " foi adicionado à lista de convidados com sucesso!"
			wa.sendMessage(message.Info.RemoteJid, invitedNoticeStr, whatsapp.ContextInfo{})
			if wa.Event.InvitesSent {
				composedWaid := (fmt.Sprintf("%v@s.whatsapp.net", waid))
				wa.sendMessage(composedWaid, inviteStr, whatsapp.ContextInfo{})
			}
		}
	}
}

func (wa *WaHandler) unconfirmAllUsers(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	wa.Event.UndoConfirmations()
	res := "Todas as confirmações foram desfeitas!"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) sendInvites(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	alreadySent := wa.Event.InvitesSent
	if alreadySent {
		res := "Os convites já foram enviados anteriormente!"
		wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
		return
	}
	for key := range wa.Event.Invited {
		str := fmt.Sprintf("     ------ 📩 𝐂𝐎𝐍𝐕𝐈𝐓𝐄 📩 ------\n\n🎈🥂 *Você foi convidado para um evento!* 🥳🎉\n\nResponda ```!sim``` para confirmar sua presença.\n\nResponda ```!lista``` para ver a lista de convidados\n\n🎫 *Atividade*: %v\n🏡 *Local*: %v\n📅 *Horário*: %v\n",
			wa.Event.Activity,
			wa.Event.Venue,
			wa.Event.Date,
		)
		composedStr := fmt.Sprintf("%v@s.whatsapp.net", key)
		wa.sendMessage(composedStr, str, whatsapp.ContextInfo{})
	}
	res := "Convites enviados!"
	wa.sendMessage(message.Info.RemoteJid, res, whatsapp.ContextInfo{})
}

func (wa *WaHandler) newEvent(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, filename string) {
	wa.Event = event.NewEvent(filename)
	res := "*EVENTO CRIADO!*\nPara configurar o evento utilize os seguintes comandos:\n\n```!nome [nome do evento]```\n```!data [data/hora do evento]```\n```!local [local do evento]```\n\nPara adicionar convidados, responda com ```!convidar``` uma mensagem contendo um *contato* ou uma *lista de contatos*\n\nPara checar o status do evento use o comando ```!status```\n\nApós adicionar todos os convidados, use o comando ```!enviar``` para enviar os convites"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)

}

func (wa *WaHandler) loadEvent(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, filename string) {
	files, err := ioutil.ReadDir("../events/")
	if err != nil {
		log.Fatalf("erro ao abrir o diretório: %v", err)
	}
	flag := false
	for _, f := range files {
		if filename+".json" == f.Name() {
			flag = true
		}
	}
	if flag {
		wa.Event, err = event.LoadEvent(filename)
		wa.sendMessage(message.Info.RemoteJid, "Evento carregado!", contextInfo)

	} else {
		res := "'" + filename + "' não é um evento existente.\n\nDiga ```!eventos``` para ver os eventos salvos ou ```!novo [filename]``` para criar um novo evento."
		wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
	}
}

func (wa *WaHandler) setName(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, name string) {
	wa.Event.SetActivity(name)
	res := fmt.Sprintf("Nome do evento alterado para: %s", name)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) setDate(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, date string) {
	wa.Event.SetActivity(date)
	res := fmt.Sprintf("Data do evento alterado para: %s", date)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) setVenue(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, venue string) {
	wa.Event.SetActivity(venue)
	res := fmt.Sprintf("Local do evento alterado para: %s", venue)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) InviteGroup(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	remoteJid := event.Waid(message.Info.RemoteJid)

	wa.Event.InviteGroup(remoteJid)
	res := "Participantes desse grupo podem entrar para a lista de convidados utilizando o comando\n```!entrar [seu nome]```"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) makeAdmin(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	remoteJid := event.Waid(message.Info.RemoteJid)

	wa.Event.AddAdmin(remoteJid)
	res := fmt.Sprintf("Você foi adicionado como administrador do evento *%s*!", wa.Event.Activity)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) isGoing(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, id event.Waid) {
	wa.Event.Going(id)
	res := "*Sua presença foi confirmada!*🥳🤩🎉\n\nResponda ```!lista``` para ver a lista de convidados"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) notGoing(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, id event.Waid) {
	wa.Event.NotGoing(id)
	res := "*Que pena que não poderá ir...*😢😔\n\nCaso mude de ideia, basta enviar ```!sim```"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) getStatus(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {

	res := wa.Event.GetStatus()
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}
