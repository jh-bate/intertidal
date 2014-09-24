package flood

import (
	"strconv"

	"github.com/jh-bate/intertidal/platform/models"
)

const (
	INTERTIDAL = "intertidal"
)

func MakeBg(bgString, date, device string) *models.BgEvent {
	bgVal, _ := strconv.ParseFloat(bgString, 64)
	return &models.BgEvent{Common: models.Common{Type: "smbg", DeviceId: device, Source: INTERTIDAL, Time: date}, Value: bgVal}
}

func MakeNote(noteString, date, device string) *models.NoteEvent {
	return &models.NoteEvent{Common: models.Common{Type: "note", Source: INTERTIDAL, DeviceId: device, Time: date}, Text: noteString, CreatorId: device}
}

func MakeCarb(carbString, date, device string) *models.FoodEvent {
	carbVal, _ := strconv.ParseFloat(carbString, 64)
	return &models.FoodEvent{Common: models.Common{Type: "food", DeviceId: device, Source: INTERTIDAL, Time: date}, Carbs: carbVal}
}

func MakeBolus(bolusString, date, device string) *models.BolusEvent {
	bolusVal, _ := strconv.ParseFloat(bolusString, 64)
	return &models.BolusEvent{Common: models.Common{Type: "bolus", DeviceId: device, Source: INTERTIDAL, Time: date}, SubType: "injected", Value: bolusVal, Insulin: "novolog"}
}

func MakeBasal(basalString, date, device string) *models.BasalEvent {
	basalVal, _ := strconv.ParseFloat(basalString, 64)
	return &models.BasalEvent{Common: models.Common{Type: "basal", DeviceId: device, Source: INTERTIDAL, Time: date}, DeliveryType: "injected", Value: basalVal, Insulin: "lantus", Duration: 86400000}
}
