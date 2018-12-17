package plugin

import (
	"testing"
)

func TestNewPluginCenter(t *testing.T) {
	pluginCenter := NewPluginCenter()
	pluginCenter.Register("001", "../plugins/count.so")
	t.Log("ok")

	go pluginCenter.Plugins["001"].Run(nil, nil)
}
