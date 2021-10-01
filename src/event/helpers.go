package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	emojiWoman  = []string{"🙋‍♀️", "🙋🏻‍♀️", "🙋🏼‍♀️", "🙋🏽‍♀️", "🙋🏾‍♀️", "🙋🏿‍♀️"}
	emojiMan    = []string{"🙋‍♂️", "🙋🏻‍♂️", "🙋🏼‍♂️", "🙋🏽‍♂️", "🙋🏾‍♂️", "🙋🏿‍♂️"}
	emojiPerson = []string{"🙋", "🙋🏻", "🙋🏼", "🙋🏽", "🙋🏾", "🙋🏿"}
)

// LoadEvent loads an event from a given json file
func LoadEvent(file string) (Event, error) {
	filePath := fmt.Sprintf("events/%s.json", file)

	e := Event{}
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return e, fmt.Errorf("unable to open '%s': %v", filePath, err)
	}

	json.Unmarshal(buf, &e)
	log.Printf("%s.json foi carregado para memória\n", file)

	return e, nil
}

// NewEvent creates a new event with a given filename
func NewEvent(file string) Event {
	e := Event{
		FileName:  file,
		Invited:   make(map[Waid]nameAndRSVP),
		AllowJoin: make([]Waid, 0),
		Admins:    make([]Waid, 0),
	}
	return e
}

// SaveEvent saves the current event state to it's file
func (e Event) SaveEvent() error {
	filePath := fmt.Sprintf("events/%s.json", e.FileName)

	buf, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return fmt.Errorf("error marshling event: %v", err)
	}
	ioutil.WriteFile(filePath, buf, 0644)
	return nil
}
