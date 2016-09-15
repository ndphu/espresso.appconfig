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

func GetV1TestConfig() *appconfig.AppConfig {
	expected := appconfig.New()
	expected.Device.Id = "0"
	expected.Server.MQTT.Host = "19november.ddns.net"
	expected.Server.MQTT.Protocol = "tcp"
	expected.Server.MQTT.Port = 5384
	expected.Server.MQTT.User = "someone"
	expected.Server.MQTT.Password = "secret"
	expected.Server.Firebase.AppName = "test-7a4ff"
	expected.Schema = "1.0"

	return expected
}

func GetV2TestConfig() *appconfig.AppConfig {
	expected := appconfig.New()
	expected.Device.Id = "0"
	expected.Server.MQTT.User = "someone"
	expected.Server.MQTT.Password = "secret"
	expected.Server.Firebase.AppName = "test-7a4ff"
	expected.Schema = "2.0"
	expected.Server.MQTT.BrokerUrl = "ws://iot.eclipse.org:80/ws"

	return expected
}

func TestGetExistingConfigFromFirebase(t *testing.T) {
	if testing.Short() {
		t.Skip("")
	}
	keyFile, exists := os.LookupEnv("KEY_FILE")
	if !exists {
		keyFile = "key.json"
	}

	expected := GetV1TestConfig()

	actual := appconfig.New()
	actual.GetConfigFromFirebase(FirebaseTestApp, "0", keyFile)

	assert.Equal(t, expected, actual, "Config should be retreive correctly.")
}

func TestParseV1ConfigFile(t *testing.T) {
	expected := GetV1TestConfig()
	actual := appconfig.New()
	actual.ParseConfigFile("test_v1config.json")

	assert.Equal(t, expected, actual, "V1 Config file parse failure")
}

func TestParseV2ConfigFile(t *testing.T) {
	expected := GetV2TestConfig()

	actual := appconfig.New()
	actual.ParseConfigFile("test_v2config.json")

	assert.Equal(t, expected, actual, "V2 Config file parse failure")

}
