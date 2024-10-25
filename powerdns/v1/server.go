package powerdns

// ref: https://doc.powerdns.com/authoritative/http-api/server.html#server
type Server struct {
	AutoPrimariesURL string `json:"autoprimaries_url"`
	ConfigURL        string `json:"config_url"`
	DaemonType       string `json:"daemon_type"`
	ID               string `json:"id"`
	Type             string `json:"server"`
	URL              string `json:"url"`
	Version          string `json:"version"`
	ZonesURL         string `json:"zones_url"`
}
