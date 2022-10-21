package util

import "github.com/garyburd/redigo/redis"


//初始化redis连接池
func NewPool() *redis.Pool {

	return &redis.Pool{
		MaxIdle:   10000,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":7963", redis.DialPassword("friendchen7963"))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
