package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

type PoolInfo struct {
	Host       string
	Port       uint16
	SubAccount string
}

func (r *PoolInfo) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if len(tmp) > 0 {
		if err := json.Unmarshal(tmp[0], &r.Host); err != nil {
			return err
		}
	}
	if len(tmp) > 1 {
		if err := json.Unmarshal(tmp[1], &r.Port); err != nil {
			return err
		}
	}
	if len(tmp) > 2 {
		if err := json.Unmarshal(tmp[2], &r.SubAccount); err != nil {
			return err
		}
	}
	return nil
}

func (r *PoolInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{r.Host, r.Port, r.SubAccount})
}

type Config struct {
	MultiUserMode               bool       `json:"multi_user_mode"`
	AgentType                   string     `json:"agent_type"`
	AlwaysKeepDownconn          bool       `json:"always_keep_downconn"`
	DisconnectWhenLostAsicboost bool       `json:"disconnect_when_lost_asicboost"`
	UseIpAsWorkerName           bool       `json:"use_ip_as_worker_name"`
	IpWorkerNameFormat          string     `json:"ip_worker_name_format"`
	SubmitResponseFromServer    bool       `json:"submit_response_from_server"`
	AgentListenIp               string     `json:"agent_listen_ip"`
	AgentListenPort             uint16     `json:"agent_listen_port"`
	PoolUseTls                  bool       `json:"pool_use_tls"`
	UseIocp                     bool       `json:"use_iocp"`
	FixedWorkerName             string     `json:"fixed_worker_name"`
	Pools                       []PoolInfo `json:"pools"`
	HTTPDebug                   struct {
		Enable bool   `json:"enable"`
		Listen string `json:"listen"`
	} `json:"http_debug"`
}

// NewConfig 创建配置对象并设置默认值
func NewConfig() (config *Config) {
	config = new(Config)
	config.DisconnectWhenLostAsicboost = true
	config.IpWorkerNameFormat = DefaultIpWorkerNameFormat
	return
}

// LoadFromFile 从文件载入配置
func (conf *Config) LoadFromFile(file string) (err error) {
	configJSON, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(configJSON, conf)
	return
}

func (conf *Config) Init() {
	if conf.MultiUserMode {
		glog.Info("[OPTION] Multi user mode: Enabled. Sub-accounts in config file will be ignored.")
	} else {
		glog.Info("[OPTION] Multi user mode: Disabled. Sub-accounts in config file will be used.")
	}

	glog.Info("[OPTION] Connect to pool server with SSL/TLS encryption: ", IsEnabled(conf.PoolUseTls))
	glog.Info("[OPTION] Always keep miner connections even if pool disconnected: ", IsEnabled(conf.AlwaysKeepDownconn))
	glog.Info("[OPTION] Disconnect if a miner lost its AsicBoost mid-way: ", IsEnabled(conf.DisconnectWhenLostAsicboost))

	if len(conf.FixedWorkerName) > 0 {
		glog.Info("[OPTION] Fixed worker name enabled, all worker name will be replaced to ", conf.FixedWorkerName, " on the server.")
	}

	for i := range conf.Pools {
		pool := &conf.Pools[i]
		if conf.MultiUserMode {
			// 如果启用多用户模式，删除矿池设置中的子账户名
			pool.SubAccount = ""
			glog.Info("add pool: ", pool.Host, ":", pool.Port, ", multi user mode")
		} else {
			glog.Info("add pool: ", pool.Host, ":", pool.Port, ", sub-account: ", pool.SubAccount)
		}
	}
}
