package data

import (
	"log"
	"time"
)

const (
	PLEDGE_COLLECTION = "pledge"
)

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
		//the target value that the pledge will evaluate against e.g. <= 8.4
		Target string `json:"target-value"`
		//what we are pledging for the outcome to be `true`
		IfMeet float64 `json:"wager"`
		//what is pledged if the outcome is `false`
		IfNotMeet float64 `json:"counterwager"`
	}
)

func New(userId, feed, target string, wager, counterWager float64, started, deadline time.Time) *Pledge {

	if started.IsZero() {
		log.Print("started date not set so we have set as now")
		started = time.Now()
	}
	return &Pledge{UserId: userId, Feed: feed, Deadline: deadline, Started: started, Target: target, IfMeet: wager, IfNotMeet: counterWager}
}
