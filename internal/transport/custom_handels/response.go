package customhandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ARM-auto-place/auth/internal/metrics"
	customErrors "github.com/ARM-auto-place/auth/pkg/models/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var DefaultInternalErrorText = "internal server error"

type RestResponse struct {
	Data             any            `json:"data"`
	Error            bool           `json:"error"`
	ErrorText        string         `json:"errorText"`
	AdditionalErrors map[string]any `json:"additionalErrors"`
}

func SendResponse(ctx *fiber.Ctx, log zerolog.Logger, data any, respError error) {
	ctx.Response().Header.Set("Content-Type", "application/json")
	ctx.Status(http.StatusOK)

	response := &RestResponse{
		Data:  data,
		Error: respError != nil,
	}

	if response.Error {
		ctx.Response().SetStatusCode(http.StatusInternalServerError)
		response.ErrorText = DefaultInternalErrorText

		var customErr *customErrors.Error
		if errors.As(respError, &customErr) {
			if customErr.GetPublicText() != "" {
				response.ErrorText = customErr.GetPublicText()
			}

			if customErr.GetStatusCode() != 0 {
				ctx.Response().SetStatusCode(customErr.GetStatusCode())
			}
		}

		if customErr != nil {
			response.AdditionalErrors = make(map[string]any)
			response.AdditionalErrors["reason"] = getPublicReason(customErr)
			if response.AdditionalErrors["reason"] == "" {
				response.AdditionalErrors["reason"] = customErr.ErrorText
			}
		}
	}

	respBody, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal response")
		return
	}

	if _, err = ctx.Write(respBody); err != nil {
		log.Error().Err(err).Msg("failed to send response")
		return
	}
}

func getPublicReason(customErr *customErrors.Error) string {
	if len(customErr.Cause) == 0 {
		return customErr.ErrorText
	}

	keys := make([]string, 0, len(customErr.Cause))
	for key := range customErr.Cause {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, "; ")
}
func HandelResponse(
	ctx *fiber.Ctx,
	method string,
	path string,
	methodName string,
	serviceName string,
	fields map[string]interface{},
	call func() (interface{}, error),
) error {
	var err error
	metrics := metrics.NewMetrics(serviceName)
	defer func(begin time.Time) {
		fields["method"] = method
		fields["path"] = path
		fields["methodName"] = methodName
		fields["serviceName"] = serviceName
		fields["took"] = time.Since(begin).String()

		if err != nil {
			if errors.Is(err, customErrors.ForbiddenError()) {
				log.Warn().Err(err).Fields(fields).Msg("call")
			} else {
				log.Error().Err(err).Fields(fields).Msg("call")
			}
			return
		}
		log.Info().Fields(fields).Msg("call")

		metrics.RequestLatency.WithLabelValues(
			serviceName,
			methodName,
			fmt.Sprint(err == nil),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())
	defer metrics.HttpCollector.WithLabelValues(
		serviceName,
		methodName,
		fmt.Sprint(err == nil),
	).Add(1)
	data, err := call()
	if err != nil {
		SendResponse(ctx, log.Logger, nil, err)
		return nil
	}

	SendResponse(ctx, log.Logger, data, nil)
	return nil
}
