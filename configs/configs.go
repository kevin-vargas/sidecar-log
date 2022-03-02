package configs

import (
	_ "embed"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

const DEV_SCOPE = "dev"

//go:embed configs/dev.yml
var dev []byte

//go:embed configs/prod.yml
var prod []byte

var configFile = map[string][]byte{
	"dev":  dev,
	"prod": prod,
}
var defaultConfig = prod

type cfg struct {
	MQTT struct {
		BROKER string `yaml:"broker"`
		PORT   int    `yaml:"port"`
		CLIENT struct {
			ID       string `yaml:"id"`
			USERNAME string `yaml:"username"`
			PASSWORD string `yaml:"password"`
		} `yaml:"client"`
	} `yaml:"mqtt"`
	K3S struct {
		NAMESPACE string `yaml:"namespace"`
		APP       string `yaml:"app"`
	} `yaml:"k3s"`
	APP string
}

var once sync.Once
var instance *cfg

func IsDev() bool {
	return os.Getenv("SCOPE") == DEV_SCOPE
}

func Get() *cfg {
	once.Do(func() {
		cfg, err := createConfig()
		if err != nil {
			panic(err)
		}
		instance = cfg
	})
	return instance
}

func isEmpty(slice []byte) bool {
	return len(slice) == 0
}
func createConfig() (*cfg, error) {
	var cfg *cfg = &cfg{}
	ymlBytes := configFile[os.Getenv("SCOPE")]
	if isEmpty(ymlBytes) {
		ymlBytes = defaultConfig
	}
	err := yaml.Unmarshal(ymlBytes, cfg)
	if err != nil {
		return nil, err
	}
	return customFields(cfg), nil
}

func customFields(c *cfg) *cfg {
	// app name
	c.APP = os.Getenv("APP")
	return c
}
