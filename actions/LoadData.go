package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jh-bate/intertidal/data"
)

const (
	baseUrl    = "https://www.trackthisfor.me/api/v1/categories/"
	tokenParam = "?access_token="
	deviceName = "www.trackthisfor.me/api/v1"
	BG         = "BLOOD SUGAR"
	CARBS      = "CARBS"
	BOLUS      = "NOVORAPID"
	BASAL      = "LANTUS"
	TRACK_THIS = "trackthis"
)

type (
	dataLoader struct {
		loaderConfig
		raw []trackThisEntries
	}
	loaderConfig struct {
		AuthToken string `json:"authToken"`
	}
	LoadData struct {
		loader  dataLoader
		Results []interface{}
	}

	trackThisCategories struct {
		Codes []struct{ Id int } `json:"categories"`
	}

	trackThisEntries struct {
		Type    string           `json:"name"`
		Units   string           `json:"symbol"`
		Entries []trackThisEntry `json:"elements"`
	}

	trackThisEntry struct {
		Value   float32 `json:"value"`
		Time    string  `json:"date"`
		Comment string  `json:"comment"`
	}
)

func LoadDataAction(calldata interface{}, store data.Store, name string) *LoadData {

	config := calldata.([]byte)

	var jsonCfg loaderConfig
	if err := json.Unmarshal(config, &jsonCfg); err != nil {
		log.Panic("the given config is invalid ", err.Error())
	}

	loader := dataLoader{loaderConfig: jsonCfg}

	return &LoadData{loader: loader}
}

func (a *LoadData) Execute() (interface{}, error) {

	a.load()
	return a.Results, nil
}

func (l *LoadData) load() {

	log.Println("loading from trackthisforme.com ...")

	url := baseUrl + tokenParam + l.loader.loaderConfig.AuthToken

	if res, err := http.Get(url); err != nil {
		log.Println("Error getting data from trackthisforme.com: ", err)
	} else {
		if data, err := ioutil.ReadAll(res.Body); err != nil {
			log.Println("Error reading data from trackthisforme.com: ", err)
		} else {
			res.Body.Close()
			var categories trackThisCategories

			if mErr := json.Unmarshal(data, &categories); mErr != nil {
				log.Println("Error parsing data from trackthisforme.com: ", mErr)
			} else {

				for i := range categories.Codes {
					l.loader.raw = append(l.loader.raw, l.loadCategory(strconv.Itoa(categories.Codes[i].Id)))
				}

				log.Println("loaded data from trackthisforme.com")
			}
		}
	}
	l.transform()
	return
}

func (l *LoadData) loadCategory(categoryId string) trackThisEntries {

	url := baseUrl + categoryId + ".json" + tokenParam + l.loader.loaderConfig.AuthToken

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	bgs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println(err)
	}

	var categoryEntries trackThisEntries

	mErr := json.Unmarshal(bgs, &categoryEntries)
	if mErr != nil {
		log.Println(mErr)
	}

	return categoryEntries
}

func (l *LoadData) transform() {

	bob := data.NewEventBuilder()

	for i := range l.loader.raw {

		switch {
		case strings.Index(strings.ToUpper(l.loader.raw[i].Type), BG) != -1:
			for en := range l.loader.raw[i].Entries {
				l.Results = append(l.Results, bob.BuildBg(fmt.Sprintf("%.1f", l.loader.raw[i].Entries[en].Value), l.loader.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(l.loader.raw[i].Type), CARBS) != -1:
			for en := range l.loader.raw[i].Entries {
				l.Results = append(l.Results, bob.BuildFood(fmt.Sprintf("%.1f", l.loader.raw[i].Entries[en].Value), l.loader.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(l.loader.raw[i].Type), BOLUS) != -1:
			for en := range l.loader.raw[i].Entries {
				l.Results = append(l.Results, bob.BuildBolus(fmt.Sprintf("%.1f", l.loader.raw[i].Entries[en].Value), l.loader.raw[i].Entries[en].Time, deviceName))
			}
			break
		case strings.Index(strings.ToUpper(l.loader.raw[i].Type), BASAL) != -1:
			for en := range l.loader.raw[i].Entries {
				l.Results = append(l.Results, bob.BuildBasal(fmt.Sprintf("%.1f", l.loader.raw[i].Entries[en].Value), l.loader.raw[i].Entries[en].Time, deviceName))
			}
			break
		default:
			for en := range l.loader.raw[i].Entries {
				log.Println("Default: ", l.loader.raw[i].Type)
				log.Println("Save: ", fmt.Sprintf("#%s %.1f", l.loader.raw[i].Type, l.loader.raw[i].Entries[en].Value))
				l.Results = append(l.Results, bob.BuildNote(fmt.Sprintf("#%s %.1f", l.loader.raw[i].Type, l.loader.raw[i].Entries[en].Value), l.loader.raw[i].Entries[en].Time, deviceName))
			}
			break
		}

	}
	return
}
