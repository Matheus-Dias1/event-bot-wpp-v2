package event

import (
	"event-bot-wpp/src/emoji"
	"fmt"
	"log"
	"strings"
)

type nameAndRSVP struct {
	Name          string `json:"name"`
	Confirmed     bool   `json:"confirmed"`
	Going         bool   `json:"going"`
	EmojiGender   string `json:"emojiGender"`
	EmojiSkinTone string `json:"emojiSkinTone"`
}

type Waid string

type Event struct {
	Activity    string               `json:"activity"`
	Venue       string               `json:"venue"`
	Date        string               `json:"date"`
	FileName    string               `json:"fileName"`
	Invited     map[Waid]nameAndRSVP `json:"invited"`
	Admins      []Waid               `json:"admins"`
	AllowJoin   []Waid               `json:"allowJoin"`
	InvitesSent bool                 `json:"invitesSent"`
}

// SetActivity sets the Event's activity title
func (e *Event) SetActivity(activity string) {
	e.Activity = activity
	log.Println("Nome do evento alterado para", e.Activity)
}

// SetVenue sets the Event's activity venue
func (e *Event) SetVenue(venue string) {
	e.Venue = venue
	log.Println("Local do evento alterado para", e.Venue)
}

// SetDate sets the Event's activity date
func (e *Event) SetDate(date string) {
	e.Date = date
	log.Println("Data do evento alterado para", e.Date)
}

// Going confirms someone's presence at the Event by their WAID
func (e *Event) Going(id Waid) {
	name := e.Invited[id].Name
	emoji := e.Invited[id].EmojiGender
	skinTone := e.Invited[id].EmojiSkinTone
	e.Invited[id] = nameAndRSVP{name, true, true, emoji, skinTone}
	log.Printf("%s (%s) confirmou presença no evento %s\n", name, id, e.Activity)
}

// NotGoing confirms someone's abscence at the Event by their WAID
func (e *Event) NotGoing(id Waid) {
	name := e.Invited[id].Name
	emoji := e.Invited[id].EmojiGender
	skinTone := e.Invited[id].EmojiSkinTone
	e.Invited[id] = nameAndRSVP{name, true, false, emoji, skinTone}
	log.Printf("%s (%s) não irá ao evento %s\n", name, id, e.Activity)
}

// Invite adds someone to the Invited slice
func (e *Event) Invite(id Waid, Name string) {
	e.Invited[id] = nameAndRSVP{Name, false, false, "PERSON", "YELLOW"}
	log.Printf("%s (%s) foi convidado para o evento %s\n", Name, id, e.Activity)
}

// InviteGroup adds the group remotejid to the AllowJoin slice
func (e *Event) InviteGroup(remotejid Waid) {
	e.AllowJoin = append(e.AllowJoin, remotejid)
	log.Printf("O grupo (%s) foi convidado para o evento %s\n", remotejid, e.Activity)
}

// AddAdmin adds the person WAID to the Admins slice
func (e *Event) AddAdmin(id Waid) {
	e.Admins = append(e.Admins, id)
	log.Printf("(%s) foi promovido a admin do evento\n", id)
}

// IsAdmin checks if a waid is in and admin
func (e Event) IsAdmin(id Waid) bool {
	for _, jid := range e.Admins {
		if jid == id {
			return true
		}
	}
	return false
}

// UndoConfirmations undo all the confirmations for the event
func (e *Event) UndoConfirmations() {
	for key, val := range e.Invited {
		e.Invited[key] = nameAndRSVP{
			Name:          val.Name,
			Confirmed:     false,
			Going:         false,
			EmojiGender:   val.EmojiGender,
			EmojiSkinTone: val.EmojiSkinTone,
		}
	}
	e.InvitesSent = false
}

// GetStatus returns a string with the event data and invited people info
func (e Event) GetStatus() string {
	template := "```Atividade:``` %v\n```Local:``` %v\n```Data:``` %v\n\n*Convidados:*"
	str := fmt.Sprintf(template, e.Activity, e.Venue, e.Date)

	for _, val := range e.Invited {

		emojiReq := emoji.PresenceEmoji{
			Gender:   val.EmojiGender,
			SkinTone: val.EmojiSkinTone,
		}

		if val.Confirmed && val.Going {
			emojiReq.Going = "IS_GOING"
		} else if val.Confirmed && !val.Going {
			emojiReq.Going = "NOT_GOING"
		} else {
			emojiReq.Going = "UNCONFIRMED"
		}

		emoji := emoji.GetEmoji(emojiReq)
		str += "\n" + val.Name + " (" + emoji + ")"

	}
	return str
}

// Checks if an event was loaded
func (e Event) IsEventLoaded() bool {
	return e.FileName != ""
}

// Checks if someone is invited
func (e Event) IsInvited(remotejid string) (bool, Waid) {
	for key := range e.Invited {
		if strings.Index(remotejid, string(key)) == 0 {
			return true, key
		}

	}
	return false, Waid("")
}

// Checks if user has a emoji configured
func (e Event) EmojiConfigured(id Waid) bool {
	person := e.Invited[id]
	return !(person.EmojiGender == "PERSON" && person.EmojiSkinTone == "YELLOW")
}

// Set persons emoji
func (e *Event) SetEmoji(id Waid, gender, skin string) {
	name := e.Invited[id].Name
	confirmed := e.Invited[id].Confirmed
	going := e.Invited[id].Going

	e.Invited[id] = nameAndRSVP{
		Name:          name,
		Confirmed:     confirmed,
		Going:         going,
		EmojiGender:   gender,
		EmojiSkinTone: skin,
	}
}
