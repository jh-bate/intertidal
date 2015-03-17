package intertidal

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const PLEDGES_COLLECTION = "pledges"

type (
	Pledge struct {
		//user that the is making the pledge
		UserId string `json:"source"`
		//the query that will give us the data to evaluate the pledge against
		Feed string `json:"feed-address"`
		//when the pledge will be evaluted and finalised
		Deadline time.Time `json:"deadline"`
		//when the peldge started
		Started time.Time `json:"started"`
		//normally `Equal`
		Type string `json:"bet-type"`
		//the target value that the pledge will evaluate against e.g. <= 8.4
		TargetValue string `json:"target-value"`
		//what we are pledging for the outcome to be `true`
		Pledge float64 `json:"wager"`
		//what is pledged if the outcome is `false`
		CounterPledge float64 `json:"counterwager"`
	}
)

func Make(userId, feed, target string, wager, counterWager float64, started, deadline time.Time) *Pledge {

	if started.IsZero() {
		log.Print("started date not set so we have set as now")
		started = time.Now()
	}
	return &Pledge{UserId: userId, Feed: feed, Deadline: deadline, Started: started, Type: "Equal", TargetValue: target, Pledge: wager, CounterPledge: counterWager}
}

// the target value can be a basic condition that we need to
// evaluate e.g.  >= 10.3
func (p *Pledge) meetTarget(actualVal float64) bool {

	targetVal := 0.0
	cmpOp := "=="

	parts := strings.Split(p.TargetValue, " ")

	if len(parts) == 2 {
		cmpOp = parts[0]
		targetVal, _ = strconv.ParseFloat(parts[1], 64)
	} else if len(parts) == 1 {
		targetVal, _ = strconv.ParseFloat(parts[0], 64)
	} else {
		return false
	}

	switch cmpOp {
	case ">":
		return actualVal > targetVal
	case ">=":
		return actualVal >= targetVal
	case "<":
		return actualVal < targetVal
	case "<=":
		return actualVal <= targetVal
	default:
		return actualVal == targetVal
	}
}

//evalute the pledge
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

	return p.meetTarget(total / float64(len(results)))
}

//save the pledge to the target store
func (p *Pledge) Save(store Store) error {
	return store.Save(PLEDGES_COLLECTION, p)
}
