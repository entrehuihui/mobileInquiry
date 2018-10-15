package web

import (
	"log"
	"net/http"
	"time"

	"github.com/drone/routes"
)

//CreateHttpServer 绑定httpServer端口和服务
func CreateHttpServer() {
	mux := routes.New()

	//查询电话号码归属地
	mux.Post("/api/checkMobile/", checkMobile)
	http.Handle("/", mux)
	if err := http.ListenAndServe(":1235", nil); err != nil {
		log.Fatal(time.Now().String()+"ListenAndServe bind 1235 port error! error:", err)
	}
}
