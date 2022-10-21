package stock

import "github.com/garyburd/redigo/redis"

// LuaScript lua脚本保证redis操作的原子性
const LuaScript = `
       local ticket_key = KEYS[1]
       local ticket_total_key = ARGV[1]
       local ticket_sold_key = ARGV[2]
       local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
       local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
		-- 查看是否还有余票,增加订单数量,返回结果值
       if(ticket_total_nums > ticket_sold_nums) then
           return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1)
       end
       return 0
`

// CloudStock 远程票库存存储结构 （hash结构）
type CloudStock struct {
	SpikeOrderHashKey  string //redis中秒杀订单hash结构key
	TotalInventoryKey  string //中总订单库存key
	QuantityOfOrderKey string //hash结构中已有订单数量key
}

// RemoteDeductionStock 远端统一扣库存
func (RemoteSpikeKeys *CloudStock) RemoteDeductionStock(conn redis.Conn) bool {
	lua := redis.NewScript(1, LuaScript)
	result, err := redis.Int(lua.Do(conn, RemoteSpikeKeys.SpikeOrderHashKey, RemoteSpikeKeys.TotalInventoryKey, RemoteSpikeKeys.QuantityOfOrderKey))
	if err != nil {
		return false
	}
	return result != 0
}
