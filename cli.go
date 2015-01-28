package main

import (
	"flag"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/jh-bate/intertidal/examples"
	"github.com/jh-bate/intertidal/store"
)

const (
	LOCAL = "local"
)

func loadData(token string, store store.Client) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()

	tt.Init(trackthis.Config{AuthToken: token}).
		Load().
		Store(store)
}

func makeStore(server bool, usr *store.User) store.Client {

	if server && usr.CanLogin() {
		return store.NewTidepoolClient(
			&store.TidepoolConfig{
				Auth:   "https://api.tidepool.io/auth",
				Upload: "https://uploads.tidepool.io/data",
				Query:  "https://api.tidepool.io/query",
			},
			usr.Name,
			usr.Pw)
	}

	return store.NewLocalClient(usr)
}

func runSync() {
	log.Println("we should sync stores")
}

func doQuery(s store.Client, qry *store.Query) {
	data, _ := s.Run(qry)
	log.Printf("%v", data)
}

func checkPledges(s *store.LocalClient) {
	pledges, _ := s.Load()
	log.Printf("found pledges %v", pledges)
}

func main() {

	/*
		### Load data
		go run cli.go -load_data -load_key=fcc72d1e11858cdce96783033c1f5af46185eec78d9233a8161b05845671e943
		### Query Data
		go run cli.go -query_data -query_types bolus,food


		### Create Pledge
		go run cli.go -make_pledge -pledge_date "2015-06-06" -pledge_value "<= 7.5" -pledge_type smbg -pledge_wager 3 -pledge_counter_wager 1
		### Query Pledge
		go run cli.go -query_pledge
	*/

	//load flags
	load := flag.Bool("load", false, "do a data load")
	loadKey := flag.String("load_key", "", "key for the df(data-feed)")
	loadInto := flag.String("load_into", "local", "local(local-store),  tp(tp-store)")

	//query flags
	query := flag.Bool("query", false, "run a query")
	querySource := flag.String("query_source", "local", "local(local-store),  tp(tp-store)")
	queryTypes := flag.String("query_types", "smbg", "query types e.g. smbg, food")
	queryPledge := flag.Bool("query_pledge", false, "query all registered pledges")

	//pledge flages
	pledge := flag.Bool("pledge", false, "make a pledge")
	pledgeDate := flag.String("pledge_date", "", "date the pledge finished e.g. 2015-11-20")
	pledgeValue := flag.String("pledge_value", "", "target value e.g. <=7.5")
	pledgeType := flag.String("pledge_type", "", "the pledge type e.g `smbg` or `weight`")
	pledgeWager := flag.Float64("pledge_wager", 0, "the wager")
	pledgeCounter := flag.Float64("pledge_counter_wager", 0, "the counter wager")

	//sync flags
	sync := flag.Bool("sync", false, "sync data")
	//syncTo := flag.String("sync_to", "local", "local(local-store),  tp(tp-store)")

	//tidepool user
	tpUsr := flag.String("tp_usr", "", "local(local-store),  tp(tp-store)")
	tpPw := flag.String("tp_pw", "", "local(local-store),  tp(tp-store)")

	flag.Parse()

	loggedInUser := &store.User{}

	if *tpUsr != "" && *tpPw != "" {
		loggedInUser.Name = *tpUsr
		loggedInUser.Pw = *tpPw
	} else {
		loggedInUser.Id = "todo2"
	}

	//loading
	if *load && *loadKey != "" {
		sendToStore := makeStore(*loadInto == "tp", loggedInUser)
		loadData(*loadKey, sendToStore)
	}

	//querying
	if *query {
		//todo the split
		justAlphaNumeric := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
		qryFromStore := makeStore(*querySource == "tp", loggedInUser)
		doQuery(qryFromStore, &store.Query{Types: strings.FieldsFunc(*queryTypes, justAlphaNumeric)})
	}
	//syncing
	if *sync && loggedInUser.CanLogin() {
		runSync()
	}

	//pledge
	if *queryPledge == true {
		checkPledges(store.NewLocalClient(loggedInUser))
	}

	if *pledge == true {

		dl, _ := time.Parse("2006-01-02", *pledgeDate)

		pledge := &store.Pledge{
			UserId:        loggedInUser.Id,
			Feed:          *pledgeType,
			Type:          "Equal",
			Deadline:      dl,
			TargetValue:   *pledgeValue,
			Pledge:        *pledgeWager,
			CounterPledge: *pledgeCounter,
		}

		log.Printf("pleadge %v", pledge)
		store := store.NewLocalClient(loggedInUser)
		store.Register(pledge)
	}

}
