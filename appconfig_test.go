package appconfig_test

import (
	"github.com/ndphu/espresso.appconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	FirebaseTestApp string = "test-7a4ff"
)

func TestGetExistingConfigFromFirebase(t *testing.T) {
	keyFile, exists := os.LookupEnv("KEY_FILE")
	if !exists {
		keyFile = "key.json"
	}

	expected := appconfig.New()
	expected.Device.Id = "0"
	expected.Server.MQTT.Host = "19november.ddns.net"
	expected.Server.MQTT.Protocol = "tcp"
	expected.Server.MQTT.Port = 5384
	expected.Server.Firebase.AppName = FirebaseTestApp
	expected.Schema = "1.0"

	actual := appconfig.New()
	actual.GetConfigFromFirebase(FirebaseTestApp, "0", keyFile)

	assert.Equal(t, expected, actual, "Config should be retreive correctly.")
}
