package util

import (
	"github.com/cloudwego/hertz-benchmark/stock"
	"github.com/garyburd/redigo/redis"
)

var (
	LocalStock stock.LocalStock
	CloudStock stock.CloudStock
	RedisPool  *redis.Pool
	RedisConn  redis.Conn
	done       chan int
	num        = 1
)

//初始化配置
func init() {
	//初始化本地票库存
	LocalStock = stock.LocalStock{
		LocalServerNumber: "001", //第一台服务器
		LocalTicketStock:  100,   //初始100张票
		LocalSalesVolume:  0,     //初始售出0张票
	}

	//定义远程票库存存储结构
	CloudStock = stock.CloudStock{
		SpikeOrderHashKey:  "ticket_hash_key",   //票的编号
		TotalInventoryKey:  "ticket_total_nums", //远程总库存量
		QuantityOfOrderKey: "ticket_sold_nums",  //远程总的订单量
	}
	//初始化redis
	RedisPool = NewPool()
	RedisConn = RedisPool.Get()
	defer RedisPool.Close()
}


// NewPool 初始化redis连接池
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   10000,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379", redis.DialPassword("friendchen"))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
