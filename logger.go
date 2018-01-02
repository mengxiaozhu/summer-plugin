package summer_plugin

import (
	"github.com/Sirupsen/logrus"
	"github.com/cocotyty/summer"
	"reflect"
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

func init() {
	summer.PluginRegister(&LogPlugin{}, summer.BeforeInit)
}
