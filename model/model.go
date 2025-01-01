package model

type upstreams struct {
	IpAddress string `json:"ipaddress"`
	Port      string `json:"port"`
}

type Rules struct {
	Upstreams  []upstreams `json:"upstreams"`
	Algorithm  string      `json:"algorithm"`
	ServerFarm string      `json:"server-farm"`
}
