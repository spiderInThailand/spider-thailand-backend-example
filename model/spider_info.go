package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	SPIDER_INFO_STATUS_ACTIVE   = "active"
	SPIDER_INFO_STATUS_INACTIVE = "inactive"
)

type SpiderInfo struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SpiderUUID   string             `json:"spider_uuid" bson:"spider_uuid"`
	Family       string             `json:"family" bson:"family"`
	Genus        string             `json:"genus" bson:"genus"`
	Species      string             `json:"species" bson:"species"`
	Author       string             `json:"author" bson:"author,omitempty"`
	PublishYear  string             `json:"publish_year" bson:"publish_year,omitempty"`
	Country      string             `json:"country" bson:"country,omitempty"`
	CountryOther string             `json:"country_other" bson:"country_other,omitempty"`
	Altitude     string             `json:"altitude" bson:"altitude,omitempty"`
	Method       string             `json:"method" bson:"method,omitempty"`
	Habital      string             `json:"habitat" bson:"habitat,omitempty"`
	Microhabital string             `json:"microhabitat" bson:"microhabitat,omitempty"`
	Designate    string             `json:"designate" bson:"designatel,omitempty"`
	Address      []Address          `json:"address" bson:"address,omitempty"`
	Paper        []string           `json:"paper" bson:"paper,omitempty"`
	Status       string             `json:"status" bson:"status"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	ImageFile    []string           `json:"image_file" bson:"image_file,omitempty"`
	CreatedBy    string             `json:"created_by" bson:"created_by,omitempty"`
}

type Address struct {
	Province string     `json:"province" bson:"province"`
	District string     `json:"district" bson:"district"`
	Locality string     `json:"locality" bson:"locality"`
	Position []Position `json:"position" bson:"position"`
}

type Position struct {
	Name      string  `json:"name" bson:"name"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type LocationResult struct {
	Province string           `json:"province"`
	Locality []LocalityResult `json:"locality"`
}

type LocalityResult struct {
	Name         string   `json:"name"`
	SubLocaltion []string `json:"sub_localtion"`
}

type GetSpiderListBySpiderTypeParam struct {
	Family  string
	Genus   string
	Species string
	Page    int32
	Size    int32
}

type SpiderImageList struct {
	Name        string `json:"name"`
	ImageBase64 string `json:"image_base64"`
}
