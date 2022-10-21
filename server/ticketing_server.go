/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz-benchmark/stock"
	"github.com/cloudwego/hertz-benchmark/util"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

const (
	port = ":8001"
)

var (
	localStock stock.LocalStock
	cloudStock stock.CloudStock
	redisPool  *redis.Pool
	redisConn  redis.Conn
	done       chan int
	mq         *util.RabbitMQ
	num        = 1
)

//初始化配置
func init() {
	//初始化本地票库存
	localStock = stock.LocalStock{
		LocalServerNumber: "001", //第一台服务器
		LocalTicketStock:  100,   //初始100张票
		LocalSalesVolume:  0,     //初始售出0张票
	}

	//定义远程票库存存储结构
	cloudStock = stock.CloudStock{
		SpikeOrderHashKey:  "ticket_hash_key",   //票的编号
		TotalInventoryKey:  "ticket_total_nums", //远程总库存量
		QuantityOfOrderKey: "ticket_sold_nums",  //远程总的订单量
	}
	//初始化redis
	redisPool = util.NewPool()
	redisConn = redisPool.Get()
	defer redisPool.Close()
	done = make(chan int, 1)
	done <- 1
}

func main() {
	//初始化rabbitMq
	mq = util.NewRabbitMQ("ticketQueue", "ticketMq", "routingKey")
	defer mq.ReleaseRes()
	opts := []config.Option{
		server.WithHostPorts(port),
	}
	h := server.New(opts...)
	h.GET("/ping", echoHandler)
	h.Spin()
}

func echoHandler(c context.Context, ctx *app.RequestContext) {
	//fmt.Print("-", num)
	num++
	<-done
	//全局读写锁
	if localStock.LocalDeductionStock() && cloudStock.RemoteDeductionStock(redisConn){
		//生成订单号  = 本地服务器编号 + 时间戳 + 该类型票编号 + 本地库存量
		order := localStock.LocalServerNumber + fmt.Sprintf("%d", time.Now().Unix()) + cloudStock.SpikeOrderHashKey + strconv.FormatInt(localStock.LocalSalesVolume, 10)
		ctx.JSON(200, utils.H{
			"message": "抢票成功",
			"orderId": order,
		})
		//将订单号发送到rabbitMq
		sendOrderInfo(order)
	} else {
		ctx.JSON(200, utils.H{
			"message": "已售罄",
		})
	}
	//将抢票状态写入到log中
	done <- 1
}

// 发送订单
func sendOrderInfo(order string) {
	err := mq.Channel.Publish(
		mq.Exchange,   // 交换器名
		mq.RoutingKey, // routing key
		false,         // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false,         // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		amqp.Publishing{ // 发送的消息，固定有消息体和一些额外的消息头，包中提供了封装对象
			ContentType: "text/plain",                    // 消息内容的类型
			Body:        []byte("订单信息: [" + order + "]"), // 消息内容
		},
	)
	if err != nil {
		return
	}
}
