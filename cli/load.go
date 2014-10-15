package main

import (
	"flag"
	"log"

	"github.com/jh-bate/intertidal/api"
	"github.com/jh-bate/intertidal/flood/sms"
	"github.com/jh-bate/intertidal/flood/trackthis"
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

func loadFromSms(token string, stash *store.BoltClient) {
	log.Println("load from sms")
	lt := sms.NewClient()

	lt.Init(sms.Config{AuthToken: token, AccountSid: token}).
		Load().
		StashLocal("test", stash)
}

func loadFromTrackThis(token string, stash *store.BoltClient) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()
	p := platform.NewClient(
		&platform.Config{
			Auth:   "https://staging-api.tidepool.io/auth",
			Upload: "https://staging-uploads.tidepool.io/data",
		},
		"jamie@tidepool.org",
		"blip4life",
	)

	p.StashUserLocal(stash)

	tt.Init(trackthis.Config{AuthToken: token}).
		Load().
		StorePlatform(p).
		StashLocal(p.User.Token, stash)
}

func main() {

	srcPtr := flag.String("s", "", "where we are getting the data from")
	authPtr := flag.String("t", "", "auth token for source")
	//destPtr := flag.String("d", "stash", "where the data will be put")

	flag.Parse()

	stash := store.NewBoltClient()

	/*if *destPtr == "stash" {
		stash := store.NewBoltClient()
	} else if *destPtr == "tp" {

	}*/

	if *srcPtr == string(api.SourceSms) {
		loadFromSms(*authPtr, stash)
	} else if *srcPtr == string(api.SourceTrackThis) {
		loadFromTrackThis(*authPtr, stash)
	}
}
