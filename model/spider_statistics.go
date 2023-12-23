package model

import "time"

type SpiderStatistics struct {
	FamilyName string       `json:"family_name" bson:"family_name"`
	Genus      []GenusGroup `json:"genus" bson:"genus"`
	CreatedAt  time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at" bson:"updated_at"`
}

type GenusGroup struct {
	GenusName string         `json:"genus_name" bson:"genus_name"`
	Species   []SpeciesGroup `json:"species" bson:"species"`
}

type SpeciesGroup struct {
	SpeciesName string `json:"species_name" bson:"species_name"`
}

type FamilyList struct {
	Family   string `json:"family"`
	Author   string `json:"author"`
	Quantity int32  `json:"quantity"`
}
