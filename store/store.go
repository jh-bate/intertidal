package store

/*
 * Generic store client
 */
type Client interface {
	Ping() error
	StoreData(key string, data []interface{}) error
	RetrieveData(key string) ([]interface{}, error)
	StoreConfig(key string, data []interface{}) error
	RetrieveConfig(key string) ([]interface{}, error)
}
