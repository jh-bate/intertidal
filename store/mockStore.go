package store

/*
 * Simple implamentation of key/value store using Bolt
 */

type MockClient struct {
	dataToReturn   []interface{}
	configToReturn []interface{}
}

func NewMockClient(data, config []interface{}) *MockClient {

	return &MockClient{
		dataToReturn:   data,
		configToReturn: config,
	}
}

func (b *MockClient) Close() error {
	return nil
}
func (b *MockClient) Ping() error {
	return nil
}
func (b *MockClient) StoreData(key string, data []interface{}) error {
	b.dataToReturn = data
	return nil
}
func (b *MockClient) RetrieveData(key string) (results []interface{}, err error) {
	return b.dataToReturn, nil
}
func (b *MockClient) StoreConfig(key string, data []interface{}) error {
	b.configToReturn = data
	return nil
}
func (b *MockClient) RetrieveConfig(key string) (results []interface{}, err error) {
	return b.configToReturn, nil
}
