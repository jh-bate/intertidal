package flood

import (
	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
)

type Feed interface {
	Init(config interface{}) *Feed
	Load() *Feed
	StashLocal(key string, local store.Client) *Feed
	StorePlatform(platform platform.Client) *Feed
}

/*
 * Mock feed for testing
 */
type MockFeed struct{}

func NewMockFeed() *MockFeed {
	return &MockFeed{}
}

func (f *MockFeed) Init(config interface{}, store *store.Client) {}
func (f *MockFeed) Load()                                        {}
func (f *MockFeed) StashLocal()                                  {}
func (f *MockFeed) StorePlatform(platform *platform.Client)      {}
