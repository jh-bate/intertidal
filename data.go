package intertidal

const (
	DATA_COLLECTION = "data"
)

/*
 * Types that match what the tidepool platform expects
 */
type (
	//common types that all others will inherit
	Common struct {
		Type     string `json:"type"`
		DeviceId string `json:"deviceId"`
		UploadId string `json:"uploadId"`
		Source   string `json:"source"`
		Time     string `json:"time"`
	}
	//a simple BG value
	BgEvent struct {
		Common
		Value float64 `json:"value"`
	}
	//a food event measured by carbs consumed
	FoodEvent struct {
		Common
		Carbs float64 `json:"carbs"`
	}
	//long acting or background insulin
	BasalEvent struct {
		Common
		DeliveryType string  `json:"deliveryType"`
		Value        float64 `json:"value"`
		Duration     int     `json:"duration"`
		Insulin      string  `json:"insulin"`
	}
	//short acting insulin
	BolusEvent struct {
		Common
		SubType string  `json:"subType"`
		Value   float64 `json:"value"`
		Insulin string  `json:"insulin"`
	}
	//a note to add context to data
	NoteEvent struct {
		Common
		CreatorId string `json:"creatorId"`
		Text      string `json:"text"`
	}
	UploadEvent struct {
		Common
		UploadId string `json:"uploadId"`
		TimeZone string `json:"timezoneName"`
		Text     string `json:"text"`
	}
	/*
			var uploadMeta = {
		          type: 'upload',
		          time: sessionInfo.start,
		          timezone: sessionInfo.tzName,
		          version: sessionInfo.version,
		          deviceId: sessionInfo.deviceId,
		          uploadId: generatedId,
		          byUser: myUserId,
		          source : 'tidepool'
		        };
	*/
)
