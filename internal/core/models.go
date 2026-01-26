package core

type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LocationData struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Accuracy float64 `json:"accuracy"`
	Speed    float64 `json:"speed"`
}

type State struct {
	ID       string       `json:"id"`
	LastSeen int64        `json:"lastSeen"`
	Location LocationData `json:"location"`
	Battery  struct {
		Level      float64 `json:"level"`
		IsCharging bool    `json:"isCharging"`
	} `json:"battery"`
}

type InitialStatePayload struct {
	Profiles []Profile `json:"profiles"`
	States   []State   `json:"states"`
}
