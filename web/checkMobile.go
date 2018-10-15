package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"../sqlServer"
	"github.com/garyburd/redigo/redis"
)

//查询号码归属地
func checkMobile(w http.ResponseWriter, r *http.Request) {
	var mobile []string
	var mobileGet []map[string]string
	clientIP := r.RemoteAddr
	ipString := splitIP(clientIP)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &mobile)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	for _, mobileNumber := range mobile {
		if _, err := strconv.Atoi(mobileNumber); err != nil || len(mobileNumber) != 11 {
			continue
		}
		mobileInfo, err := getRedisValue(string([]byte(mobileNumber)[:7]))
		if err != nil {
			continue
		}
		mobileInfo["checkIp"] = ipString
		mobileInfo["mobile"] = mobileNumber
		//传送到MYSQL模板进行保存
		sqlServer.ChanMobile <- mobileInfo
		mobileGet = append(mobileGet, mobileInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	mobileGetString, _ := json.Marshal(mobileGet)
	fmt.Fprintf(w, string(mobileGetString))
}

func getRedisValue(key string) (map[string]string, error) {
	mobileInfo := make(map[string]string)
	connection, err := redis.DialTimeout("tcp", ":6379", 0, 1*time.Second, 1*time.Second)
	if err != nil {
		fmt.Println("redis.DialTimeout failed:", err)
		return nil, err
	}
	defer connection.Close()
	redisInfo, err := redis.Bytes(connection.Do("GET", key))
	if err != nil {
		mobileInfo["addr"] = "未知"
		mobileInfo["operator"] = "0"
		mobileInfo["checkStatus"] = "0"
		return mobileInfo, nil
	}
	err = json.Unmarshal(redisInfo, &mobileInfo)
	if err != nil {
		fmt.Println("getRedisValue json.Unmarshal failed:", err)
		return nil, err
	}
	mobileInfo["checkStatus"] = "1"
	return mobileInfo, err
}

//拼接IP使其成为12位整数
func splitIP(clientIP string) string {
	clientIPArrar := strings.Split(clientIP, ":")
	clientIPArrayNum := strings.Split(clientIPArrar[0], ".")
	ipString := ""
	for _, clientIPNumber := range clientIPArrayNum {
		if len(clientIPNumber) == 1 {
			clientIPNumber = "0" + "0" + clientIPNumber
		} else if len(clientIPNumber) == 2 {
			clientIPNumber = "0" + clientIPNumber
		}
		ipString = ipString + clientIPNumber
	}
	return ipString
}
