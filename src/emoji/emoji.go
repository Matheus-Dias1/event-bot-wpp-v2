package emoji

import (
	"fmt"
	"log"
)

type Gender string
type SkinTone string

type PresenceEmoji struct {
	Going    string
	Gender   string
	SkinTone string
}

const (
	isGoing     string = "\U0001F64B%s\U0000200D%s\U0000FE0F"
	notGoing    string = "\U0001F645%s\U0000200D%s\U0000FE0F"
	unconfirmed string = "\U0001F937%s\U0000200D%s\U0000FE0F"

	person Gender = ""
	man    Gender = "\U00002642"
	woman  Gender = "\U00002640"

	yellow              SkinTone = ""
	lightSkinTone       SkinTone = "\U0001F3FB"
	mediumLightSkinTone SkinTone = "\U0001F3FC"
	mediumSkinTone      SkinTone = "\U0001F3FD"
	mediumDarkSkinTone  SkinTone = "\U0001F3FE"
	darkSkinTone        SkinTone = "\U0001F3FF"
)

var GoingMap = map[string]string{
	"IS_GOING":    isGoing,
	"NOT_GOING":   notGoing,
	"UNCONFIRMED": unconfirmed,
}

var GenderMap = map[string]Gender{
	"PERSON": person,
	"MAN":    man,
	"WOMAN":  woman,
}

var SkinToneMap = map[string]SkinTone{
	"YELLOW":                 yellow,
	"LIGHT_SKIN_TONE":        lightSkinTone,
	"MEDIUM_LIGHT_SKIN_TONE": mediumLightSkinTone,
	"MEDIUM_SKIN_TONE":       mediumSkinTone,
	"MEDIUM_DARK_SKIN_TONE":  mediumDarkSkinTone,
	"DARK_SKIN_TONE":         darkSkinTone,
}

// GetEmoji returns an emoji based on the params
func GetEmoji(req PresenceEmoji) string {

	format, ok := GoingMap[req.Going]
	if !ok {
		log.Printf("Unexpected emoji going type: %v", req.Going)
		return "⚠️"
	}

	skinTone, ok := SkinToneMap[req.SkinTone]
	if !ok {
		log.Printf("Unexpected emoji skintone type: %v", req.SkinTone)
		skinTone = yellow
	}

	gender, ok := GenderMap[req.Gender]
	if !ok {
		log.Printf("Unexpected emoji gender type: %v", req.Gender)
		gender = person
	}

	return fmt.Sprintf(format, skinTone, gender)
}
