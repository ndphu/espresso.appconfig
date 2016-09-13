package appconfig

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/ndphu/espresso.helper.firebase"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ServerConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type DeviceConfig struct {
	Id string `json:"id"`
}

type AppConfig struct {
	Schema string `json:"schema"`
	Server ServerConfig
	Device DeviceConfig
}

var (
	confFile    = kingpin.Flag("config-file", "Config file").String()
	keyFile     = kingpin.Flag("key-file", "Google OAuth key file").String()
	logFile     = kingpin.Flag("log-file", "The path to the log file").String()
	deviceId    = kingpin.Flag("device-id", "Required to use online config").String()
	skipSSL     = kingpin.Flag("skip-ssl", "Skip SSL verification. Should be used only in test mode").Bool()
	firebaseApp = kingpin.Flag("firebase-app", "The Firebase app name").Default("rpictl").String()
)

func New() *AppConfig {
	return &AppConfig{}
}

func (appConfig *AppConfig) Load() {
	log.Println("Loading arguments...")
	kingpin.Parse()
	// Initialize log
	useLogFile := false
	if *logFile != "" {
		useLogFile = true
		log.Printf("Using log file: %s\n", *logFile)
	}
	var f_logFile *os.File = nil
	if useLogFile {
		f_logFile, err := os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		log.SetOutput(f_logFile)
	}
	defer f_logFile.Close()

	// Parse config
	if *confFile == "" {
		if *deviceId == "" {
			panic("Device ID is not defined")
		}
		if *keyFile == "" {
			panic("Key file is not defined")
		}
		appConfig.GetConfigFromFirebase(*firebaseApp, *deviceId, *keyFile)
	} else {
		log.Printf("Using config file: %s\n", *confFile)
		appConfig.ParseConfigFile(*confFile)
	}
}

func (appConfig *AppConfig) GetConfigFromFirebase(firebaseApp string, deviceId string, keyFile string) {
	client := firebase_helper.NewFirebaseClient(keyFile)
	if *skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}
	body := firebase_helper.GetData(fmt.Sprintf("https://%s.firebaseio.com/config/%s", firebaseApp, deviceId), client)

	if fmt.Sprintf("%s", body) == "null" {
		log.Fatal("Cannot get online config for device %s", deviceId)
	} else {
		appConfig.ParseConfig([]byte(body))
	}
}

func (appConfig *AppConfig) ParseConfigFile(confPath string) {
	// Parse config file
	raw, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	appConfig.ParseConfig(raw)
}

func (appConfig *AppConfig) ParseConfig(raw []byte) {
	err := json.Unmarshal(raw, appConfig)
	if err != nil {
		log.Print("Error:", err)
		os.Exit(1)
	}
}
