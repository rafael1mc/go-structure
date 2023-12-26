package ping

type PingResponse struct {
	Message string `json:"ping"`
	Time    int64  `json:"time"`
}
