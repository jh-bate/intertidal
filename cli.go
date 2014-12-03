package main

import (
	"flag"
	"log"

	"github.com/jh-bate/intertidal/api"
	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
	"github.com/jh-bate/intertidal/examples"
)

func loadFromTrackThis(token string, tp *platform.TidepoolClient, stash *store.BoltClient) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()

	tp.StashUserLocal(stash)

	tt.Init(trackthis.Config{AuthToken: token}).
		Load().
		StorePlatform(tp).
		StashLocal(tp.User.Token, stash)
}

func main() {

	srcPtr := flag.String("s", "", "where we are getting the data from")
	usrPtr := flag.String("u", "", "tidepool username")
	pwPtr := flag.String("p", "", "tidepool password")
	srcTokenPtr := flag.String("t", "", "auth token for source")

	flag.Parse()

	stash := store.NewBoltClient()

	tp := platform.NewClient(
		&platform.Config{
			Auth:   "https://api.tidepool.io/auth",
			Upload: "https://uploads.tidepool.io/data",
		},
		*usrPtr,
		*pwPtr,
	)

	if *srcPtr == string(api.SourceTrackThis) {
		loadFromTrackThis(*srcTokenPtr, tp, stash)
	}
}
