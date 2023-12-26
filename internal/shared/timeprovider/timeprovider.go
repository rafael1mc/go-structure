package timeprovider

import "time"

type TimeProvider interface {
	ProvideUTCMilli() int64
	ProvideUTCSec() int64
}

type TimeProviderImpl struct{}

func NewTimeProviderImpl() TimeProvider {
	return TimeProviderImpl{}
}

func (t TimeProviderImpl) ProvideUTCMilli() int64 {
	return time.Now().UTC().UnixMilli()
}

func (t TimeProviderImpl) ProvideUTCSec() int64 {
	return time.Now().UTC().Unix()
}
