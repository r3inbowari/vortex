package vortex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const configFileSizeLimit = 10 << 20

/**
 * Load File
 * @param path 文件路径
 * @param dist 存放目标
 */
func LoadConfig(path string, dist interface{}) error {
	configFile, err := os.Open(path)
	if err != nil {
		Fatal("Failed to open config file.", logrus.Fields{"path": path, "err": err})
		return err
	}

	fi, _ := configFile.Stat()
	if size := fi.Size(); size > (configFileSizeLimit) {
		Fatal("Config file size exceeds reasonable limited", logrus.Fields{"path": path, "size": size})
		return errors.New("limited")
	}

	if fi.Size() == 0 {
		Fatal("Config file is empty, skipping", logrus.Fields{"path": path, "size": 0})
		return errors.New("empty")
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)
	buffer, err = StripComments(buffer)
	if err != nil {
		Fatal("Failed to strip comments from json", logrus.Fields{"err": err})
		return err
	}

	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, &dist)
	if err != nil {
		Fatal("Failed unmarshalling json", logrus.Fields{"err": err})
		return err
	}
	return nil
}

// StripComments
func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0)
	lines := bytes.Split(data, []byte("\n"))
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}
	return bytes.Join(filtered, []byte("\n")), nil
}

func GetIP(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

func PortInUse(port int) bool {
	checkStatement := fmt.Sprintf("lsof -i:%d ", port)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if len(output) > 0 {
		return true
	}
	return false
}

/**
 * get map key
 */
func GetKeys(m map[int]int) []int {
	i := 0
	keys := make([]int, len(m))
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
