package util

type JsonMessage struct {
	Status  string            `json:"Status"`
	Message string            `json:"Message"`
	Config  map[string]string `json:"Config"`
}
