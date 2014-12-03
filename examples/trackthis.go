package trackthis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jh-bate/intertidal/backend/platform"
	"github.com/jh-bate/intertidal/backend/store"
	"github.com/jh-bate/intertidal/flood"
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

func (c *Client) Init(config interface{}) *Client {
	c.config = config.(Config)
	return c
}

func (c *Client) Load() *Client {

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
	return c
}

func (c *Client) StashLocal(key string, local store.Client) *Client {

	if len(c.processed) > 0 {

		err := local.StoreUserData(key, c.processed)

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
		return c
	}
	log.Println("No data to send to the platform")
	return c
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

	bob := flood.NewEventBuilder()

	for i := range c.raw {

		switch {
		case strings.Index(strings.ToUpper(c.raw[i].Type), BG) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, bob.BuildBg(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), CARBS) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, bob.BuildFood(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), BOLUS) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, bob.BuildBolus(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(c.raw[i].Type), BASAL) != -1:
			for en := range c.raw[i].Entries {
				c.processed = append(c.processed, bob.BuildBasal(fmt.Sprintf("%.1f", c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		default:
			for en := range c.raw[i].Entries {
				log.Println("Default: ", c.raw[i].Type)
				log.Println("Save: ", fmt.Sprintf("#%s %.1f", c.raw[i].Type, c.raw[i].Entries[en].Value))
				c.processed = append(c.processed, bob.BuildNote(fmt.Sprintf("#%s %.1f", c.raw[i].Type, c.raw[i].Entries[en].Value), c.raw[i].Entries[en].Time, deviceName))
			}
			break
		}

	}
	return
}
