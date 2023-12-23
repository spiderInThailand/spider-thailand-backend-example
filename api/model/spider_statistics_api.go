package model

import "spider-go/model"

// ==================================================
// get spider statistice type
// ==================================================
type SpiderStatisticsRequester interface{}
type SpiderStatisticsResponser struct {
	Header ResponseHeader `json:"header"`
	Data   interface{}    `json:"data"`
}

// **************************************************

// ==================================================
// get family list
// ==================================================
type GetFamilyListRequester struct {
	Header RequestUserHeader        `json:"header"`
	Data   GetFamilyListRequestData `json:"data"`
}

type GetFamilyListRequestData struct {
	Page int32 `json:"page" validate:"min=0"`
	Size int32 `json:"size" validate:"min=1"`
}

type GetFamilyListResponser struct {
	Header ResponseHeader            `json:"header"`
	Data   GetFamilyListResponseData `json:"data"`
}

type GetFamilyListResponseData struct {
	FamilyList []model.FamilyList `json:"family_list"`
}
