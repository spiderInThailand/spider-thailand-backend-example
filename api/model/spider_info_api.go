package model

import "spider-go/model"

// ==================================================
// get one spider info
// ==================================================
type GetOneSpiderInfoRequester struct {
	Header RequestUserHeader    `json:"header"`
	Data   GetOnespiderInfoData `json:"data"`
}

type GetOneSpiderInfoResponsor struct {
	Header ResponseHeader `json:"header"`
	Data   SpiderInfo     `json:"data"`
}

type GetOnespiderInfoData struct {
	SpiderUUID string `json:"spider_uuid"`
}

// **************************************************

// ==================================================
// get image spider
// ==================================================
type GetSpiderImageRequester struct {
	Header RequestUserHeader         `json:"header"`
	Data   GetSpiderImageRequestData `json:"data"`
}

type GetSpiderImageResponsor struct {
	Header ResponseHeader     `json:"header"`
	Data   GetSpiderImageDate `json:"data"`
}

type GetSpiderImageRequestData struct {
	SpiderImageList []string `json:"spider_image_list"`
}
type GetSpiderImageDate struct {
	SpiderImageList []SpiderImageList `json:"spider_image_list"`
}
type SpiderImageList struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}

// **************************************************

// ==================================================
// get spider list manager
// ==================================================
type GetSpiderInfoListManagerRequester struct {
	Header RequestUserHeader                   `json:"header"`
	Data   GetSpiderInfoListManagerRequestData `json:"data"`
}

type GetSpiderInfoListManagerResponsor struct {
	Header ResponseHeader                       `json:"header"`
	Data   GetSpiderInfoListManagerResponseData `json:"data"`
}

type GetSpiderInfoListManagerRequestData struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type GetSpiderInfoListManagerResponseData struct {
	SpiderInfoList []SpiderInfo `json:"spider_info_list"`
}

// **************************************************

// ==================================================
// get spider info list filter geographine
// ==================================================
type GetSpiderInfoByGeographiesRequester struct {
	Header RequestUserHeader                     `json:"header"`
	Data   GetSpiderInfoByGeographiesRequestData `json:"data"`
}

type GetSpiderInfoByGeographieResponsor struct {
	Header ResponseHeader                          `json:"header"`
	Data   GetSpiderInfoByGeographiesResponsetData `json:"data"`
}

type GetSpiderInfoByGeographiesRequestData struct {
	Province string `json:"province"`
	District string `json:"district"`
	Position string `json:"position"`
}

type GetSpiderInfoByGeographiesResponsetData struct {
	SpiderInfoList []SpiderInfo `json:"spider_info_list"`
}

// **************************************************

// ==================================================
// get geographies by spider type
// ==================================================
type GetGeoGraphiesBySpiderTypeRequester struct {
	Header RequestUserHeader                     `json:"header"`
	Data   GetGeoGraphiesBySpiderTypeRequestData `json:"data"`
}

type GetGeoGraphiesBySpiderTypeResponser struct {
	Header ResponseHeader                         `json:"header"`
	Data   GetGeoGraphiesBySpiderTypeResponseData `json:"data"`
}

type GetGeoGraphiesBySpiderTypeRequestData struct {
	Family  string `json:"family"`
	Genus   string `json:"genus"`
	Species string `json:"species"`
}

type GetGeoGraphiesBySpiderTypeResponseData struct {
	LocationResult []model.LocationResult `json:"location_result"`
}

// **************************************************

// ==================================================
// get geographies by spider type
// ==================================================
type GetSpiderInfoByLocalityRequester struct {
	Header RequestUserHeader                  `json:"header"`
	Data   GetSpiderInfoByLocalityRequestData `json:"data"`
}

type GetSpiderInfoByLocalityResponser struct {
	Header ResponseHeader                      `json:"header"`
	Data   GetSpiderInfoByLocalityResponseData `json:"data"`
}

type GetSpiderInfoByLocalityRequestData struct {
	LocalityName string `json:"locality_name" validate:"required"`
	Page         int32  `json:"page" validate:"min=0"`
	Size         int32  `json:"size" validate:"min=1"`
}

type GetSpiderInfoByLocalityResponseData struct {
	SpiderInfoList []SpiderInfo `json:"spider_info_list"`
}

// **************************************************

// ==================================================
// get spider list by spider type
// ==================================================
type GetSpiderListBySpiderTypeRequester struct {
	Header RequestUserHeader             `json:"header"`
	Data   GetSpiderListBySpiderTypeData `json:"data"`
}

type GetSpiderListBySpiderTypeResponser struct {
	Header ResponseHeader                      `json:"header"`
	Data   GetSpiderInfoByLocalityResponseData `json:"data"`
}

type GetSpiderListBySpiderTypeData struct {
	Family  string `json:"family"`
	Genus   string `json:"genus"`
	Species string `json:"species"`
	Page    int32  `json:"page" validate:"min=0"`
	Size    int32  `json:"size" validate:"min=1"`
}
