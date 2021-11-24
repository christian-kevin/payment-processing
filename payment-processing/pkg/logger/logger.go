package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SLogger interface {
	BuildContextDataAndSetValue(country string, contextID string) (cctx context.Context)
	GetEntry() *log.Entry
	Infof(ctx context.Context, message string, args ...interface{})
	Errorf(ctx context.Context, message string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	GetContextData(ctx context.Context) *ContextData
}

// safe typing https://golang.org/pkg/context/#WithValue
type contextDataMapKeyType string

// add key here for future request based value
var (
	// data map to contain values
	ContextDataMapKey contextDataMapKeyType = "value"

	// context key data added to map
	ContextCountryKey = "country"
	ContextIdKey      = "context_id"
)

type SLog struct {
	entry *log.Entry
}

// context key data added to map
type ContextData struct {
	Country   string
	ContextId string
}

// map contains value need to be displayed
type fields log.Fields

const (
	maximumCallerDepth int = 25
	minimumCallerDepth int = 1
)

func NewSLogger(service string) SLogger {
	entry, _ := getEntryAndLogger(service)
	return &SLog{entry}
}

func getEntryAndLogger(service string) (*log.Entry, *log.Logger) {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	entry := log.NewEntry(logger)
	entry = entry.WithField("service", service)
	return entry, logger
}

func (l *SLog) BuildContextDataAndSetValue(country string, contextId string) (ctx context.Context) {
	data := make(map[string]string, 0)
	data[ContextCountryKey] = country
	data[ContextIdKey] = contextId

	ctx = context.WithValue(context.Background(), ContextDataMapKey, data)

	return ctx
}

func (l *SLog) GetEntry() *log.Entry {
	return l.entry
}

func (l *SLog) Infof(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, l))).Infof(message, args...)
}

func (l *SLog) Errorf(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, l))).Errorf(message, args...)
}

func (l *SLog) Info(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, l))).Info(args...)
}

func (l *SLog) Error(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, l))).Error(args...)
}

func (l *SLog) GetContextData(ctx context.Context) *ContextData {
	dataMap := ctx.Value(ContextDataMapKey)
	var result *ContextData

	if dataMap != nil {
		if data, ok := dataMap.(map[string]string); ok {
			result = &ContextData{
				Country:   data[ContextCountryKey],
				ContextId: data[ContextIdKey],
			}
		}
	}
	return result
}

func (f fields) getFieldsFromContext(ctx context.Context, l *SLog) {
	countryData := ""
	contextIDdata := ""

	data := l.GetContextData(ctx)
	if data != nil {
		countryData = data.Country
		contextIDdata = data.ContextId
	}

	f[ContextCountryKey] = countryData
	f[ContextIdKey] = contextIDdata
}

func getDefaultData(ctx context.Context, log *SLog) fields {
	data := fields(make(map[string]interface{}))
	data.getCallStackTrace()
	data.getFieldsFromContext(ctx, log)
	return data
}

// currently not supporting caller package named 'log'
func (f fields) getCallStackTrace() {
	f.getCaller(getFrame())
}

func (f fields) getCaller(caller *runtime.Frame) {
	if caller == nil {
		return
	}

	funcVal := caller.Function
	fileVal := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	if funcVal != "" {
		f["func"] = funcVal
	}
	if fileVal != "" {
		f["file"] = fileVal
	}
}

// getCaller retrieves the name of the first non this package calling function
func getFrame() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// hard code here because package runtime Caller behaves differently between golang version
		thisPackageName := "log"
		lenPkgName := len(pkg)
		thisPkgNameIndex := strings.Index(pkg, "/"+thisPackageName)

		// If the caller isn't part of this package, we're done
		if thisPkgNameIndex == -1 || thisPkgNameIndex != lenPkgName-4 {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
