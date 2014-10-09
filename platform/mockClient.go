package platform

import (
	"log"
)

type (
	MockClient struct {
		token string
	}
)

func NewMockClient(usr, pw string) (*MockClient, error) {

	client := &MockClient{}

	if tkn, err := client.login(usr, pw); err != nil {
		log.Println("Error init client: ", err)
		return nil, err
	} else {
		client.token = tkn
		return client, nil
	}
}

func (mock *MockClient) login(usr, pw string) (token string, err error) {
	return "fairy.dust.as.a.token", nil
}

func (mock *MockClient) Signup(name, pw, contact string) (User, error) {
	return &User{Name: name, Contact: contact, Id: "123.456.777", Token: mock.token}, nil
}

func (mock *MockClient) LoadInto(data *[]interface{}) error {
	return nil
}

func (mock *MockClient) LoadFrom(userid string) ([]interface{}, error) {
	return nil, nil
}
