package vortex

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"sync"
)

const Network = "tcp"

var listener net.Listener

func RunDTUConnector() {
	var err error
	port := GetConfig().VortexPort
	if port == nil {
		port = new(int)
		*port = 8000
	}
	for PortInUse(*port) {
		*port++
	}
	if listener, err = net.Listen(Network, ":"+strconv.Itoa(*port)); err != nil {
		Fatal("Listen failed", logrus.Fields{"port": *port, "err": err})
	} else {
		Info("DTU listen succeed", logrus.Fields{"port": *port})
	}
	defer func() { _ = listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			Warn("Accept failed", logrus.Fields{"err": err})
			break
		}
		go HandleProcessor(conn)
	}
}

type DTUSession struct {
	readChan  chan []byte
	writeChan chan []byte
	stopChan  chan bool
	conn      net.Conn
}

var SessionsCollection sync.Map

func RegDTUSession(conn net.Conn) DTUSession {
	var ds DTUSession
	ds.readChan = make(chan []byte)  // 读
	ds.writeChan = make(chan []byte) // 写
	ds.stopChan = make(chan bool)    // 停
	ds.conn = conn                   // 连接
	addr := GetIP(conn)
	SessionsCollection.Store(addr, ds)
	Info("connected", logrus.Fields{"addr": addr})
	return ds
}

/**
 * 释放一个session
 */
func (ds *DTUSession) ReleaseDTUSession() {
	SessionsCollection.Delete(GetIP(ds.conn))
}

func GetDTUSessionsKey() []string {
	var ret []string
	SessionsCollection.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(string))
		return true
	})
	return ret
}

func HandleProcessor(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	session := RegDTUSession(conn)
	go session.readConn()
	// go session.writeConn()

	for {
		select {
		case stop := <-session.stopChan:
			if stop {
				Info("disconnected", logrus.Fields{"addr": GetIP(conn)})
				session.ReleaseDTUSession()
				return
			}
		case read := <-session.readChan:
			fmt.Println(string(read))
		}
	}
}

func (ds *DTUSession) readConn() {
	for {
		data := make([]byte, 64)
		n, err := ds.conn.Read(data)
		if err != nil {
			break
		}
		ds.readChan <- data[:n]
	}
	ds.stopChan <- true
}

/**
 * 写
 */
func (ds *DTUSession) Write(b []byte) error {
	if _, err := ds.conn.Write(b); err != nil {
		return err
	}
	return nil
}

/**
 * write
 */
func (ds *DTUSession) writeConn() {
	for {
		data := <-ds.writeChan
		if _, err := ds.conn.Write(data); err != nil {
			break
		}
	}
	ds.stopChan <- true
}
