package plugin

import (
	"log"
	"testing"
)

func TestNewPluginCenter(t *testing.T) {
	pluginCenter := NewPluginCenter()
	pluginCenter.Register("001", "../plugins/count.so")
	t.Log("ok")

	go pluginCenter.Plugins["001"].Run(nil, nil)
	config := pluginCenter.Plugins["001"].Config
	log.Println("plugin config", config)
}
