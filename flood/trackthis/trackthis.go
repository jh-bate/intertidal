package trackthis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jh-bate/intertidal/flood"
	"github.com/jh-bate/intertidal/platform"
	"github.com/jh-bate/intertidal/store"
)

const (
	baseUrl    = "https://www.trackthisfor.me/api/v1/categories/"
	tokenParam = "?access_token="
	deviceName = "www.trackthisfor.me/api/v1"
	BG         = "BLOOD SUGAR"
	CARBS      = "CARBS"
	BOLUS      = "NOVORAPID"
	BASAL      = "LANTUS"
)

type (
	Client struct {
		config    Config
		raw       []TrackThisEntries
		processed []interface{}
		store     store.Client
	}

	Config struct {
		AuthToken string `json:"authToken"`
	}

	TrackThisCategories struct {
		Codes []struct{ Id int } `json:"categories"`
	}

	TrackThisEntries struct {
		Type    string           `json:"name"`
		Units   string           `json:"symbol"`
		Entries []TrackThisEntry `json:"elements"`
	}

	TrackThisEntry struct {
		Value   float32 `json:"value"`
		Time    string  `json:"date"`
		Comment string  `json:"comment"`
	}
)

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Init(config interface{}, store store.Client) {

	c.config = config.(Config)
	c.store = store
}

func (c *Client) Load() {

	log.Println("loading from trackthisforme.com ...")

	url := baseUrl + tokenParam + c.config.AuthToken

	if res, err := http.Get(url); err != nil {
		log.Println("Error getting data from trackthisforme.com: ", err)
	} else {
		if data, err := ioutil.ReadAll(res.Body); err != nil {
			log.Println("Error reading data from trackthisforme.com: ", err)
		} else {
			res.Body.Close()
			var categories TrackThisCategories

			if mErr := json.Unmarshal(data, &categories); mErr != nil {
				log.Println("Error parsing data from trackthisforme.com: ", mErr)
			} else {

				for i := range categories.Codes {
					c.raw = append(c.raw, c.loadCategory(strconv.Itoa(categories.Codes[i].Id)))
				}

				log.Println("loaded data from trackthisforme.com")
			}
		}
	}
	c.transform()
	return
}

func (c *Client) StashLocal() {

	if len(c.processed) > 0 {

		err := c.store.StoreData("123", c.processed)

		if err != nil {
			log.Println("Error statshing data ", err)
		}
		return
	}
	log.Println("No data to stash")
	return
}

func (c *Client) StorePlatform(platform platform.Client) {

	if len(c.processed) > 0 {

		err := platform.LoadInto(c.processed)

		if err != nil {
			log.Println("Error sending to platform ", err)
		}
	}
	log.Println("No data to send to the platform")
}

func (c *Client) loadCategory(categoryId string) TrackThisEntries {

	url := baseUrl + categoryId + ".json" + tokenParam + c.config.AuthToken

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	bgs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
	}

	var category TrackThisEntries

	mErr := json.Unmarshal(bgs, &category)
	if mErr != nil {
		log.Println(mErr)
	}

	return category
}

func (c *Client) transform() {
	for i := range c.raw {

		switch {
		case strings.Index(strings.ToUpper(c.raw[i].Type), BG) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, flood.MakeBg(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), CARBS) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, flood.MakeCarb(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), BOLUS) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, flood.MakeBolus(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), BASAL) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, flood.MakeBasal(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		default:
			for en := range c.raw[i].Entries {
				log.Println("Default: ", c.raw[i].Type)
				log.Println("Save: ", fmt.Sprintf("#%s %.1f", c.raw[i].Type, c.raw[i].Entries[en].Value))
				c.processed = append(c.processed, flood.MakeNote(fmt.Sprintf("#%s %.1f", c.raw[i].Type, c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		}

	}
	return
}
