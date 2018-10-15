package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"./sqlServer"
	"./web"
	"github.com/garyburd/redigo/redis"
)

var dirfile *string

func main() {
	//选项：是否重新载入归属地配置表
	renew := flag.Bool("r", false, "yes：重新载入数据库数据，默认不重新载入")
	dirfile = flag.String("d", "./mobile.txt", "dirload：载入归属地配置文件目录")
	//刷新缓存
	flag.Parse()
	if *renew {
		chanQueue()
	}
	go sqlServer.QueuePost()
	web.CreateHttpServer()
}

var chanAddrInfo = make(chan string, 1)

func resolveAddrInfo() {
	file, err := os.Open(*dirfile)
	if err != nil {
		close(chanAddrInfo)
		return
	}
	defer file.Close()
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		chanAddrInfo <- string(a)
	}
	close(chanAddrInfo)
}

func chanQueue() {
	conn, err := redis.DialTimeout("tcp", ":6379", 0, 1*time.Second, 1*time.Second)
	if err != nil {
		fmt.Println("redis.DialTimeout failed:", err)
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup
	chanQuequeList := make(chan int, 1)
	//开携程
	go resolveAddrInfo()
	//等待读取文件
	for info := range chanAddrInfo {
		wg.Add(1)
		chanQuequeList <- 1
		func(info string) {
			setInfo(info, conn)
			<-chanQuequeList
			defer wg.Done()
		}(info)
	}
	wg.Wait()
}

func setInfo(info string, conn redis.Conn) {

	infoArray := strings.Split(info, ",")
	if len(infoArray) != 7 {
		return
	}

	var infoMap = make(map[string]string)
	infoMap["addr"] = infoArray[2] + infoArray[3]
	infoMap["mobile"] = infoArray[1]
	infoMap["operator"] = "0"
	if infoArray[4] == "中国移动" {
		infoMap["operator"] = "1"
	} else if infoArray[4] == "中国联通" {
		infoMap["operator"] = "2"
	} else if infoArray[4] == "中国电信" {
		infoMap["operator"] = "3"
	}
	mjson, _ := json.Marshal(infoMap)
	_, err := conn.Do("SET", infoArray[1], string(mjson))
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
}
