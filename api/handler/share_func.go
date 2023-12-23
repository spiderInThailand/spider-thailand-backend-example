package handler

import (
	api_model "spider-go/api/model"
	"spider-go/model"
)

func mapSpiderInfoModel(data *model.SpiderInfo) *api_model.SpiderInfo {

	var address []api_model.Address

	for _, v := range data.Address {
		tempAddress := api_model.Address{
			Province: v.Province,
			District: v.District,
			Locality: v.Locality,
		}

		var tempPosition []api_model.Position

		for _, v2 := range v.Position {
			thisPosition := api_model.Position{
				Latitude:  v2.Latitude,
				Longitude: v2.Longitude,
				Name:      v2.Name,
			}

			tempPosition = append(tempPosition, thisPosition)
		}

		tempAddress.Position = tempPosition

		address = append(address, tempAddress)
	}

	RespSpiderInfo := api_model.SpiderInfo{
		SpiderUUID:     data.SpiderUUID,
		Family:         data.Family,
		Genus:          data.Genus,
		Species:        data.Species,
		Author:         data.Author,
		PublishYear:    data.PublishYear,
		Country:        data.Country,
		OtherCountries: data.CountryOther,
		Altitude:       data.Altitude,
		Method:         data.Method,
		Habital:        data.Habital,
		Microhabital:   data.Microhabital,
		Designate:      data.Designate,
		Address:        address,
		Paper:          data.Paper,
		Image:          data.ImageFile,
	}

	return &RespSpiderInfo
}
