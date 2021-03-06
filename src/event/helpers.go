package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var (
	emojiWoman  = []string{"πββοΈ", "ππ»ββοΈ", "ππΌββοΈ", "ππ½ββοΈ", "ππΎββοΈ", "ππΏββοΈ"}
	emojiMan    = []string{"πββοΈ", "ππ»ββοΈ", "ππΌββοΈ", "ππ½ββοΈ", "ππΎββοΈ", "ππΏββοΈ"}
	emojiPerson = []string{"π", "ππ»", "ππΌ", "ππ½", "ππΎ", "ππΏ"}
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
	log.Printf("%s.json foi carregado para memΓ³ria\n", file)

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
