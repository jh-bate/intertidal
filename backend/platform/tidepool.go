package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jh-bate/intertidal/store"
)

const (
	TP_SESSION_TOKEN = "x-tidepool-session-token"
)

type (
	Client interface {
		LoadFrom(userid string) ([]interface{}, error)
		LoadInto(data []interface{}) error
		login(usr, pw string) (string, error)
		Profile(name string) error
		Signup(name, pw, contact string) error
		StashUserLocal(local store.Client)
	}
	TidepoolClient struct {
		config     *Config
		token      string
		httpClient *http.Client
		User       *User
	}
	User struct {
		Token   string `json:"-"`
		Id      string `json:"-"`
		Pw      string `json:"password"`
		Name    string `json:"username"`
		Contact string `json:"emails"`
		Profile profile
	}
	profile struct {
		FullName string `json:"fullName"`
	}
	Config struct {
		Auth   string `json:"auth"`
		Upload string `json:"upload"`
	}

	Common struct {
		Type     string `json:"type"`
		DeviceId string `json:"deviceId"`
		Source   string `json:"source"`
		Time     string `json:"time"`
	}
	BgEvent struct {
		Common
		Value float64 `json:"value"`
	}
	FoodEvent struct {
		Common
		Carbs float64 `json:"carbs"`
	}
	BasalEvent struct {
		Common
		DeliveryType string  `json:"deliveryType"`
		Value        float64 `json:"value"`
		Duration     int     `json:"duration"`
		Insulin      string  `json:"insulin"`
	}
	BolusEvent struct {
		Common
		SubType string  `json:"subType"`
		Value   float64 `json:"value"`
		Insulin string  `json:"insulin"`
	}
	NoteEvent struct {
		Common
		CreatorId string `json:"creatorId"`
		Text      string `json:"text"`
	}
)

func NewClient(cfg *Config, usrName, pw string) *PlatformClient {

	client := &PlatformClient{config: cfg, httpClient: &http.Client{}}

	if tkn, err := client.login(usrName, pw); err != nil {
		log.Panicf("Error init client: ", err)
		return nil
	} else {
		client.User = &User{Token: tkn, Name: usrName}
		log.Printf("user [%v]", client.User)
		return client
	}
}

func (c *PlatformClient) StashUserLocal(local store.Client) {

	err := local.StoreUser(c.User.Token, c.User.Name)

	if err != nil {
		log.Println("Error statshing data ", err)
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

func (tc *PlatformClient) Signup(name, pw, contact string) error {

	newUsr := &User{Name: name, Pw: pw, Contact: contact}
	jsonBlock, _ := json.Marshal(newUsr)

	req, _ := http.NewRequest("POST", tc.config.Upload, bytes.NewBufferString(string(jsonBlock)))
	req.Header.Set("content-type", "application/json")

	if resp, err := tc.httpClient.Do(req); err != nil {
		log.Println("Error doing signup: ", err)
		return err
	} else {

		var signup struct {
			id string `json:"userid"`
		}

		defer req.Body.Close()
		if err := json.NewDecoder(req.Body).Decode(&signup); err != nil {
			log.Printf("Error after trying to signup: %v", err)
			return err
		}

		newUsr.Id = signup.id
		newUsr.Token = resp.Header.Get(TP_SESSION_TOKEN)
		tc.User = newUsr

		if err := tc.Profile(name); err != nil {
			log.Println("Error adding profile for newly signedup user")
		}

		return nil
	}
}

func (tc *PlatformClient) Profile(name string) error {

	type profile struct {
		FullName string `json:"fullName"`
	}

	prf := &profile{FullName: name}
	jsonBlock, _ := json.Marshal(prf)

	req, _ := http.NewRequest("POST", tc.config.Upload, bytes.NewBufferString(string(jsonBlock)))
	req.Header.Add(TP_SESSION_TOKEN, tc.token)
	req.Header.Set("content-type", "application/json")

	if resp, err := tc.httpClient.Do(req); err != nil {
		log.Println("Error doing signup: ", err)
		return err
	} else {

		if resp.StatusCode == http.StatusOK {
			tc.User.Profile.FullName = name
		} else {
			log.Printf("Issue adding profile [%s] [%s]", resp.StatusCode, resp.Status)
		}

		return nil
	}
}

func (tc *PlatformClient) LoadInto(data []interface{}) error {

	jsonBlock, _ := json.Marshal(data)

	log.Println(" block to load ", bytes.NewBufferString(string(jsonBlock)))
	//log.Println(" token ", tc.token)

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

func (tc *PlatformClient) LoadFrom(userid string) ([]interface{}, error) {
	return nil, nil
}
