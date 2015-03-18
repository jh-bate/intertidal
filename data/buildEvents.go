package data

import (
	"log"
	"strconv"
	"strings"
)

type (
	EventBuilder struct {
		registerPatterns map[Pattern]string
		unprocessed      []UnprocessedEvent
		events           []interface{}
	}
	UnprocessedEvent struct {
		Text, Date, Device string
	}
	//Enum type's
	Pattern string
)

const (
	SRC_INTERTIDAL = "intertidal"
	//Available Pattern's
	PatternSmbg  Pattern = "smbg"
	PatternFood  Pattern = "food"
	PatternBolus Pattern = "bolus"
	PatternBasal Pattern = "basal"
)

func NewEventBuilder() *EventBuilder {
	return &EventBuilder{}
}

//update the registered patterns we use for matching
func (self *EventBuilder) AddPatterns(newPatterns map[Pattern]string) *EventBuilder {
	self.registerPatterns = newPatterns
	return self
}

//update the registered patterns we use for matching
func (self *EventBuilder) AddUnprocessedEvents(unprocessed []UnprocessedEvent) *EventBuilder {
	self.unprocessed = unprocessed
	return self
}

//e.g. text is #bg=9.9 and pattern is #bg= gives value of 9.9
func getVal(pattern, text string) string {
	if strings.Index(strings.ToUpper(text), strings.ToUpper(pattern)) != -1 {
		strVal := strings.Split(text, pattern)[1]
		log.Printf("val %s", strVal)
		return strVal
	}
	return ""
}

//update the registered patterns we use for matching
func (self *EventBuilder) Process() (events []interface{}) {
	for i := range self.unprocessed {
		if self.registerPatterns[PatternFood] != "" {
			if strVal := getVal(self.registerPatterns[PatternFood], self.unprocessed[i].Text); strVal != "" {
				events = append(events, self.BuildFood(strVal, self.unprocessed[i].Date, self.unprocessed[i].Device))
			}
		}
		if self.registerPatterns[PatternSmbg] != "" {
			if strVal := getVal(self.registerPatterns[PatternSmbg], self.unprocessed[i].Text); strVal != "" {
				events = append(events, self.BuildBg(strVal, self.unprocessed[i].Date, self.unprocessed[i].Device))
			}
		}
		if self.registerPatterns[PatternBolus] != "" {
			if strVal := getVal(self.registerPatterns[PatternBolus], self.unprocessed[i].Text); strVal != "" {
				events = append(events, self.BuildBolus(strVal, self.unprocessed[i].Date, self.unprocessed[i].Device))
			}
		}
		if self.registerPatterns[PatternBasal] != "" {
			if strVal := getVal(self.registerPatterns[PatternBasal], self.unprocessed[i].Text); strVal != "" {
				events = append(events, self.BuildBasal(strVal, self.unprocessed[i].Date, self.unprocessed[i].Device))
			}
		}
	}
	return events
}

//create note type from text
func (self *EventBuilder) BuildNote(text, date, device string) *NoteEvent {
	return &NoteEvent{
		Common:    Common{Type: "note", Source: SRC_INTERTIDAL, DeviceId: device, Time: date},
		Text:      text,
		CreatorId: device,
	}
}

//create bg type from text
func (self *EventBuilder) BuildBg(val, date, device string) *BgEvent {
	bgVal, _ := strconv.ParseFloat(val, 64)

	return &BgEvent{
		Common: Common{Type: "smbg", DeviceId: device, Source: SRC_INTERTIDAL, Time: date},
		Value:  bgVal,
	}
}

//create food type from text
func (self *EventBuilder) BuildFood(val, date, device string) *FoodEvent {
	carbVal, _ := strconv.ParseFloat(val, 64)
	return &FoodEvent{
		Common: Common{Type: "food", DeviceId: device, Source: SRC_INTERTIDAL, Time: date},
		Carbs:  carbVal,
	}
}

//create bolus type from text
func (self *EventBuilder) BuildBolus(val, date, device string) *BolusEvent {
	bolusVal, _ := strconv.ParseFloat(val, 64)
	return &BolusEvent{
		Common:  Common{Type: "bolus", DeviceId: device, Source: SRC_INTERTIDAL, Time: date},
		SubType: "injected",
		Value:   bolusVal,
		Insulin: "novolog",
	}
}

//create basal type from text
func (self *EventBuilder) BuildBasal(val, date, device string) *BasalEvent {
	basalVal, _ := strconv.ParseFloat(val, 64)
	return &BasalEvent{
		Common:       Common{Type: "basal", DeviceId: device, Source: SRC_INTERTIDAL, Time: date},
		DeliveryType: "injected",
		Value:        basalVal,
		Insulin:      "lantus",
		Duration:     86400000,
	}
}
