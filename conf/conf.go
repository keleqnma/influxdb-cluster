package conf

import (
	"context"
	"encoding/json"
	"influxcluster/logging"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	logger     = logging.GetLogger("conf")
	cli        *clientv3.Client
	appConfig  APPConfig
	configPath string
	mu         sync.RWMutex
)

type StorageConfig struct {
	URL           string
	Zone          string
	Interval      int
	Timeout       int
	TimeoutQuery  int
	MaxRowLimit   int
	CheckInterval int
}

type NodeConfig struct {
	ListenAddr   string
	Zone         string
	Interval     int
	WriteTracing int
	QueryTracing int
}

type APPConfig struct {
	StorageCfgs []StorageConfig
	NodeCfg     NodeConfig
}

func init() {
	mu = sync.RWMutex{}
	configPath = "conf_" + os.Getenv("APP_ENV")
	loadConfig(configPath, "conf")
	cfg := clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: time.Second * time.Duration(viper.GetInt("etcd.header_timeout")),
	}

	var err error
	cli, err = clientv3.New(cfg)
	if err != nil {
		logger.Fatal(cfg, err)
	}
	initConfig(configPath, &appConfig)
}

func initConfig(key string, v interface{}) {
	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		logger.Error(resp, err)
		return
	}

	if resp != nil && len(resp.Kvs) >= 1 {
		mu.Lock()
		err = json.Unmarshal(resp.Kvs[0].Value, v)
		mu.Unlock()
		if err != nil {
			logger.Error(resp, err)
			return
		}
	}
	watchAndUpdate(key, v)
}

func watchAndUpdate(key string, v interface{}) {
	rch := cli.Watch(context.Background(), key, clientv3.WithProgressNotify())
	go func() {
		// watch 该节点下的每次变化
		for wresp := range rch {
			for _, ev := range wresp.Events {
				mu.Lock()
				err := json.Unmarshal(ev.Kv.Value, v)
				mu.Unlock()
				if err != nil {
					continue
				}
			}
		}
	}()
}

func loadConfig(in string, paths ...string) {
	viper.SetConfigName(in)
	path, err := os.Getwd()
	if err != nil {
		logger.Fatal("fail to get current path:", err)
	}
	viper.AddConfigPath(path)
	for _, configPath := range paths {
		viper.AddConfigPath(configPath)
	}
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		logger.Fatal("fail to load config file:", err)
	}
}

func GetConfig() APPConfig {
	mu.RLock()
	defer mu.RUnlock()
	return appConfig
}
