package store

/*
 * Generic store client
 */
type Client interface {
	Ping() error
	StoreUser(key, token string) error
	RetrieveUser(key string) ([]interface{}, error)
	StoreUserData(usr string, data []interface{}) error
	RetrieveUserData(usr string) ([]interface{}, error)
	StoreConfig(key string, data []interface{}) error
	RetrieveConfig(key string) ([]interface{}, error)
}

/*
 * Simple mocked key/value store for testing
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
