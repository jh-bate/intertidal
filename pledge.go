package intertidal

import (
	"time"
)

const PLEDGES_COLLECTION = "pledges"

type (
	Pledge struct {
		UserId        string    `json:"source"`
		Feed          string    `json:"feed-address"`
		Deadline      time.Time `json:"deadline"`
		Type          string    `json:"bet-type"` //
		TargetValue   string    `json:"target-value"`
		Pledge        float64   `json:"wager"`
		CounterPledge float64   `json:"counterwager"`
	}
)

func Make(userId, feed, target string, wager, counterWager float64, deadline time.Time) *Pledge {
	return &Pledge{UserId: userId, Feed: feed, Deadline: deadline, Type: "Equal", TargetValue: target, Pledge: wager, CounterPledge: counterWager}
}

func (p *Pledge) Evaluate(store Store) bool {
	return false
}

func (p *Pledge) Save(store Store) error {
	return store.Save(PLEDGES_COLLECTION, p)
}
