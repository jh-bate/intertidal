package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const (
	TP_SESSION_TOKEN = "x-tidepool-session-token"
)

type (
	Client interface {
		login(usr, pw string) (string, error)
		LoadInto(data []interface{}) error
		LoadFrom(userid string) ([]interface{}, error)
	}
	PlatformClient struct {
		config     *Config
		token      string
		httpClient *http.Client
		user       []map[string]string
	}
	Config struct {
		Auth   string `json:"auth"`
		Upload string `json:"upload"`
		Acct   string `json:"acct"`
	}
)

func NewClient(cfg *Config, usr, pw string) (*PlatformClient, error) {

	client := &PlatformClient{config: cfg, httpClient: &http.Client{}}

	if tkn, err := client.login(usr, pw); err != nil {
		log.Println("Error init client: ", err)
		return nil, err
	} else {
		//log.Println("here it is ", tkn)
		client.token = tkn
		return client, nil
	}
}

func (tc *PlatformClient) login(usr, pw string) (token string, err error) {

	req, err := http.NewRequest("POST", tc.config.Auth+"/login", nil)
	req.SetBasicAuth(usr, pw)
	if resp, err := tc.httpClient.Do(req); err != nil {
		return "", err
	} else {
		if resp.StatusCode == http.StatusOK {
			return resp.Header.Get(TP_SESSION_TOKEN), nil
		}
		return "", errors.New("Issue logging in: " + string(resp.StatusCode))
	}
}

func (tc *PlatformClient) LoadInto(data []interface{}) error {

	jsonBlock, _ := json.Marshal(data)

	log.Println(" block to load ", bytes.NewBufferString(string(jsonBlock)))
	//log.Println(" token ", tc.token)

	req, _ := http.NewRequest("POST", tc.config.Upload, bytes.NewBufferString(string(jsonBlock)))
	req.Header.Add(TP_SESSION_TOKEN, tc.token)
	req.Header.Set("content-type", "application/json")

	if resp, err := tc.httpClient.Do(req); err != nil {
		log.Println("Error loading messages: ", err)
		return err
	} else {
		log.Printf("all good? [%d] [%s] ", resp.StatusCode, resp.Status)
		updatedToken := resp.Header.Get(TP_SESSION_TOKEN)
		if updatedToken != "" && tc.token != updatedToken {
			tc.token = updatedToken
			log.Println("updated the token")
		}
	}

	return nil
}

func (tc *PlatformClient) LoadFrom(userid string) (*[]interface{}, error) {
	return nil, nil
}
