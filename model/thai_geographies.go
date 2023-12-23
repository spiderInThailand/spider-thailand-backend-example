package model

import "time"

type Province struct {
	NameTH    string     `json:"name_th" bson:"name_th"`
	NameEN    string     `json:"name_en" bson:"name_en"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	Amphure   []District `json:"amphure" bson:"amphure,omitempty"`
}

type District struct {
	NameTH    string    `json:"name_th" bson:"name_th"`
	NameEN    string    `json:"name_en" bson:"name_en"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
