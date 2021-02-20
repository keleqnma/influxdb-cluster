package server

import (
	"flag"
	"log"
)

var (
	ConfigFile  string
	NodeName    string
	RedisAddr   string
	RedisPwd    string
	RedisDb     int
	LogFilePath string
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	flag.StringVar(&LogFilePath, "log-file-path", "/var/log/influx-proxy.log", "output file")
	flag.StringVar(&ConfigFile, "config", "", "config file")
	flag.StringVar(&NodeName, "node", "l1", "node name")
	flag.StringVar(&RedisAddr, "redis", "localhost:6379", "config file")
	flag.StringVar(&RedisPwd, "redis-pwd", "", "config file")
	flag.IntVar(&RedisDb, "redis-db", 0, "config file")
	flag.Parse()
}
