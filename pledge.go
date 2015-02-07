package intertidal

import (
	"log"
	"time"
)

const PLEDGES_COLLECTION = "pledges"

type (
	Pledge struct {
		UserId        string    `json:"source"`
		Feed          string    `json:"feed-address"`
		Deadline      time.Time `json:"deadline"`
		Started       time.Time `json:"started"`
		Type          string    `json:"bet-type"` //
		TargetValue   string    `json:"target-value"`
		Pledge        float64   `json:"wager"`
		CounterPledge float64   `json:"counterwager"`
	}
)

func Make(userId, feed, target string, wager, counterWager float64, started, deadline time.Time) *Pledge {
	return &Pledge{UserId: userId, Feed: feed, Deadline: deadline, Started: started, Type: "Equal", TargetValue: target, Pledge: wager, CounterPledge: counterWager}
}

func (p *Pledge) Evaluate(store Store) bool {

	pledgeQry := &Query{
		UserId:   p.UserId,
		Types:    []string{p.Feed},
		FromTime: p.Started.String(),
	}

	results, _ := store.Query(pledgeQry)

	total := 0.0

	for i := range results {
		total = total + results[i]["value"].(float64)
	}

	log.Printf("count %d total %.1f ave %.1f target %s", len(results), total, total/float64(len(results)), p.TargetValue)

	return false
}

func (p *Pledge) Save(store Store) error {
	return store.Save(PLEDGES_COLLECTION, p)
}
