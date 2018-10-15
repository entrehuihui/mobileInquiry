package sqlServer

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/GO-SQL-Driver/MYSQL"
)

var (
	mu sync.RWMutex
	db = &sql.DB{}
)

//MyInit 连接数据库
func init() {
	//打开数据库连接
	var err error
	db, err = sql.Open("mysql", "账号密码:3306)/entredb?charset=utf8")
	if err != nil {
		log.Fatal(err.Error(), err)
	}
}
