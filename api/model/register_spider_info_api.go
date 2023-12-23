package model

type RegisterSpiderInfoRequester struct {
	Header RequestUserHeader `json:"header"`
	Data   SpiderInfo        `json:"data"`
}

type RegisterSpiderInfoResponser struct {
	Header ResponseHeader `json:"header"`
	Data   SpiderInfo     `json:"data"`
}

type SpiderInfo struct {
	SpiderUUID     string    `json:"spider_uuid,omitempty"`
	Family         string    `json:"family"`
	Genus          string    `json:"genus"`
	Species        string    `json:"species"`
	Author         string    `json:"author"`
	PublishYear    string    `json:"publish_year"`
	Country        string    `json:"country"`
	OtherCountries string    `json:"other_countries"`
	Altitude       string    `json:"altitude"`
	Method         string    `json:"method"`
	Habital        string    `json:"habital"`
	Microhabital   string    `json:"microhabital"`
	Designate      string    `json:"designate"`
	Address        []Address `json:"address"`
	Paper          []string  `json:"paper"`
	Image          []string  `json:"image"`
}

type Address struct {
	Province string     `json:"province"`
	District string     `json:"district"`
	Locality string     `json:"locality"`
	Position []Position `json:"position"`
}

type Position struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
