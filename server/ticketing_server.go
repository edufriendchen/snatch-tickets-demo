/*
 *  FriendChen Authors
 *  2022-10-15
 */

package main

import (
	"github.com/cloudwego/hertz-benchmark/handler"
	"github.com/cloudwego/hertz-benchmark/util"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
)

const (
	port = ":8001"
)


func main() {
	//初始化rabbitMq
	util.Mq = util.NewRabbitMQ("ticketQueue", "ticketMq", "routingKey")
	defer util.Mq.ReleaseRes()
	opts := []config.Option{
		server.WithHostPorts(port),
	}
	h := server.New(opts...)
	h.GET("/ping", handler.EchoHandler)
	h.Spin()
}
