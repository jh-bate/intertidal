package flood

import (
	"strconv"

	"github.com/jh-bate/intertidal/backend/platform"
)

const (
	INTERTIDAL = "intertidal"
)

func MakeBg(bgString, date, device string) *platform.BgEvent {
	bgVal, _ := strconv.ParseFloat(bgString, 64)
	return &platform.BgEvent{Common: platform.Common{Type: "smbg", DeviceId: device, Source: INTERTIDAL, Time: date}, Value: bgVal}
}

func MakeNote(noteString, date, device string) *platform.NoteEvent {
	return &platform.NoteEvent{Common: platform.Common{Type: "note", Source: INTERTIDAL, DeviceId: device, Time: date}, Text: noteString, CreatorId: device}
}

func MakeCarb(carbString, date, device string) *platform.FoodEvent {
	carbVal, _ := strconv.ParseFloat(carbString, 64)
	return &platform.FoodEvent{Common: platform.Common{Type: "food", DeviceId: device, Source: INTERTIDAL, Time: date}, Carbs: carbVal}
}

func MakeBolus(bolusString, date, device string) *platform.BolusEvent {
	bolusVal, _ := strconv.ParseFloat(bolusString, 64)
	return &platform.BolusEvent{Common: platform.Common{Type: "bolus", DeviceId: device, Source: INTERTIDAL, Time: date}, SubType: "injected", Value: bolusVal, Insulin: "novolog"}
}

func MakeBasal(basalString, date, device string) *platform.BasalEvent {
	basalVal, _ := strconv.ParseFloat(basalString, 64)
	return &platform.BasalEvent{Common: platform.Common{Type: "basal", DeviceId: device, Source: INTERTIDAL, Time: date}, DeliveryType: "injected", Value: basalVal, Insulin: "lantus", Duration: 86400000}
}
