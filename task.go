package vortex

import (
	"github.com/sirupsen/logrus"
	"time"
)

var tw *TimeWheel

func GetTimeWheel() *TimeWheel {
	if tw == nil {
		tw = New(time.Second, 60).Start()
	}
	return tw
}

func (ds *DTUSession) LoadTask() {
	ss := GetConfig().GetSensorSetByAttach(GetIP(ds.conn))

	for i := 0; i < len(ss); i++ {
		Info("add "+ss[i].ID, logrus.Fields{"attach": GetIP(ds.conn)})
	}
}

func (ds *DTUSession) GetTask(id string) []Order {
	ss := GetConfig().GetSensorSetByAttach(GetIP(ds.conn))

	for i := 0; i < len(ss); i++ {
		if ss[i].ID == id {
			return ss[i].OrderSet
		}
	}
	return nil
}

type TaskKey struct {
	id   string
	name string
}

func (si *SensorInfo) CreateTask(queueChannel chan Order) {
	for i := 0; i < len(si.OrderSet); i++ {
		if *si.OrderSet[i].Auto {
			si.OrderSet[i].ID = si.ID
			key := TaskKey{
				id:   si.ID,
				name: si.OrderSet[i].Name,
			}
			data := TaskData{"data": si.OrderSet[i], "channel": queueChannel}

			b := time.Duration(*si.OrderSet[i].Interval) * time.Second

			if err := GetTimeWheel().AddTask(b, -1, key, data, TaskPush); err != nil {
				Warn("add task error", logrus.Fields{"id": si.ID, "order_name": si.OrderSet[i].Name})
			}
			Info("add task", logrus.Fields{"id": si.ID, "order_name": si.OrderSet[i].Name})
		}
	}
}

func (si *SensorInfo) RemoveTask() {
	for i := 0; i < len(si.OrderSet); i++ {
		key := TaskKey{si.ID, si.OrderSet[i].Name}
		if err := GetTimeWheel().RemoveTask(key); err != nil {
			Warn("remove a non existent key", logrus.Fields{"id": si.ID, "order_name": si.OrderSet[i].Name})
		}
	}
	Info("remove task", logrus.Fields{"id": si.ID})
}

func (or *Order) RemoveTask() {
	key := TaskKey{or.ID, or.Name}
	if err := GetTimeWheel().RemoveTask(key); err != nil {
		Warn("remove a non existent key", logrus.Fields{"id": or.ID, "order_name": or.Name})
	} else {
		Info("remove task", logrus.Fields{"id": or.ID, "order_name": or.Name})
	}
}

/**
 * 任务
 * @param queueChannel
 */
func TaskPush(data TaskData) {
	body := data["data"].(Order)
	queueChannel := data["channel"].(chan Order)
	defer func() {
		if recover() != nil {
			Warn("task channel has been close", logrus.Fields{"id": body.ID, "order_name": body.Name})
			key := TaskKey{body.ID, body.Name}
			if err := GetTimeWheel().RemoveTask(key); err != nil {
				Warn("remove a non existent key", logrus.Fields{"id": body.ID, "order_name": body.Name})
			} else {
				Info("remove task", logrus.Fields{"id": body.ID, "order_name": body.Name})
			}
		}
	}()
	queueChannel <- body
}

/**
 * task executor set
 * @param queueChannel 任务队列
 * @param handler 处理函数
 */
func (ds *DTUSession) setTaskExecutor(queueChannel chan Order, handler func(order Order)) {
	go func() {
		for v := range queueChannel {

			handler(v)
		}
		Info("queueChannel has been close", logrus.Fields{"addr": GetIP(ds.conn)})
	}()
}

func (ds *DTUSession) TaskSetup() chan Order {
	Info("task init", logrus.Fields{"addr": GetIP(ds.conn)})

	ch := make(chan Order, 10)
	// 这个pop每个dtu有且只有一个, 生命周期应与tcp挂钩
	ds.setTaskExecutor(ch, func(order Order) {
		if result, err := ds.WriteAndWaitRead(order.Operation); err != nil {
			order.RemoveTask()
		} else {
			println(result)
		}
	})
	for _, v := range GetConfig().GetSensorSetByAttach(GetIP(ds.conn)) {
		v.CreateTask(ch)
	}
	return ch
}
