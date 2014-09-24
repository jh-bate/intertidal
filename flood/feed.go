package flood

import (
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

type Source interface {
	Init(config interface{}, store *store.Client)
	Load()
	StashLocal()
	StorePlatform(platform *platform.Client)
}
