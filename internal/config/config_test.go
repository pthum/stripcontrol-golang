package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoad(t *testing.T) {
	testConf := `
server:
  host: localhost
  port: 8080
  mode: debug
database:
  type: sqlite
  user: sa
  pass: password
  host: stripcontrol.sqlite
  port: 5432
  name: stripcontrol
messaging:
  host: mqtthost
  port: 1234
  striptopic: ledstripz
  profiletopic: profilez
  disabled: true
`
	conf := &Config{}
	err := conf.readConf([]byte(testConf))
	assert.Nil(t, err)
	assert.Equal(t, "localhost", conf.Server.Host)
	assert.Equal(t, "8080", conf.Server.Port)
	assert.Equal(t, "debug", conf.Server.Mode)
	assert.Equal(t, "mqtthost", conf.Messaging.Host)
	assert.Equal(t, "1234", conf.Messaging.Port)
	assert.Equal(t, "ledstripz", conf.Messaging.StripTopic)
	assert.Equal(t, "profilez", conf.Messaging.ProfileTopic)
	assert.Equal(t, true, conf.Messaging.Disabled)
}

func TestConfigLoadError(t *testing.T) {
	// tabs are not allowed in yaml
	testConf := `
server:
	host: localhost
	port: 8080
	mode: debug
`
	conf := &Config{}
	err := conf.readConf([]byte(testConf))
	assert.NotNil(t, err)
}

func TestConfigLoadNonExistingFile(t *testing.T) {
	// tabs are not allowed in yaml
	testFile := "/shouldnotexist"
	_, err := InitConfig(testFile)
	assert.NotNil(t, err)
}
