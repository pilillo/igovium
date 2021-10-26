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
	Type         string  `yaml:"type,omitempty"`
	Mode         string  `yaml:"mode,omitempty"`
	HostAddress  string  `yaml:"host-address,omitempty"`
	Password     string  `yaml:"password,omitempty"`
	K8sDiscovery *string `yaml:"k8s-discovery,omitempty"`
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
	// Output file format
	Format string `yaml:"format,omitempty"`
	// Path to temporary output file
	TmpDir string `yaml:"tmp-dir,omitempty"`
	// format to be used for organizing the output files in folders
	DatePartitioner string `yaml:"date-partitioner,omitempty"`
	// embedded remote volume configuration
	// yaml behaves slightly differently than json unmarshalling
	// https://github.com/go-yaml/yaml/issues/63
	RemoteVolumeConfig `yaml:",inline"`
}

// RemoteVolumeConfig holds the configuration for remote volumes
type RemoteVolumeConfig struct {
	DeleteLocal bool      `yaml:"delete-local"`
	S3Config    *S3Config `yaml:"s3,omitempty"`
}

type S3Config struct {
	// Endpoint of target distributed volume (e.g. s3)
	Endpoint string `yaml:"endpoint"`
	// whether to use a secure connection (https)
	UseSSL bool `yaml:"use-ssl"`
	// Target Bucket
	Bucket string `yaml:"bucket"`
	// credentials for target bucket
	AccessKeyVarName string `yaml:"access-key-varname"`
	SecretKeyVarName string `yaml:"secret-key-varname"`
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
