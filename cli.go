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

	query := flag.Bool("query", false, "run a query")

	flag.Parse()

	//loading data
	if *load {
		log.Println("loading ...")
		dl := actions.LoadDataAction([]byte(fmt.Sprintf(`{"authToken":"%s"}`, *loadKey)))
		dl.Execute()
		log.Printf("loaded [%d] records", len(dl.Results))
		if *save == true {
			log.Print("saving results ...")
			dData := actions.AddDataAction(dl.Results, str, data.DATA_COLLECTION)
			dData.Execute()
			log.Print("saved results")
		} else {
			log.Print("asked not to save the results")
		}
	}

	//add a pledge
	if *pledge && *query == false {
		log.Println("adding pledge ...")
		log.Printf("pledge %v", pledgeData)
		p := actions.AddPledgeAction(str, data.PLEDGE_COLLECTION, pledgeData)
		p.Execute()
		log.Print("added pledge")
	}

	//query a pledge
	if *pledge && *query {
		log.Println("querying pledge ...")
		p := actions.CheckPledgeAction(str, data.PLEDGE_COLLECTION)
		p.Execute()
		log.Print("queryed pledge. ")
		log.Printf("target meet? %t", p.Yay)
	}
}
