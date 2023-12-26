package middleware

import (
	"errors"
	errorutil "gomodel/internal/shared/util/error"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpErrorHandler struct {
	statusCodes map[error]int
}

func NewHttpErrorHandler(
	errorStatusCodeMaps map[error]int,
) *HttpErrorHandler {
	return &HttpErrorHandler{
		statusCodes: errorStatusCodeMaps,
	}
}

func (h *HttpErrorHandler) getStatusCode(err error) int {
	for k, v := range h.statusCodes {
		if errors.Is(err, k) {
			return v
		}
	}

	return http.StatusInternalServerError
}

func (h HttpErrorHandler) Handler(err error, c echo.Context) {
	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		if httpErr.Internal != nil {
			if internalErr, ok := httpErr.Internal.(*echo.HTTPError); ok {
				httpErr = internalErr
			}
		}
	} else {
		httpErr = &echo.HTTPError{
			Code:    h.getStatusCode(err),
			Message: errorutil.UnwrapRecursive(err).Error(),
		}
	}

	code := httpErr.Code
	message := httpErr.Message
	if _, ok := httpErr.Message.(string); ok {
		message = map[string]any{"message": err.Error()}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(httpErr.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}
