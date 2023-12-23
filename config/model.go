package config

import "time"

type Root struct {
	API         API          `mapstructure:"api"`
	JWT         JWT          `mapstructure:"jwt"`
	Mongo       MongoConfig  `mapstructure:"mongo"`
	Redis       RedisConfig  `mapstructure:"redis"`
	RedisOption RedisOptions `mapstructure:"redis_options"`
	RSAOption   RSAOption    `mapstructure:"rsa_option"`
	File        File         `mapstructure:"file"`
}

type API struct {
	RunningPort    string `mapstructure:"running_port"`
	MaxRequestSize int64  `mapstructure:"maximum_request_size"`
}

type JWT struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
	Issure     string        `mapstructure:"issure"`
}

type MongoConfig struct {
	HostPort   string `mapstructure:"host_port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	DbName     string `mapstructure:"db_name"`
	AuthSource string `mapstructure:"auth_source"`
}

type RedisConfig struct {
	HostPort string `mapstructure:"host_port"`
	Index    int    `mapstructure:"index"`
	Password string `mapstructure:"passowrd"`
}

type RedisOptions struct {
	RSA   RedisOption `mapstructure:"rsa"`
	Login RedisOption `mapstructure:"login"`
}

type RedisOption struct {
	KeyFormat string        `mapstructure:"key_format"`
	TTL       time.Duration `mapstructure:"ttl"`
}

type RSAOption struct {
	RandomKeySize int `mapstructure:"random_key_size"`
	RSASize       int `mapstructure:"rsa_size"`
}

type File struct {
	SpiderImage   string `mapstructure:"spider_image"`
	FileImagePath string `mapstructure:"file_image_path"`
}
