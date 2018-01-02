package summer_plugin

import (
	"github.com/sirupsen/logrus"
	"github.com/cocotyty/summer"
	"reflect"
	"gopkg.in/macaron.v1"
	"time"
	"net/http"
	"os"
)

func NewLog(name string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"component": name,
	})
}

type LogPlugin struct {
}

// look up the value which field wanted
func (l *LogPlugin) Look(h *summer.Holder, path string, sf *reflect.StructField) reflect.Value {

	if path == "*" {
		h.Basket.EachHolder(func(name string, holder *summer.Holder) bool {
			if holder == h {
				path = name
				return true
			}
			return false
		})
	}
	mlog := NewLog(path)
	return reflect.ValueOf(mlog)
}

// tell  summer the plugin prefix
func (l *LogPlugin) Prefix() string {
	return "~"
}

// zIndex represent the sequence of plugins
func (l *LogPlugin) ZIndex() int {
	return 0
}

var macaronEntry *logrus.Entry

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv("debug") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	macaronEntry = logrus.WithField("component", "macaron")
}
func init() {
	summer.PluginRegister(&LogPlugin{}, summer.BeforeInit)
}

func MacaronLogger() func(ctx *macaron.Context) {
	return func(ctx *macaron.Context) {
		start := time.Now()
		rw := ctx.Resp.(macaron.ResponseWriter)
		ctx.Next()
		macaronEntry.WithFields(logrus.Fields{
			"method":     ctx.Req.Method,
			"requestURI": ctx.Req.RequestURI,
			"remoteAddr": ctx.RemoteAddr(),
			"status":     rw.Status(),
			"statusText": http.StatusText(rw.Status()),
			"use":        time.Since(start),
		}).Info()
	}

}
