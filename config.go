package vortex

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

/**
 * 单条指令
 */
type Order struct {
	Interval  *int   `json:"interval"`  // 自动测量间隔时间(秒)
	Auto      *bool  `json:"auto"`      // 自动测量
	Operation []byte `json:"operation"` // 操作指令
	Name      string `json:"name"`      // 指令名称
	ID        string `json:"-"`
}

/**
 * 传感器参数 包含自定义任务和状态等信息
 */
type SensorInfo struct {
	// runtime
	Status int `json:"-"` // 传感器状态

	ID     string `json:"id"`     // 传感器ID
	Addr   byte   `json:"addr"`   // 传感器设备地址
	Type   byte   `json:"type"`   // 传感器类型
	Attach string `json:"attach"` // 传感器附着的透传设备IP

	OrderSet []Order `json:"order_set"` // 传感器指令集
}

func (lc *LocalConfig) GetSensorSetByAttach(ip string) []SensorInfo {
	var retSet = make([]SensorInfo, 0)
	for i := 0; i < len(lc.SensorSet); i++ {
		if lc.SensorSet[i].Attach == ip {
			retSet = append(retSet, *lc.SensorSet[i])
		}
	}
	return retSet
}

/**
 * local config struct
 */
type LocalConfig struct {
	Name           string  `json:"name"`             // 收集器名称
	LoggerLevel    *string `json:"log_level"`        // 日志等级
	VortexPort     *int    `json:"vortex_port"`      // 服务端口
	BrokerHost     *string `json:"broker_host"`      // 中间件地址
	BrokerScheme   *string `json:"broker_scheme"`    // 中间件协议
	BrokerUsername *string `json:"broker_username"`  // 中间件用户名
	BrokerPassword *string `json:"broker_password"`  // 中间件密码
	BrokerClientID *string `json:"broker_client_id"` // ClientID

	SensorSet []*SensorInfo `json:"sensor_set"` // 传感器集合

	CacheDeadline time.Time `json:"-"`
}

var config = new(LocalConfig)

/**
 * load cnf/conf.json
 */
func GetConfig() *LocalConfig {
	if config.CacheDeadline.Before(time.Now()) {
		if err := LoadConfig("cnf/conf.json", config); err != nil {
			Fatal("loading file failed")
			return nil
		}
		config.CacheDeadline = time.Now().Add(time.Second * 60)
	}
	return config
}

/**
 * save cnf/conf.json
 */
func (lc *LocalConfig) SetConfig() error {
	fp, err := os.Create("cnf/conf.json")
	if err != nil {
		Fatal("loading file failed", logrus.Fields{"err": err})
	}
	defer fp.Close()
	data, err := json.Marshal(lc)
	if err != nil {
		Fatal("marshal file failed", logrus.Fields{"err": err})
	}
	n, err := fp.Write(data)
	if err != nil {
		Fatal("write file failed", logrus.Fields{"err": err})
	}
	Info("already update config file", logrus.Fields{"size": n})
	return nil
}
