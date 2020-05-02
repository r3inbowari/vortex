package vortex

import (
	"github.com/sirupsen/logrus"
	"sync"
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

func (si *SensorInfo) CreateTask(times int, queueChannel chan Order) {
	for i := 0; i < len(si.OrderSet); i++ {
		if *si.OrderSet[i].Auto {
			key := TaskKey{
				id:   si.ID,
				name: si.OrderSet[i].Name,
			}
			data := TaskData{"data": si.OrderSet[i], "channel": queueChannel, "id": si.ID}
			b := time.Duration(*si.OrderSet[i].Interval) * time.Second
			if err := GetTimeWheel().AddTask(b, -1, key, data, TaskPush); err != nil {
				Warn("add task error", logrus.Fields{"id": si.ID, "operation_name": si.OrderSet[i].Name})
			}
		}
	}
}

/**
 * 单DTU任务阻塞队列的压入
 * @param queueChannel 单DTU内任务的阻塞队列
 */
func TaskPush(data TaskData) {
	body := data["data"].(Order)
	queueChannel := data["channel"].(chan Order)
	defer func() {
		if recover() != nil {
			Warn("task channel has been close", logrus.Fields{"id": data["id"].(string), "operation_name": body.Name})
			key := TaskKey{data["id"].(string), body.Name}
			if err := GetTimeWheel().RemoveTask(key); err != nil {
				Warn("remove a non existent key", logrus.Fields{"id": data["id"].(string), "operation_name": body.Name})
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
func (ds *DTUSession) SetTaskExecutor(queueChannel chan Order, handler func(order Order, wg *sync.WaitGroup)) {
	go func() {
		var wg sync.WaitGroup

		for v := range queueChannel {
			wg.Add(1)
			handler(v, &wg)
			wg.Wait()
		}
		Info("queueChannel has been close", logrus.Fields{"addr": GetIP(ds.conn)})
	}()
}
