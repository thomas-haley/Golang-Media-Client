package util

type ServerConfig struct {
	ALLOW_FILE_READ      bool `json:"ALLOW_FILE_READ"`
	ALLOW_FILE_WRITE     bool `json:"ALLOW_FILE_WRITE"`
	ALLOW_FILE_DOWNLOAD  bool `json:"ALLOW_FILE_DOWNLOAD"`
	ALLOW_FILE_TRANSCODE bool `json:"ALLOW_FILE_TRANSCODE"`
	ENFORCE_WL           bool `json:"ENFORCE_WL"`
	ALLOW_WATCH_PARTY    bool `json:"ALLOW_WATCH_PARTY"`
	ALLOW_ACCESS         bool `json:"ALLOW_ACCESS"`
}
