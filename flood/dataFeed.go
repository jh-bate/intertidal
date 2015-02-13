package flood

import (
	"github.com/jh-bate/intertidal"
)

type DataFeed interface {
	//Note: all functions return an instance of the feed to allow chaining

	//initialise your data feed with config such as oauth key and source url
	Init(config interface{}) *DataFeed

	//load from your initialise data feed
	Load() *DataFeed

	//save the loaded and parsed data into the provided store
	Store(store *intertidal.Store) *DataFeed
}
