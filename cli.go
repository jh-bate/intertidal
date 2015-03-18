package main

import (
	"flag"
	"fmt"
	"log"
	//"strings"
	"time"
	//"unicode"

	"github.com/jh-bate/intertidal/actions"
	"github.com/jh-bate/intertidal/data"
)

const (
	LOCAL = "local"
)

func initStore(usr *data.User) data.Store {
	log.Print("creating local store")
	return data.NewLocalStore(usr)
}

func main() {

	user := &data.User{}
	str := initStore(user)
	mgr := actions.CreateActionManager()

	//load flags
	load := flag.Bool("load", false, "do a data load")
	loadKey := flag.String("load_key", "", "key for the df(data-feed)")
	save := flag.Bool("save", true, "save the results?")

	//pledge flags
	pledge := flag.Bool("pledge", false, "make a pledge")

	pledgeData := &data.Pledge{}

	days := 0
	flag.IntVar(&days, "pledge_length", 90, "number of days before pledge finished e.g. 90")
	pledgeData.Deadline = time.Now().AddDate(0, 0, days)

	daysFromNow := 0
	flag.IntVar(&daysFromNow, "pledge_start", 0, "number of days from now the pledge will start e.g. -30")
	pledgeData.Started = time.Now().AddDate(0, 0, daysFromNow)

	flag.StringVar(&pledgeData.Target, "pledge_value", "", "target value e.g. <=7.5")
	flag.StringVar(&pledgeData.Feed, "pledge_feed", "", "the pledge feed e.g `smbg` or `weight`")
	flag.Float64Var(&pledgeData.IfMeet, "pledge_wager", 0, "the wager")
	flag.Float64Var(&pledgeData.IfNotMeet, "pledge_counter_wager", 0, "the counter wager")

	//query
	query := flag.Bool("query", false, "run a query")

	flag.Parse()

	//loading data
	if *load {
		res, _ := mgr.Execute(actions.ActionTypeLoadData, []byte(fmt.Sprintf(`{"authToken":"%s"}`, *loadKey)), str, data.DATA_COLLECTION)
		if *save == true {
			mgr.Execute(actions.ActionTypeAddData, res, str, data.DATA_COLLECTION)
		} else {
			log.Print("asked not to save the results")
		}
	}

	//add a pledge
	if *pledge && *query == false {
		log.Printf("pledge %v", pledgeData)
		mgr.Execute(actions.ActionTypeAddPledge, pledgeData, str, data.PLEDGE_COLLECTION)
	}

	//query a pledge
	if *pledge && *query {
		res, _ := mgr.Execute(actions.ActionTypeCheckPledge, nil, str, data.PLEDGE_COLLECTION)
		log.Printf("target meet? %t", res)
	}

	//query
	if *query && *pledge == false {
		qry := &data.Query{Types: []string{"smbg"}, UserId: user.Id}
		res, _ := mgr.Execute(actions.ActionTypeQueryData, qry, str, data.DATA_COLLECTION)
		log.Printf("query results? %s", res)
	}
}
