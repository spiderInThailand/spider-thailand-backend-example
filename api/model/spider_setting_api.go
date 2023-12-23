package model

// upload spider image
type SpiderImageSettingRequester struct {
	Header RequestUserHeader      `json:"header"`
	Data   SpiderImageSettingData `json:"data"`
}

type SpiderImageSettingResponser struct {
	Header ResponseHeader `json:"header"`
	Data   struct{}       `json:"data"`
}

type SpiderImageSettingData struct {
	SpiderUUID      string   `json:"spider_uuid"`
	ListImageEncode []string `json:"list_image_encode"`
}

// delete spider info
type DeleteSpiderRequester struct {
	Header RequestUserHeader `json:"header"`
	Data   DeleteSpiderData  `json:"data"`
}

type DeleteSpiderData struct {
	SpiderUUID string `json:"spider_uuid"`
}

type DeleteSpiderResponser struct {
	Header ResponseHeader `json:"header"`
	Data   struct{}       `json:"data"`
}

// edit spider info
type EditSpiderInfoRequester struct {
	Header RequestUserHeader `json:"header"`
	Data   SpiderInfo        `json:"data"`
}

type EditSpiderInfoResponser struct {
	Header ResponseHeader `json:"header"`
	Data   struct{}       `json:"data"`
}

// remove spider image
type RemoveSpiderImageRequester struct {
	Header RequestUserHeader     `json:"header"`
	Data   RemoveSpiderImageRequestData `json:"data"`
}

type RemoveSpiderImageRequestData struct {
	SpiderUUID string `json:"spider_uuid"`
	SpiderImageList []string `json:"spider_image_list"`
}

type RemoveSpiderImageResponse struct {
	Header ResponseHeader `json:"header"`
	Data   struct{}       `json:"data"`
}
