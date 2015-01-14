package main

import (
	"flag"
	"log"

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

func main() {

	//-do=l -key=xxxx [-to=cs -u= -p=]

	//-do=q [-from=cs -u= -p=]

	//-do=s [-from=cs -u= -p=]

	//what
	load := flag.Bool("l", false, "do a l(oad)")
	query := flag.Bool("q", false, "do a q(uery)")
	sync := flag.Bool("s", false, "do a s(ync)")
	//where from
	from := flag.String("from", "ls", "cs(central-store), ls(local-store), df(data-feed)")
	key := flag.String("key", "", "key for the df(data-feed)")
	//where to
	to := flag.String("to", "ls", "cs(central-store), ls(local-store)")
	//creds
	un := flag.String("u", "", "cs(central-store) username")
	pw := flag.String("p", "", "cs(central-store) password")

	flag.Parse()

	loggedInUser := &store.User{}

	if *un != "" && *pw != "" {
		loggedInUser.Name = *un
		loggedInUser.Pw = *pw
	} else {
		loggedInUser.Id = "todo2"
	}

	toStore := makeStore(*to == "cs", loggedInUser)
	fromStore := makeStore(*from == "cs", loggedInUser)

	if *load && *key != "" {
		loadData(*key, toStore)
	}
	if *query {
		doQuery(fromStore, &store.Query{})
	}
	if *sync && loggedInUser.CanLogin() {
		runSync()
	}

}
