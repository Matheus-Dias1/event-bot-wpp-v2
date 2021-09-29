package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// LoadEvent loads an event from a given json file
func LoadEvent(file string) (Event, error) {
	filePath := fmt.Sprintf("../events/%s.json", file)

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
		FileName: file,
	}
	return e
}

// SaveEvent saves the current event state to it's file
func (e Event) SaveEvent() error {
	filePath := fmt.Sprintf("../events/%s.json", e.FileName)

	buf, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("error marshling event: %v", err)
	}
	ioutil.WriteFile(filePath, buf, 0644)
	log.Printf("O evento %v foi salvo em %v.json\n", e.Activity, e.FileName)
	return nil
}
