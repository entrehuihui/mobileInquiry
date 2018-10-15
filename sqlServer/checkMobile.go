package sqlServer

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

var ChanMobile = make(chan map[string]string, 20)

//MobileInfoPost 将查询到的信息插入MYSQL中
func MobileInfoPost(mobileInfoList []map[string]string) {
	var buf bytes.Buffer
	buf.WriteString("replace into mobileCheckAddr (mobileNumber, mobileAddr, createTime, operator, checkIp, checkStatus) values ")
	for _, mobileInfo := range mobileInfoList {
		now := time.Now().Unix()
		buf.WriteString("(")
		buf.WriteString("'")
		buf.WriteString(mobileInfo["mobile"])
		buf.WriteString("',")
		buf.WriteString("'")
		buf.WriteString(mobileInfo["addr"])
		buf.WriteString("',")
		buf.WriteString("'")
		buf.WriteString(strconv.FormatInt(now, 10))
		buf.WriteString("',")
		buf.WriteString("'")
		buf.WriteString(mobileInfo["operator"])
		buf.WriteString("',")
		buf.WriteString("'")
		buf.WriteString(mobileInfo["checkIp"])
		buf.WriteString("',")
		buf.WriteString(mobileInfo["checkStatus"])
		buf.WriteString("),")
	}
	SQL := buf.String()
	SQL = SQL[0 : len(SQL)-1]
	db.Exec(SQL)
}

//MobileInfoGet 查询数据
func MobileInfoGet(mobileInfoList map[string]string) (int64, error) {
	var mobile int64
	rows := db.QueryRow("select mobileNumber from mobileCheckAddr where mobileNumber = " + mobileInfoList["mobile"])
	err := rows.Scan(&mobile)
	if err != nil {
		fmt.Println("QueryRow error", err)
		return 0, err
	}
	return mobile, err
}

//MobileInfoDel 删除数据
func MobileInfoDel() {

}

//QueuePost 插入队列
func QueuePost() {
	queueNumberList := 0
	var mobileInfoList []map[string]string
	for mobileInfo := range ChanMobile {
		if _, err := MobileInfoGet(mobileInfo); err == nil {
			continue
		}
		mobileInfoList = append(mobileInfoList, mobileInfo)
		queueNumberList++
		if queueNumberList >= 20 {
			queueNumberList = 0
			go MobileInfoPost(mobileInfoList)
			//清空数组
			mobileInfoList = []map[string]string{}
		}
	}
}
