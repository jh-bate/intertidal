package main

import (
	"flag"
	"log"
	"strings"
	"time"
	"unicode"

	"github.com/jh-bate/intertidal"
	"github.com/jh-bate/intertidal/examples/datafeed"
)

const (
	LOCAL = "local"
)

func loadData(usr *intertidal.User, dest, key string) {

	log.Println("load from trackthis")
	s := makeStore(dest == "tp", usr)

	s.Login()

	tt := trackthis.NewClient()

	tt.Init(trackthis.Config{AuthToken: key}).
		Load().
		Store(s)
}

func makeStore(server bool, usr *intertidal.User) intertidal.Store {

	/*if server && usr.CanLogin() {
		return intertidal.NewTidepoolStore(
			&intertidal.TidepoolConfig{
				Auth:   "https://api.tidepool.io/auth",
				Upload: "https://uploads.tidepool.io/data",
				Query:  "https://api.tidepool.io/query",
			},
			usr.Name,
			usr.Pw)
	}*/

	return intertidal.NewLocalStore(usr)
}

func doSync(usr *intertidal.User, from, to string) {

	fs := makeStore(from == "tp", usr)
	ts := makeStore(to == "tp", usr)

	fs.Login()
	ts.Login()

	if usr.CanLogin() {
		log.Println("we should sync from [%s] to [%s]", from, to)
		return
	}
	log.Println("we cannot sync as the user doesn't have valid creds")
}

func doQuery(usr *intertidal.User, storeName, types string) {

	//todo the split of the types we want to query for
	justAlphaNumeric := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	s := makeStore(storeName == "tp", usr)
	qryToRun := &intertidal.Query{Types: strings.FieldsFunc(types, justAlphaNumeric)}

	s.Login()
	results, _ := s.Query(qryToRun)
	log.Printf("query returned %v", results)
}

func checkPledges(usr *intertidal.User) {
	s := intertidal.NewLocalStore(usr)
	s.Login()
	var pledges intertidal.Pledge
	s.Find(intertidal.PLEDGES_COLLECTION, &pledges)

	good := pledges.Evaluate(s)
	log.Printf("won? %t", good)
}

func makePledge(usr *intertidal.User, p *intertidal.Pledge) {
	p.UserId = usr.Id
	log.Printf("pleadge %v", p)
	s := intertidal.NewLocalStore(usr)
	s.Login()
	s.Save(intertidal.PLEDGES_COLLECTION, p)
}

func main() {

	/*
		### Load data
		go run examples/cli.go -load_data -load_key=fcc72d1e11858cdce96783033c1f5af46185eec78d9233a8161b05845671e943
		### Query Data
		go run examples/cli.go -query_data -query_types bolus,food


		### Create Pledge
		go run examples/cli.go -make_pledge -pledge_date "2015-06-06" -pledge_value "<= 7.5" -pledge_type smbg -pledge_wager 3 -pledge_counter_wager 1
		### Query Pledge
		go run examples/cli.go -query_pledge
	*/

	user := &intertidal.User{}

	//load flags
	load := flag.Bool("load", false, "do a data load")
	loadKey := flag.String("load_key", "", "key for the df(data-feed)")
	loadInto := flag.String("load_into", "local", "local(local-store),  tp(tp-store)")

	//query flags
	query := flag.Bool("query", false, "run a query")
	queryFrom := flag.String("query_from", "local", "local(local-store),  tp(tp-store)")
	queryTypes := flag.String("query_types", "smbg", "query types e.g. smbg, food")
	queryPledge := flag.Bool("query_pledge", false, "query all registered pledges")

	//pledge flages
	pledgeData := &intertidal.Pledge{Type: "Equal"}
	pledge := flag.Bool("pledge", false, "make a pledge")

	days := 0
	flag.IntVar(&days, "pledge_length", 90, "number of days before pledge finished e.g. 90")
	pledgeData.Deadline = time.Now().AddDate(0, 0, days)

	daysFromNow := 0
	flag.IntVar(&daysFromNow, "pledge_start", 0, "number of days from now the pledge will start e.g. -30")
	pledgeData.Started = time.Now().AddDate(0, 0, daysFromNow)

	flag.StringVar(&pledgeData.TargetValue, "pledge_value", "", "target value e.g. <=7.5")
	flag.StringVar(&pledgeData.Feed, "pledge_feed", "", "the pledge feed e.g `smbg` or `weight`")
	flag.Float64Var(&pledgeData.Pledge, "pledge_wager", 0, "the wager")
	flag.Float64Var(&pledgeData.CounterPledge, "pledge_counter_wager", 0, "the counter wager")

	//sync flags
	sync := flag.Bool("sync", false, "sync data")
	syncFrom := flag.String("sync_from", "local", "local(local-store),  tp(tp-store)")
	syncTo := flag.String("sync_to", "local", "local(local-store),  tp(tp-store)")

	//tidepool user
	flag.StringVar(&user.Name, "tp_usr", "", "local(local-store),  tp(tp-store)")
	flag.StringVar(&user.Pw, "tp_pw", "", "local(local-store),  tp(tp-store)")

	flag.Parse()

	//loading data
	if *load {
		loadData(user, *loadInto, *loadKey)
	}

	//querying data
	if *query {
		doQuery(user, *queryFrom, *queryTypes)
	}

	//syncing data
	if *sync {
		doSync(user, *syncFrom, *syncTo)
	}

	//pledge
	if *queryPledge {
		checkPledges(user)
	}

	if *pledge {
		log.Printf("pledge %v", pledgeData)
		makePledge(user, pledgeData)
	}

}
