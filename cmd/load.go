package main

import (
	"flag"
	"log"

	"github.com/jh-bate/intertidal/api"
	"github.com/jh-bate/intertidal/flood/trackthis"
	"github.com/jh-bate/intertidal/flood/twilio"
	"github.com/jh-bate/intertidal/store"
)

func loadFromTwilio(token string, stash *store.BoltClient) {
	log.Println("load from twilio")
	tw := twilio.NewClient()

	tw.Init(twilio.Config{AuthToken: token, AccountSid: token}, stash)

	tw.Load()
	tw.StashLocal()
}

func loadFromTrackThis(token string, stash *store.BoltClient) {
	log.Println("load from trackthis")

	tt := trackthis.NewClient()

	tt.Init(trackthis.Config{AuthToken: token}, stash)

	tt.Load()
	tt.StashLocal()
}

func main() {

	srcPtr := flag.String("s", "", "where we are getting the data from")
	authPtr := flag.String("t", "", "auth token for source")
	destPtr := flag.String("d", "stash", "where the data will be put")

	flag.Parse()

	stash := store.NewBoltClient()

	log.Println("flags:", *srcPtr)
	log.Println("flags:", *authPtr)
	log.Println("flags:", *destPtr)

	if *srcPtr == string(api.SourceTwilio) {
		loadFromTwilio(*authPtr, stash)
	} else if *srcPtr == string(api.SourceTrackThis) {
		loadFromTrackThis(*authPtr, stash)
	}
}
