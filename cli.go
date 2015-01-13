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

func load(token string, store store.Client) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()

	tt.Init(trackthis.Config{AuthToken: token}).
		Load().
		Store(store)
}

func makeStore(server bool) store.Client {
	/*store.NewTidepoolClient(
	&store.Config{
		Auth:   "https://api.tidepool.io/auth",
		Upload: "https://uploads.tidepool.io/data",
		Query:  "https://api.tidepool.io/query",
	},
	*usr,
	*pw)*/
	return store.NewLocalClient(&store.User{Id: "todo2"})
}

func sync() {
	log.Println("sync stores")
}

func query() {
	log.Println("query")
}

func main() {

	//-from trackthis -frm_key XXXXX -server

	//-sync -u an@email.org -p 123x43

	//

	//incomming data
	//src := flag.String("frm", trackthis.TRACK_THIS, "where we are getting the data from")
	srcToken := flag.String("frm_key", "", "auth token for source")

	//platfrom
	//usr := flag.String("u", "", "tidepool username")
	//pw := flag.String("p", "", "tidepool password")

	//storage
	//sync := flag.Bool("sync", false, "to you want to sync local with server")
	server := flag.Bool("server", false, "send to server, default is local")

	flag.Parse()

	//load data from trackthis
	load(*srcToken, makeStore(*server))

	/*tp := platform.NewClient(
		&platform.Config{
			Auth:   "https://api.tidepool.io/auth",
			Upload: "https://uploads.tidepool.io/data",
			Query:  "https://api.tidepool.io/query",
		},
		*usr,
		*pw,
	)*/

}
