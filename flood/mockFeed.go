package flood

import (
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

type MockFeed struct {
}

func NewMockFeed() *MockFeed {
	return &MockFeed{}
}

func (f *MockFeed) Init(config interface{}, store *store.Client) {}
func (f *MockFeed) Load()                                        {}
func (f *MockFeed) StashLocal()                                  {}
func (f *MockFeed) StorePlatform(platform *platform.Client)      {}
