package model

type upstreams struct {
	IpAddress string `json:"ipaddress"`
	Name      string `json:"name"`
	Port      string `json:"port"`
}

type Rules struct {
	Upstreams []upstreams `json:"upstreams"`
	Algorithm string      `json:"algorithm"`
}
