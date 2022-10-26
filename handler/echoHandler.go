package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz-benchmark/util"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)

var (
	done chan int
	num  = 1
)

func EchoHandler(c context.Context, ctx *app.RequestContext) {
	fmt.Print("-", num)
	num++
	//全局读写锁
	if util.LocalStock.LocalDeductionStock() && util.CloudStock.RemoteDeductionStock(util.RedisConn) {
		//生成订单号  = 本地服务器编号 + 时间戳 + 该类型票编号 + 本地库存量
		order := util.LocalStock.LocalServerNumber + fmt.Sprintf("%d", time.Now().Unix()) + util.CloudStock.SpikeOrderHashKey + strconv.FormatInt(util.LocalStock.LocalSalesVolume, 10)
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
}

// 发送订单
func sendOrderInfo(order string) {
	err := util.Mq.Channel.Publish(
		util.Mq.Exchange,   // 交换器名
		util.Mq.RoutingKey, // routing key
		false,              // 是否返回消息(匹配队列)，如果为true, 会根据binding规则匹配queue，如未匹配queue，则把发送的消息返回给发送者
		false,              // 是否返回消息(匹配消费者)，如果为true, 消息发送到queue后发现没有绑定消费者，则把发送的消息返回给发送者
		amqp.Publishing{ // 发送的消息，固定有消息体和一些额外的消息头，包中提供了封装对象
			ContentType: "text/plain",                    // 消息内容的类型
			Body:        []byte("订单信息: [" + order + "]"), // 消息内容
		},
	)
	if err != nil {
		return
	}
}
