package actions

import (
	"log"
	"strconv"
	"strings"

	"github.com/jh-bate/intertidal/data"
)

type (
	CheckPledge struct {
		store  data.Store
		pledge *data.Pledge
		Yay    bool
		name   string
	}
)

func CheckPledgeAction(calldata interface{}, store data.Store, name string) *CheckPledge {
	return &CheckPledge{store: store, name: name}
}

func (a *CheckPledge) Execute() (interface{}, error) {

	a.store.Find(a.name, &a.pledge)

	//query for this pledge
	qry := &data.Query{
		UserId:   a.pledge.UserId,
		Types:    []string{a.pledge.Feed},
		FromTime: a.pledge.Started.String(),
	}

	results, _ := a.store.Query(data.DATA_COLLECTION, qry)

	total := 0.0

	for i := range results {
		total = total + results[i]["value"].(float64)
	}

	log.Printf("pledge was for %s of %s result is %.1f", a.pledge.Feed, a.pledge.Target, (total / float64(len(results))))

	return hasMeetTarget(a.pledge, total), nil
}

// the target value can be a basic condition that we need to
// evaluate e.g.  >= 10.3
func hasMeetTarget(p *data.Pledge, actualVal float64) bool {

	targetVal := 0.0
	cmpOp := "=="

	parts := strings.Split(p.Target, " ")

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
