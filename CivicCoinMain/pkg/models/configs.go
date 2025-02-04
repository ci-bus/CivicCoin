package models

type Configs struct {
	Keys struct {
		Me    string   `json:"me"`
		Nodes []string `json:"nodes"`
	} `json:"keys"`
	Websocket struct {
		Address string `json:"address"`
	} `json:"websocket"`
	Redis struct {
		Addr string `json:"addr"`
		Pass string `json:"pass"`
		Db   int    `json:"db"`
	} `json:"redis"`
}
