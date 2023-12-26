package ping

import (
	"gomodel/internal/shared/timeprovider"
)

type Ping struct {
	timeProvider timeprovider.TimeProvider
}

func NewPing(
	timeProvider timeprovider.TimeProvider,
) *Ping {
	return &Ping{
		timeProvider: timeProvider,
	}
}

func (p Ping) Respond() PingResponse {
	return PingResponse{
		Message: "Pong",
		Time:    p.timeProvider.ProvideUTCMilli(),
	}
}
