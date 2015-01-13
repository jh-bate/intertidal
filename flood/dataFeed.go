package flood

import (
	"github.com/jh-bate/intertidal/store"
)

type DataFeed interface {
	Init(config interface{}) *DataFeed
	Load() *DataFeed
	Store(store *store.Client) *DataFeed
}

/*
 * Mock feed for testing
 */
type MockFeed struct{}

func NewMockFeed() *MockFeed {
	return &MockFeed{}
}

func (f *MockFeed) Init(config interface{})   {}
func (f *MockFeed) Load()                     {}
func (f *MockFeed) Store(store *store.Client) {}
