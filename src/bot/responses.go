package bot

import (
	"event-bot-wpp/src/emoji"
	"event-bot-wpp/src/event"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

func (wa *WaHandler) listEvents(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	files, err := ioutil.ReadDir("events/")
	if err != nil {
		log.Fatalf("error listing events: %v", err)
	}
	strFiles := ""
	for _, f := range files {
		strFiles = strFiles + "\nβ " + f.Name()[:len(f.Name())-5]
	}
	res := "πππππππ πππππππ" + strFiles + "\n\nPara abrir um evento, use o comando\n```!abrir [nome do evento]```"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) inviteUser(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	singleContact := message.ContextInfo.QuotedMessage.ContactMessage

	// when used without quoting a contact
	if message.ContextInfo.QuotedMessage == nil || (message.ContextInfo.QuotedMessage.ContactsArrayMessage == nil && singleContact == nil) {
		invalidUsageStr := "UtilizaΓ§Γ£o do comando:\n\n```!convidar``` como *RESPOSTA* a um *CONTATO* ou uma *LISTA DE CONTATOS*"
		wa.sendMessage(message.Info.RemoteJid, invalidUsageStr, contextInfo)
		return
	}

	inviteStr := fmt.Sprintf("     ------ π© πππππππ π© ------\n\nππ₯ *VocΓͺ foi convidado para um evento!* π₯³π\n\nResponda ```!sim``` para confirmar sua presenΓ§a.\n\nResponda  ```!nao``` caso nΓ£o possa ir\n\nResponda ```!lista``` para ver a lista de convidados\n\nπ« *Atividade*: %v\nπ‘ *Local*: %v\nπ *HorΓ‘rio*: %v\n",
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
			invitedNoticeStr := displayName + " foi adicionado Γ  lista de convidados com sucesso!"
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
			invitedNoticeStr := displayName + " foi adicionado Γ  lista de convidados com sucesso!"
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
	res := "Todas as confirmaΓ§Γ΅es foram desfeitas!"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) sendInvites(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {
	alreadySent := wa.Event.InvitesSent
	if alreadySent {
		res := "Os convites jΓ‘ foram enviados anteriormente!"
		wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
		return
	}
	for key := range wa.Event.Invited {
		str := fmt.Sprintf("     ------ π© πππππππ π© ------\n\nππ₯ *VocΓͺ foi convidado para um evento!* π₯³π\n\nResponda ```!sim``` para confirmar sua presenΓ§a.\n\nResponda  ```!nao``` caso nΓ£o possa ir\n\nResponda ```!lista``` para ver a lista de convidados\n\nπ« *Atividade*: %v\nπ‘ *Local*: %v\nπ *HorΓ‘rio*: %v\n",
			wa.Event.Activity,
			wa.Event.Venue,
			wa.Event.Date,
		)
		composedStr := fmt.Sprintf("%v@s.whatsapp.net", key)
		wa.sendMessage(composedStr, str, whatsapp.ContextInfo{})
	}
	res := "Convites enviados!"
	wa.sendMessage(message.Info.RemoteJid, res, whatsapp.ContextInfo{})
	wa.Event.InvitesSent = true
}

func (wa *WaHandler) newEvent(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, filename string) {
	wa.Event = event.NewEvent(filename)
	res := "*EVENTO CRIADO!*\nPara configurar o evento utilize os seguintes comandos:\n\n```!nome [nome do evento]```\n```!data [data/hora do evento]```\n```!local [local do evento]```\n\nPara adicionar convidados, responda com ```!convidar``` uma mensagem contendo um *contato* ou uma *lista de contatos*\n\nPara checar o status do evento use o comando ```!status```\n\nApΓ³s adicionar todos os convidados, use o comando ```!enviar``` para enviar os convites"
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)

}

func (wa *WaHandler) loadEvent(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, filename string) {
	files, err := ioutil.ReadDir("events/")
	if err != nil {
		log.Fatalf("erro ao abrir o diretΓ³rio: %v", err)
	}
	flag := false
	for _, f := range files {
		if filename+".json" == f.Name() {
			flag = true
		}
	}
	if flag {
		wa.Event, err = event.LoadEvent(filename)
		if err != nil {
			log.Printf("error loading event: %v", err)
			return
		}
		wa.sendMessage(message.Info.RemoteJid, "Evento carregado!", contextInfo)

	} else {
		res := "'" + filename + "' nΓ£o Γ© um evento existente.\n\nDiga ```!eventos``` para ver os eventos salvos ou ```!novo [filename]``` para criar um novo evento."
		wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
	}
}

func (wa *WaHandler) setName(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, name string) {
	wa.Event.SetActivity(name)
	res := fmt.Sprintf("Nome do evento alterado para: %s", name)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) setDate(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, date string) {
	wa.Event.SetDate(date)
	res := fmt.Sprintf("Data do evento alterado para: %s", date)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) setVenue(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, venue string) {
	wa.Event.SetVenue(venue)
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
	res := fmt.Sprintf("VocΓͺ foi adicionado como administrador do evento *%s*!", wa.Event.Activity)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) isGoing(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, id event.Waid) {
	wa.Event.Going(id)
	res := "*Sua presenΓ§a foi confirmada!*π₯³π€©π\n\nResponda ```!lista``` para ver a lista de convidados"
	if !wa.Event.EmojiConfigured(id) {
		res = res + "\n\nVocΓͺ pode configurar seu emoji que aparecerΓ‘ na lista de convidados com ```!emoji```"
	}
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) notGoing(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, id event.Waid) {
	wa.Event.NotGoing(id)
	res := "*Que pena que nΓ£o poderΓ‘ ir...*π’π\n\nCaso mude de ideia, basta enviar\n```!sim```"
	if !wa.Event.EmojiConfigured(id) {
		res = res + "\n\nVocΓͺ pode configurar seu emoji que aparecerΓ‘ na lista de convidados com ```!emoji```"
	}
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) getStatus(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {

	res := wa.Event.GetStatus()
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) configEmoji(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, id event.Waid, params string) {
	//params parsing
	paramList := strings.Split(params, " ")
	if len(paramList) != 2 {
		wa.emojiHelp(message, contextInfo)
		return
	}

	gender, err := strconv.Atoi(paramList[0])
	if err != nil || gender < 1 || gender > 3 {
		wa.emojiHelp(message, contextInfo)
		return
	}

	skin, err := strconv.Atoi(paramList[1])
	if err != nil || skin < 1 || skin > 5 {
		wa.emojiHelp(message, contextInfo)
		return
	}

	genderStr := GenderMap[gender]
	skinStr := skinToneMap[skin]
	wa.Event.SetEmoji(id, genderStr, skinStr)

	req := emoji.PresenceEmoji{
		Going:    "IS_GOING",
		Gender:   genderStr,
		SkinTone: skinStr,
	}

	updatedEmoji := emoji.GetEmoji(req)
	res := "Seu emoji foi atualizado com sucesso! " + updatedEmoji

	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) emojiHelp(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo) {

	res := `Para configurar seu emoji envie o comando
%s!emoji [gΓͺnero] [pele]%s usando os seguintes atributos:

π¨ - 1
π© - 2
π§ - 3

π» - 1
πΌ - 2
π½ - 3
πΎ - 4
πΏ - 5

ex: %s!emoji 1 3%s = π¨π½`

	codeQuotes := "```"
	res = fmt.Sprintf(res, codeQuotes, codeQuotes, codeQuotes, codeQuotes)
	wa.sendMessage(message.Info.RemoteJid, res, contextInfo)
}

func (wa *WaHandler) JoinFromGroup(message whatsapp.TextMessage, contextInfo whatsapp.ContextInfo, name string) {
	inviteStr := fmt.Sprintf("     ------ π© πππππππ π© ------\n\nππ₯ *VocΓͺ foi convidado para um evento!* π₯³π\n\nResponda ```!sim``` para confirmar sua presenΓ§a.\n\nResponda  ```!nao``` caso nΓ£o possa ir\n\nResponda ```!lista``` para ver a lista de convidados\n\nπ« *Atividade*: %v\nπ‘ *Local*: %v\nπ *HorΓ‘rio*: %v\n",
		wa.Event.Activity,
		wa.Event.Venue,
		wa.Event.Date,
	)
	senderId := message.Info.SenderJid
	atIndex := strings.Index(senderId, "@")
	waid := event.Waid(senderId[0:atIndex])
	wa.Event.Invite(waid, name)

	res := name + " foi adicionado aos convidados!"
	wa.sendMessage(message.Info.RemoteJid, res, whatsapp.ContextInfo{})

	if wa.Event.InvitesSent {
		sendTo := string(waid) + "@s.whatsapp.net"
		wa.sendMessage(sendTo, inviteStr, whatsapp.ContextInfo{})
	}
}
