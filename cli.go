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

func loadData(usr *store.User, dest, key string) {

	log.Println("load from trackthis")
	s := makeStore(dest == "tp", usr)

	usr.Login(s)

	tt := trackthis.NewClient()

	tt.Init(trackthis.Config{AuthToken: key}).
		Load().
		Store(s)
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

func doSync(usr *store.User, from, to string) {

	fs := makeStore(from == "tp", usr)
	ts := makeStore(to == "tp", usr)

	usr.Login(fs)
	usr.Login(ts)

	if usr.CanLogin() {
		log.Println("we should sync from [%s] to [%s]", from, to)
		return
	}
	log.Println("we cannot sync as the user doesn't have valid creds")
}

func doQuery(usr *store.User, storeName, types string) {

	//todo the split of the types we want to query for
	justAlphaNumeric := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	s := makeStore(storeName == "tp", usr)
	qryToRun := &store.Query{Types: strings.FieldsFunc(types, justAlphaNumeric)}

	usr.Login(s)

	data, _ := s.Run(qryToRun)
	log.Printf("%v", data)
}

func checkPledges(usr *store.User) {
	s := store.NewLocalClient(usr)
	usr.Login(s)
	pledges, _ := s.Load()
	log.Printf("found pledges %v", pledges)
}

func makePledge(usr *store.User, p *store.Pledge) {
	p.UserId = usr.Id
	log.Printf("pleadge %v", p)
	s := store.NewLocalClient(usr)
	usr.Login(s)
	s.Register(p)
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

	user := &store.User{}

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
	pledgeData := &store.Pledge{Type: "Equal"}
	pledge := flag.Bool("pledge", false, "make a pledge")
	pledgeData.Deadline, _ = time.Parse("2006-01-02", *flag.String("pledge_date", "", "date the pledge finished e.g. 2015-11-20"))
	pledgeData.TargetValue = *flag.String("pledge_value", "", "target value e.g. <=7.5")
	pledgeData.Feed = *flag.String("pledge_feed", "", "the pledge feed e.g `smbg` or `weight`")
	pledgeData.Pledge = *flag.Float64("pledge_wager", 0, "the wager")
	pledgeData.CounterPledge = *flag.Float64("pledge_counter_wager", 0, "the counter wager")

	//sync flags
	sync := flag.Bool("sync", false, "sync data")
	syncFrom := flag.String("sync_from", "local", "local(local-store),  tp(tp-store)")
	syncTo := flag.String("sync_to", "local", "local(local-store),  tp(tp-store)")

	//tidepool user
	user.Name = *flag.String("tp_usr", "", "local(local-store),  tp(tp-store)")
	user.Pw = *flag.String("tp_pw", "", "local(local-store),  tp(tp-store)")

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
		makePledge(user, pledgeData)
	}

}
