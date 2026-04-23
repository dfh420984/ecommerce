package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Wechat   WechatConfig   `yaml:"wechat"`
	Alipay   AlipayConfig   `yaml:"alipay"`
	Cors     CorsConfig     `yaml:"cors"`
}

type AppConfig struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Mode       string `yaml:"mode"`
	UploadPath string `yaml:"upload_path"`
	ServerURL  string `yaml:"server_url"`
	JWTSecret  string `yaml:"jwt_secret"`
	JWTExpire  int    `yaml:"jwt_expire"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type WechatConfig struct {
	AppID     string `yaml:"appid"`
	Secret    string `yaml:"secret"`
	MchID     string `yaml:"mchid"`
	APIKey    string `yaml:"apikey"`
	CertPath  string `yaml:"cert_path"`
	KeyPath   string `yaml:"key_path"`
	NotifyURL string `yaml:"notify_url"`
}

type AlipayConfig struct {
	AppID      string `yaml:"appid"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
	NotifyURL  string `yaml:"notify_url"`
}

type CorsConfig struct {
	AllowOrigins []string `yaml:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods"`
	AllowHeaders []string `yaml:"allow_headers"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		var data []byte
		data, err = os.ReadFile(path)
		if err != nil {
			return
		}
		cfg = &Config{}
		err = yaml.Unmarshal(data, cfg)
	})
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func Get() *Config {
	return cfg
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
