package vortex

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type TimeWheel struct {
	interval       time.Duration
	ticker         *time.Ticker
	slots          []*list.List
	currentPos     int
	slotNum        int
	addTaskChannel chan *task
	stopChannel    chan bool
	taskRecord     *sync.Map
}

type Job func(TaskData)

type TaskData map[interface{}]interface{}

type task struct {
	interval time.Duration
	times    int
	circle   int
	key      interface{}
	job      Job
	taskData TaskData
}

func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:       interval,
		slots:          make([]*list.List, slotNum),
		currentPos:     0,
		slotNum:        slotNum,
		addTaskChannel: make(chan *task),
		stopChannel:    make(chan bool),
		taskRecord:     &sync.Map{},
	}
	tw.init()
	return tw
}

func (tw *TimeWheel) Start() *TimeWheel {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
	return tw
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(task)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) AddTask(interval time.Duration, times int, key interface{}, data TaskData, job Job) error {
	if interval <= 0 || key == nil || job == nil || times < -1 || times == 0 {
		return errors.New("非法的参数")
	}
	_, ok := tw.taskRecord.Load(key)
	if ok {
		return errors.New("重复的Key")
	}
	tw.addTaskChannel <- &task{interval: interval, times: times, key: key, taskData: data, job: job}
	return nil
}

func (tw *TimeWheel) RemoveTask(key interface{}) error {
	if key == nil {
		return nil
	}
	value, ok := tw.taskRecord.Load(key)
	if !ok {
		return errors.New("不存在的Key")
	} else {
		task := value.(*task)
		task.times = 0
		tw.taskRecord.Delete(task.key)
	}
	return nil
}

func (tw *TimeWheel) UpdateTask(key interface{}, interval time.Duration, taskData TaskData) error {
	if key == nil {
		return errors.New("非法的Key")
	}
	value, ok := tw.taskRecord.Load(key)
	if !ok {
		return errors.New("不存在的Key")
	}
	task := value.(*task)
	task.taskData = taskData
	task.interval = interval
	return nil
}

func (tw *TimeWheel) init() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAddRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

func (tw *TimeWheel) addTask(task *task) {
	if task.times == 0 {
		return
	}
	pos, circle := tw.getPositionAndCircle(task.interval)
	task.circle = circle
	tw.slots[pos].PushBack(task)
	tw.taskRecord.Store(task.key, task)
}

func (tw *TimeWheel) scanAddRunTask(l *list.List) {
	if l == nil {
		return
	}
	for item := l.Front(); item != nil; {
		task := item.Value.(*task)
		if task.times == 0 {
			next := item.Next()
			l.Remove(item)
			tw.taskRecord.Delete(task.key)
			item = next
			continue
		}
		if task.circle > 0 {
			task.circle--
			item = item.Next()
			continue
		}
		go task.job(task.taskData)
		next := item.Next()
		l.Remove(item)
		item = next
		if task.times == 1 {
			task.times = 0
			tw.taskRecord.Delete(task.key)
		} else {
			if task.times > 0 {
				task.times--
			}
			tw.addTask(task)
		}
	}
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.slotNum)
	pos = int(tw.currentPos+delaySeconds/intervalSeconds) % tw.slotNum
	return
}
