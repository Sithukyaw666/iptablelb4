package model

type Upstreams struct {
	IpAddress string `json:"ipaddress"`
	Port      string `json:"port"`
}

type Request struct {
	Upstreams  []Upstreams `json:"upstreams"`
	Algorithm  string      `json:"algorithm"`
	ServerFarm string      `json:"server_farm"`
	Port       string      `json:"port"`
}
