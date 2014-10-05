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
