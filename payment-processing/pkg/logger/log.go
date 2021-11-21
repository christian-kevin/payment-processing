package log

import (
"context"
"github.com/rs/xid"
"gitlab.com/opaper/goutils/log"
)

var logger log.SLogger

func InitializeLogger(appName string) {
	logger = log.NewSLogger(appName)
}

func Get() log.SLogger {
	return logger
}

func GetEmptyContext() context.Context {
	return logger.BuildContextDataAndSetValue("", "")
}

func BuildContextLoggerWithCountry(ctx context.Context, country string) context.Context {
	data := make(map[string]string, 0)
	if value := ctx.Value(log.ContextDataMapKey); value != nil {
		if dataMap, ok := value.(map[string]string); ok {
			data[log.ContextIdKey] = dataMap[log.ContextIdKey]
		}
	} else {
		guid := xid.New()
		data[log.ContextIdKey] = guid.String()
	}

	data[log.ContextCountryKey] = country

	return context.WithValue(ctx, log.ContextDataMapKey, data)
}

func BuildContextLoggerWithRequestID(ctx context.Context, requestID string) context.Context {
	data := make(map[string]string, 0)
	if value := ctx.Value(log.ContextDataMapKey); value != nil {
		if dataMap, ok := value.(map[string]string); ok {
			data[log.ContextCountryKey] = dataMap[log.ContextCountryKey]
		}
	}

	data[log.ContextIdKey] = requestID

	return context.WithValue(ctx, log.ContextDataMapKey, data)
}
