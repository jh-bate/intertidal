package intertidal

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
	TidepoolStore struct {
		config     *TidepoolConfig
		httpClient *http.Client
		User       *User
	}
	TidepoolConfig struct {
		Auth   string `json:"auth"`
		Upload string `json:"upload"`
		Query  string `json:"query"`
	}
)

func NewTidepoolStore(cfg *TidepoolConfig, usrName, pw string) *TidepoolStore {

	client := &TidepoolStore{config: cfg, httpClient: &http.Client{}, User: &User{Name: usrName, Pw: pw}}

	if err := client.Login(); err != nil {
		log.Panicf("Error init client: ", err)
		return nil
	}
	return client
}

// we need to login to the platform to be able to us it
func (tc *TidepoolStore) Login() (err error) {

	req, err := http.NewRequest("POST", tc.config.Auth+"/login", nil)
	req.SetBasicAuth(tc.User.Name, tc.User.Pw)
	if resp, err := tc.httpClient.Do(req); err != nil {
		return err
	} else {
		if resp.StatusCode == http.StatusOK {
			tc.User.Token = resp.Header.Get(TP_SESSION_TOKEN)
			return nil
		}
		return errors.New("Issue logging in: " + string(resp.StatusCode))
	}
}

func (tc *TidepoolStore) Ping() error {
	req, _ := http.NewRequest("GET", tc.config.Auth+"/status", nil)

	if resp, err := tc.httpClient.Do(req); err != nil {
		return errors.New("Issue with the tidepool platform: " + err.Error())
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("Issue with the tidepool platform: " + string(resp.StatusCode))
	}
	return nil
}

func (tc *TidepoolStore) Save(data []interface{}) error {

	jsonBlock, _ := json.Marshal(data)

	log.Println(" block to load ", bytes.NewBufferString(string(jsonBlock)))

	req, _ := http.NewRequest("POST", tc.config.Upload, bytes.NewBufferString(string(jsonBlock)))
	req.Header.Add(TP_SESSION_TOKEN, tc.User.Token)
	req.Header.Set("content-type", "application/json")

	if resp, err := tc.httpClient.Do(req); err != nil {
		log.Println("Error loading messages: ", err)
		return err
	} else {
		log.Printf("all good? [%d] [%s] ", resp.StatusCode, resp.Status)
		updatedToken := resp.Header.Get(TP_SESSION_TOKEN)
		if updatedToken != "" && tc.User.Token != updatedToken {
			tc.User.Token = updatedToken
			log.Println("updated the token")
		}
	}

	return nil
}

func (tc *TidepoolStore) Run(qry *Query) ([]interface{}, error) {

	if qry.UserId == "" {
		qry.UserId = tc.User.Id
	}

	log.Println(" query to run ", qry.ToString())

	req, _ := http.NewRequest("POST", tc.config.Upload, bytes.NewBufferString(qry.ToString()))
	req.Header.Add(TP_SESSION_TOKEN, tc.User.Token)
	req.Header.Set("content-type", "application/json")

	if resp, err := tc.httpClient.Do(req); err != nil {
		log.Println("Error running query: ", err)
		return nil, err
	} else {
		log.Printf("all good? [%d] [%s] ", resp.StatusCode, resp.Status)
		updatedToken := resp.Header.Get(TP_SESSION_TOKEN)
		if updatedToken != "" && tc.User.Token != updatedToken {
			tc.User.Token = updatedToken
			log.Println("updated the token")
		}
	}

	return nil, nil
}
