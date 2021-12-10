package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var configFileName string

type WebAdminConfig struct {
	Addr        string `yaml:"addr"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	JwtKey      string `yaml:"jwtkey"`
	IdentityKey string `yaml:"identitykey"`
}

type DataBase struct {
	Type        string `mapstructure:"type"`
	Host        string `mapstructure:"host"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	TablePrefix string `mapstructure:"table-prefix"`
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	From     string `yaml:"from"`
	Password string `yaml:"password"`
	To       string `yaml:"to"`
	// Others   string `yaml:"others"`
	// Interval int64  `yaml:"interval"`
}

type Ssh struct {
	Server   []string `yaml:"server"`
	Fw       string   `yaml:"fw"`
	Ipset    string   `yaml:"ipset"`
	Host     string   `yaml:"host"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type Es struct {
	Host     []string `yaml:"host"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	GwIndex  string   `yaml:"gwindex"`
	WsIndex  string   `yaml:"wsindex"`
}

type Config struct {
	Mode     string          `yaml:"mode"`
	WebAdmin *WebAdminConfig `yaml:"web"`
	Database *DataBase `mapstructer:"database"`
	Secret   string `yaml:"secret"`
	LogPath  string `yaml:"log_path"`
	LogLevel string `yaml:"log_level"`

	Ssh []Ssh `yaml:"ssh"`

	DbPath string `yaml:"db_path"`

	Redis *Redis `yaml:"redis"`

	Es *Es `yaml:"es"`

	Email *EmailConfig `yaml:"email"`
}

var cfg Config

func ParseConfigData(data []byte) (*Config, error) {
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseConfigFile(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	configFileName = fileName

	return ParseConfigData(data)
}

func GetEncryptKey() (string, error) {
	if len(cfg.Secret) == 0 {
		return "", fmt.Errorf("%s", "config secret is empty")
	}

	return cfg.Secret, nil
}
