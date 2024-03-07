package config

type Config struct {
	Api Api `json:"api"`
}

type Api struct {
	Addr string `json:"addr"`
}
