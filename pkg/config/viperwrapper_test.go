package config

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/spf13/viper"
)

var sampleConfigW = []byte(`
--- 
apiVersion: v1
server: 
  delay: 3.4
  host: 127.0.0.1
  port: 9090
services: 
  db: 
    connection: "postgres://"
    engine: postgres
`)

func writeConfig(data []byte) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fileName := path.Join(dir, "config.yaml")

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
	return fileName
}

func TestMain(t *testing.M) {

	fileName := writeConfig(sampleConfigW)

	code := t.Run()

	err := os.Remove(fileName)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func Test_initViper(t *testing.T) {

	var settingChanged bool = false
	onChange := func() error {
		settingChanged = true
		return nil
	}

	err := initViper("config", "yaml", "test", onChange)
	if err != nil {
		t.Errorf("init viper failed %v", err)
	}

	serverHost := viper.GetString("server.host")
	if serverHost != "127.0.0.1" {
		t.Errorf("viper settings is not loaded correct")
	}

	var sampleConfigChanged = []byte(`
--- 
apiVersion: v1
server: 
  delay: 3.4
  host1: "0.0.0.0"
  port: 9080
services: 
  db: 
    connection: "postgres://"
    engine: postgres
	`)

	writeConfig(sampleConfigChanged)
	if settingChanged != true {
		t.Error("setting file changed but watch not worked")
	}
}
