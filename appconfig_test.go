package appconfig_test

import (
	"github.com/ndphu/espresso.appconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetExistingConfigFromFirebase(t *testing.T) {
	keyFile, exists := os.LookupEnv("KEY_FILE")
	if !exists {
		t.Error("KEY_FILE is undefined")
		return
	}

	expected := appconfig.New()
	expected.Device.Id = "0"
	expected.Server.Host = "19november.ddns.net"
	expected.Server.Protocol = "tcp"
	expected.Server.Port = 5384
	expected.Schema = "1.0"

	actual := appconfig.New()
	actual.GetConfigFromFirebase("test-7a4ff", "0", keyFile)

	assert.Equal(t, expected, actual, "Config should be retreive correctly.")
}
