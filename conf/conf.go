package conf

import (
	"encoding/json"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"influxcluster/logging"
	"os"
	"time"
)

var configPath string
var logger = logging.GetLogger("conf")
var cli *clientv3.Client

type BackendConfig struct{
	URL             string
	Zone            string
	Interval        int
	Timeout         int
	TimeoutQuery    int
	MaxRowLimit     int
	CheckInterval   int
}

type NodeConfig struct{
	ListenAddr   string
	DB           string
	Zone         string
	Nexts        string //the backends keys, will accept all data, split with ','
	Interval     int
	IdleTimeout  int
	WriteTracing int
	QueryTracing int
}

var appConfig NodeConfig

func init() {
	configPath = "conf_" + os.Getenv("APP_ENV")
	loadConfig(configPath)
	cfg := clientv3.Config{
		Endpoints:               viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: time.Second*time.Duration(viper.GetInt("etcd.header_timeout")),
	}

	var err error
	cli, err = clientv3.New(cfg)
	if err != nil {
		logger.Fatal(cfg,err)
	}

	logger.Info(cli.Endpoints())
	initConfig()
}

func initConfig() {
	resp, err := cli.Get(context.Background(), configPath)
	if err != nil{
		logger.Error(resp,err)
	}else if resp != nil && len(resp.Kvs) >= 1{
		err = json.Unmarshal(resp.Kvs[0].Value, &appConfig)
		if err != nil {
			logger.Error(resp,err)
		}
	}
	watchAndUpdate()
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
	if err != nil {             // Handle errors reading the config file
		logger.Fatal("fail to load config file:", err)
	}
}

func watchAndUpdate() {
	rch := cli.Watch(context.Background(), configPath, clientv3.WithProgressNotify())
	go func() {
		// watch 该节点下的每次变化
		for wresp := range rch {
			for _, ev := range wresp.Events {
				err := json.Unmarshal(ev.Kv.Value, &appConfig)
				if err != nil {
					//logger.Info(resp,err)
					continue
				}
			}
		}
	}()
}

func GetConfig() NodeConfig{
	return appConfig
}
