package event

import (
	"fmt"
	"log"
)

type nameAndRSVP struct {
	Name      string `json:"name"`
	Confirmed bool   `json:"confirmed"`
	Going     bool   `json:"going"`
}

type waid string

type Event struct {
	Activity    string               `json:"activity"`
	Venue       string               `json:"venue"`
	Date        string               `json:"date"`
	FileName    string               `json:"fileName"`
	Invited     map[waid]nameAndRSVP `json:"invited"`
	Admins      []waid               `json:"admins"`
	AllowJoin   []waid               `json:"allowJoin"`
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
func (e *Event) Going(id waid) {
	name := e.Invited[id].Name
	e.Invited[id] = nameAndRSVP{name, true, true}
	log.Printf("%s (%s) confirmou presença no evento %s\n", name, id, e.Activity)
}

// NotGoing confirms someone's abscence at the Event by their WAID
func (e *Event) NotGoing(id waid) {
	name := e.Invited[id].Name
	e.Invited[id] = nameAndRSVP{name, true, false}
	log.Printf("%s (%s) confirmou presença no evento %s\n", name, id, e.Activity)
}

// Invite adds someone to the Invited slice
func (e *Event) Invite(id waid, Name string) {
	e.Invited[id] = nameAndRSVP{Name, false, false}
	log.Printf("%s (%s) foi convidado para o evento %s\n", Name, id, e.Activity)
}

// InviteGroup adds the group remotejid to the AllowJoin slice
func (e *Event) InviteGroup(remotejid waid) {
	e.AllowJoin = append(e.AllowJoin, remotejid)
	log.Printf("O grupo (%s) foi convidado para o evento %s\n", remotejid, e.Activity)
}

// AddAdmin adds the person WAID to the Admins slice
func (e *Event) AddAdmin(id waid) {
	e.Admins = append(e.Admins, id)
	log.Printf("(%s) foi promovido a admin do evento\n", id)
}

// IsAdmin checks if a waid is in and admin
func (e Event) IsAdmin(id waid) bool {
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
			Name:      val.Name,
			Confirmed: false,
			Going:     false,
		}
	}
	e.InvitesSent = false
}

// GetStatus returns a string with the event data and invited people info
func (e Event) GetStatus() string {
	template := "```Atividade:``` %v\n```Local:``` %v\n```Data:``` %v\n\n*Convidados:*"
	str := fmt.Sprintf(template, e.Activity, e.Venue, e.Date)

	var val nameAndRSVP
	for _, val = range e.Invited {
		str += "\n" + val.Name
		if val.Confirmed {
			str += " (👍)"
		} else {
			str += " (❓)"
		}
	}
	return str
}
