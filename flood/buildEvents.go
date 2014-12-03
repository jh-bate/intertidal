package flood

import (
	"log"
	"strconv"
	"strings"

	"github.com/jh-bate/intertidal/backend/platform"
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
	intertidal = "intertidal"
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
func (self *EventBuilder) BuildNote(text, date, device string) *platform.NoteEvent {
	return &platform.NoteEvent{
		Common:    platform.Common{Type: "note", Source: intertidal, DeviceId: device, Time: date},
		Text:      text,
		CreatorId: device,
	}
}

//create bg type from text
func (self *EventBuilder) BuildBg(val, date, device string) *platform.BgEvent {
	bgVal, _ := strconv.ParseFloat(val, 64)

	return &platform.BgEvent{
		Common: platform.Common{Type: "smbg", DeviceId: device, Source: intertidal, Time: date},
		Value:  bgVal,
	}
}

//create food type from text
func (self *EventBuilder) BuildFood(val, date, device string) *platform.FoodEvent {
	carbVal, _ := strconv.ParseFloat(val, 64)
	return &platform.FoodEvent{
		Common: platform.Common{Type: "food", DeviceId: device, Source: intertidal, Time: date},
		Carbs:  carbVal,
	}
}

//create bolus type from text
func (self *EventBuilder) BuildBolus(val, date, device string) *platform.BolusEvent {
	bolusVal, _ := strconv.ParseFloat(val, 64)
	return &platform.BolusEvent{
		Common:  platform.Common{Type: "bolus", DeviceId: device, Source: intertidal, Time: date},
		SubType: "injected",
		Value:   bolusVal,
		Insulin: "novolog",
	}
}

//create basal type from text
func (self *EventBuilder) BuildBasal(val, date, device string) *platform.BasalEvent {
	basalVal, _ := strconv.ParseFloat(val, 64)
	return &platform.BasalEvent{
		Common:       platform.Common{Type: "basal", DeviceId: device, Source: intertidal, Time: date},
		DeliveryType: "injected",
		Value:        basalVal,
		Insulin:      "lantus",
		Duration:     86400000,
	}
}
