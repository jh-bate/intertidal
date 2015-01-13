package store

/*
 * Types that match what the tidepool platform expects
 */
type (
	Common struct {
		Type     string `json:"type"`
		DeviceId string `json:"deviceId"`
		UploadId string `json:"uploadId"`
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
