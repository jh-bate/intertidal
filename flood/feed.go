package flood

import (
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

type Feed interface {
	Init(config interface{}) *Feed
	Load() *Feed
	StashLocal(key string, local store.Client) *Feed
	StorePlatform(platform platform.Client) *Feed
}
