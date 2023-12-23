package model

// ==================================================
// get province
// ==================================================
type GetProvinceRequester struct {
	Header RequestUserHeader `json:"header"`
	Data   struct{}          `json:"data"`
}

type GetProvinceResponser struct {
	Header ResponseHeader `json:"header"`
	Data   ProvinceList   `json:"data"`
}

type ProvinceList struct {
	Province []Province `json:"province_list"`
}

type Province struct {
	Number int16  `json:"number,omitempty"`
	NameTH string `json:"name_th"`
	NameEN string `json:"name_en"`
}

// **************************************************

// ==================================================
// get distance
// ==================================================

type GetDistrictRequester struct {
	Header RequestUserHeader  `json:"header"`
	Data   GetDistrictRequest `json:"data"`
}

type GetDistrictRequest struct {
	ProvinceNameEN string `json:"province_name_en"`
}

type GetDistrictResponser struct {
	Header ResponseHeader `json:"header"`
	Data   DistrictList   `json:"data"`
}

type DistrictList struct {
	District []District `json:"district_list"`
}

type District struct {
	Number int16  `json:"number,omitempty"`
	NameTH string `json:"name_th"`
	NameEN string `json:"name_en"`
}
