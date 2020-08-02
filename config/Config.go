package config

import (
	"time"
)

var Cfg *Config

type Config struct {
	DiscoveryDelay time.Duration
}

func Init() {
	Cfg = &Config{DiscoveryDelay: time.Second * 5}
}
