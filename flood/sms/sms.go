package sms

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/carlosdp/twiliogo"
	"github.com/jh-bate/intertidal/flood"
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

const (
	//Health
	TIME     = "T="
	ACTIVITY = "A="
	BG       = "G="
	CARB     = "C="
	BASAL    = "L="
	BOLUS    = "S="
	NOTE     = "N="

	//Calcs
	LOG_LOW = "#LG"

	MMOLL = "mmol/L"

	TWILIO = "twilio"
)

type (
	ApiClient interface {
		Configure(config interface{}) error
		Send(sms string) error
		Recieve(account string) ([]interface{}, error)
	}
	Client struct {
		config    Config
		raw       []TextData
		processed []interface{}
		api       ApiClient
	}
	TextData struct {
		text, date, device string
	}
	Config struct {
		UserId     string `json:"-"`
		AccountSid string `json:"accountSid"`
		AuthToken  string `json:"authToken"`
	}
)

var (
	//local testing
	testMessage = twiliogo.Message{
		Sid:         "testsid",
		DateCreated: time.Now().Format(time.RFC3339Nano),
		DateUpdated: time.Now().Format(time.RFC3339Nano),
		DateSent:    time.Now().Format(time.RFC3339Nano),
		AccountSid:  "AC3TestAccount",
		From:        "+15555555555",
		To:          "+16666666666",
		Body:        "G=6.7 C=90 S=10 L=20 #lg",
		NumSegments: "1",
		Status:      "queued",
		Direction:   "outbound-api",
		Price:       "2",
		PriceUnit:   "cents",
		ApiVersion:  "2008-04-01",
		Uri:         "/2010-04-01/Accounts/AC3TestAccount/Messages/testsid.json",
	}
	testMessages = twiliogo.MessageList{
		Messages: []twiliogo.Message{testMessage},
	}
)

func NewClient() *Client {
	return &Client{}
}

func newTextData(text, date, device string) TextData {
	return TextData{
		text:   text,
		date:   date,
		device: device,
	}
}

func (c *Client) Init(config interface{}) *Client {
	c.config = config.(Config)
	return c
}

func (c *Client) AttachApi(api ApiClient) *Client {
	c.api = api
	return c
}

func (c *Client) Load() *Client {

	log.Println("loading from twilio")

	smsClient := new(twiliogo.MockClient) //twilio.NewClient(c.config.AccountSid, c.config.AuthToken)

	messagesJson, _ := json.Marshal(testMessages)
	smsClient.On("get", url.Values{}, smsClient.RootUrl()+"/SMS/Messages.json").Return(messagesJson, nil)

	if messages, err := twiliogo.GetMessageList(smsClient); err != nil {
		log.Println(err)
	} else {
		for i := range messages.Messages {
			msg := messages.Messages[i]
			c.raw = append(c.raw, newTextData(msg.Body, msg.DateSent, msg.From))
		}
	}
	c.transform()
	return c
}

func (c *Client) transform() {
	log.Println("transform text from sms client")

	for i := range c.raw {

		smsTxt := strings.Split(c.raw[i].text, " ")

	outer:
		for en := range smsTxt {

			log.Println("text ", smsTxt[en])

			switch {
			case strings.Index(strings.ToUpper(smsTxt[en]), BG) != -1:
				bg := strings.Split(smsTxt[en], BG)
				c.processed = append(c.processed, flood.MakeBg(bg[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), CARB) != -1:
				carb := strings.Split(smsTxt[en], CARB)
				c.processed = append(c.processed, flood.MakeCarb(carb[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), BASAL) != -1:
				basal := strings.Split(smsTxt[en], BASAL)
				c.processed = append(c.processed, flood.MakeBasal(basal[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), BOLUS) != -1:
				bolus := strings.Split(smsTxt[en], BOLUS)
				c.processed = append(c.processed, flood.MakeBolus(bolus[1], c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), LOG_LOW) != -1:
				//hard code 'LOW'
				c.processed = append(c.processed, flood.MakeBg("3.9", c.raw[i].date, c.raw[i].device))
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), ACTIVITY) != -1:
				log.Println("Will be an activity ", c.raw[i])
				break
			case strings.Index(strings.ToUpper(smsTxt[en]), NOTE) != -1:
				c.processed = append(c.processed, flood.MakeNote(smsTxt[en], c.raw[i].date, c.raw[i].device))
				break
			default:
				c.processed = append(c.processed, flood.MakeNote(c.raw[i].text, c.raw[i].date, c.raw[i].device))
				break outer
			}
		}
	}
	return
}

func (c *Client) StashLocal(local store.Client) *Client {

	if len(c.processed) > 0 {

		log.Printf("to stash: [%v]", c.processed)

		err := local.StoreData("999", c.processed)

		if err != nil {
			log.Println("Error statshing data ", err)
		}
		return c
	}
	log.Println("No data to stash")
	return c
}

func (c *Client) StorePlatform(platform platform.Client) *Client {

	if len(c.processed) > 0 {

		err := platform.LoadInto(c.processed)

		if err != nil {
			log.Println("Error sending to platform ", err)
		}
	}
	log.Println("No data to send to the platform")
	return c
}
