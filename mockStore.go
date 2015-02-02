package intertidal

/*
 * Simple mocked key/value store for testing
 */
type MockStoreClient struct {
	dataToReturn []interface{}
}

func NewMockClient(data []interface{}) *MockStoreClient {

	return &MockStoreClient{
		dataToReturn: data,
	}
}

func (b *MockStoreClient) Ping() error {
	return nil
}

func (b *MockStoreClient) Save(userid string, data []interface{}) error {
	b.dataToReturn = data
	return nil
}

func (b *MockStoreClient) Query(qry *Query) (data []interface{}, err error) {
	return b.dataToReturn, nil
}
