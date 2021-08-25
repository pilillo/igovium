package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Args ... Arguments provided either as env vars or string args
var Args struct {
	Config string `required:"true" arg:"-c,required"`
}

// Config ... Defines a model for the input config files
type Config struct {
	RESTConfig    *RESTConfig    `yaml:"rest,omitempty"`
	GRPCConfig    *GRPCConfig    `yaml:"grpc,omitempty"`
	DMCacheConfig *DMCacheConfig `yaml:"dm-cache,omitempty"`
	DBCacheConfig *DBCacheConfig `yaml:"db-cache,omitempty"`
}

type RESTConfig struct {
	Port int `yaml:"port,omitempty"`
}

type GRPCConfig struct {
	Port int `yaml:"port,omitempty"`
}

type DMCacheConfig struct {
	Type        string `yaml:"type,omitempty"`
	Mode        string `yaml:"mode,omitempty"`
	HostAddress string `yaml:"host-address,omitempty"`
	Password    string `yaml:"password,omitempty"`
}

type DBCacheConfig struct {
	DriverName               string             `yaml:"driver-name,omitempty"`
	DataSourceName           string             `yaml:"data-source-name,omitempty"`
	MaxLocalCacheElementSize int                `yaml:"local-cache-size,omitempty"`
	Historicize              *HistoricizeConfig `yaml:"historicize,omitempty"`
}

type HistoricizeConfig struct {
	// Cron schedule to trigger the historization process for
	Schedule string `yaml:"schedule,omitempty"`
	// Endpoint of target distributed volume (e.g. s3)
	Endpoint string `yaml:"endpoint,omitempty"`
	// whether to use a secure connection (https)
	UseSSL bool `yaml:"use-ssl,omitempty"`
	// Target Bucket
	Bucket string `yaml:"bucket,omitempty"`
	// Output file format
	Format string `yaml:"format,omitempty"`
	// Path to temporary output file
	TmpDir string `yaml:"tmp-dir,omitempty"`
}

func LoadCfg() *Config {
	err := envconfig.Process("igovium", &Args)
	if err != nil {
		log.Printf("Impossible to parse from env vars - %v", err.Error())
		log.Printf("Attempting parsing string arguments")
		arg.MustParse(&Args)
	}
	// load config from file
	return Load(Args.Config)
}

// Load ... load configuration from file path
func Load(filename string) *Config {
	if !fileExists(filename) {
		log.Fatalf("Configuration file %s does not exist (or is a directory)", filename)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	config, err := parseCfg(data)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}
	config, err = validateCfg(config)
	if err != nil {
		log.Fatalf("Error - %v", err)
	}

	return config
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseCfg(data []byte) (*Config, error) {
	cfg := &Config{}

	err := yaml.Unmarshal(data, &cfg)
	log.Println("Successfully loaded config")

	return cfg, err
}

func validateCfg(cfg *Config) (*Config, error) {
	// todo add validation of input config
	return cfg, nil
}
