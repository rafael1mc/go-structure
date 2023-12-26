package consts

import "log/slog"

func SlogError(err error) slog.Attr {
	return slog.Any("error", err)
}
