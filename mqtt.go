package vortex

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// scheme -> ws/ssl/tcp
// var scheme = "tcp"
// var host = "106.13.79.157"
// var port = "1883"

// ClientID 随机acm0-bjd2-fdi1-am81
// var ClientID = bson.NewObjectId().String()
// var Username = "r3inb"
// var Password = "159463"
// var base = "r3inbowari.top"

var client mqtt.Client = nil

/**
 * 获取MQ连接单例
 */
func GetMQTTInstance() mqtt.Client {
	if client == nil || !client.IsConnectionOpen() {
		var err error
		client, err = processMQTTClient()
		if err != nil {
			Warn("mqtt broker connect failed")
		} else {
			Info("connected to mq", logrus.Fields{"host": *GetConfig().BrokerHost})
		}
	}
	return client
}

/**
 * mqtt
 */
func processMQTTClient() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*GetConfig().BrokerScheme + "://" + *GetConfig().BrokerHost)
	if GetConfig().BrokerClientID != nil {
		opts.SetClientID(*GetConfig().BrokerClientID)
	}
	opts.SetUsername(*GetConfig().BrokerUsername)
	opts.SetPassword(*GetConfig().BrokerPassword)
	// opts.SetKeepAlive(2 * time.Second)
	// 默认消费方式
	//opts.SetDefaultPublishHandler(defaultPublishHandler)
	// ping超时
	//opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}

/**
 * 订阅主题映射
 */
func MQTTMapping(topic string, callback mqtt.MessageHandler) bool {
	if token := GetMQTTInstance().Subscribe(topic, 1, callback); token.Wait() && token.Error() != nil {
		Warn("subscribe failed", logrus.Fields{"topic": topic})
		return false
	}
	Info("subscribed successfully", logrus.Fields{"topic": topic})
	return true
}

/**
 * 发布主题消息
 */
func MQTTPublish(topic string, payload interface{}) {
	token := GetMQTTInstance().Publish(topic, 1, false, payload)
	token.Wait()
}
