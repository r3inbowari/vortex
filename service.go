package vortex

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
)

type Order struct {
	Interval  *int64 `json:"interval"`  // 自动测量间隔时间(秒)
	Auto      *bool  `json:"auto"`      // 自动测量
	Operation []byte `json:"operation"` // 操作指令
	Name      string `json:"name"`      // 指令名称
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

	OrderSet []Order `json:"order"` // 传感器指令集
}

/**
 * 本地服务
 */
type LocalService struct {
	Name           string  `json:"name"`             // 收集器名称
	LoggerLevel    string  `json:"loggerLevel"`      // 日志等级
	VortexPort     *int    `json:"vortexPort"`       // 服务端口
	BrokerIP       *string `json:"broker_ip"`        // 中间件地址
	BrokerPort     *string `json:"broker_port"`      // 中间件端口
	BrokerScheme   *string `json:"broker_scheme"`    // 中间件协议
	BrokerUsername *string `json:"broker_username"`  // 中间件用户名
	BrokerPassword *string `json:"broker_password"`  // 中间件密码
	BrokerClientID *string `json:"broker_client_id"` // ClientID

	SensorSet []*SensorInfo `json:"sensorInformation"` // 传感器集合
}

/**
 * Load cnf/conf.json
 */
func GetConfig() *LocalService {
	config := new(LocalService)
	if err := LoadConfig("cnf/conf.json", config); err != nil {
		Fatal("Loading File Failed")
		return nil
	}
	return config
}

/**
 * Save cnf/conf.json
 */
func (ls *LocalService) SetConfig() error {
	fp, err := os.Create("cnf/conf.json")
	if err != nil {
		Fatal("Loading File Failed", logrus.Fields{"err": err})
	}
	defer fp.Close()
	data, err := json.Marshal(ls)
	if err != nil {
		Fatal("Marshal File Failed", logrus.Fields{"err": err})
	}
	n, err := fp.Write(data)
	if err != nil {
		Fatal("Write File Failed", logrus.Fields{"err": err})
	}
	Info("已更新CONFIG文件", logrus.Fields{"size": n})
	return nil
}
