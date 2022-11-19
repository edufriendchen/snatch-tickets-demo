# Snatch-Tickets-Demo

#### 解：字节校园镜像技术项目实战活动 —— [【后端】如果有一千万个人抢票怎么办？](https://bytedancecampus1.feishu.cn/docx/doxcnmevBUDWUE9egG6pSBaoVKf)  



由于本人服务器配置和数量资源有限，只能通过单机下的并发性能进一步合理推演集群下的并发性能。

- 框架  
-放票服务ticketing_server使用的字节Golang 微服务 HTTP 框架——[Hertz](https://www.cloudwego.io/zh/docs/hertz/)搭建。  
-客户服务client同样使用Hertz搭建，并参考[官方hertz压测示例项目](https://github.com/cloudwego/hertz-benchmark)实现了模拟高并发请求。    
-订单信息消费服务order_server使用了MQ消息队列，实现订单与放票服务器的解耦也起到削峰的作用。  
-利用go的协程以及Hertz的netpoll对网络的优化来尽可能的提升单机下的并发。

- 系统高并发  
首先我们清楚单机承受的并发是很有限的，并且单机下系统的可用性不能得到保证。所以高并发的系统都是靠服务器集群共同承受的，通过负载均衡服务器将用户请求量均衡到服务器集群上，这样单机所承受的并发量就小了很多，就算某一台机器宕机也不会直接导致系统瘫痪，系统的可用性得到了保证。  
- 单机高并发  
下一步是优化单机下的并发性能，首先我将系统中票的总库存量分配到本地服务器的内存中，直接在内存中减库存，然后将用户抢票成功的信息发送到MQ消息队列中，在另一个系统中异步创建订单并持久化，这样就避免了放票服务器对数据库频繁的IO操作抢占CPU资源影响并发性能。  
- 防止超卖少卖  
然后就是防止超卖和少卖的问题，防止超卖我的思路是这样的，将票的总库存量放在Redis集群中，用户请求在服务器本地内存扣除票库存成功后再去Redis上扣除一下票库存，只有都扣库存成功才算抢票成功，Redis上票的总库存是一定的，所以票就不会超卖。Redis单机的写并发再8w左右、读在10w左右，Redis集群并发量更强大，可以根据系统的并发数搭建Redis集群,还有就是从Redis上扣除库存其实有一个坑，Redis的事务不具备原子性，我们可以将批量操作Redis的指令写在lua脚本中。最后是防止少买，要想避免少卖首先我们得知道发生少卖的原因，当系统中某一台服务器发生宕机时，分配到这台服务器内存上的票库存就不会被卖出，于是少卖就发生了。最简单的解决办法就是给每台服务器内存上多分配一些票库存。例如：票总库存1000张，有十台服务器。可以给每台服务器内存上分配120张票，这样如果1台服务器宕机，剩下9台服务器上的票库存共1080张，就不会发生少买。如果2台服务器宕机了怎么办，这只能看你对服务器的信任值了，信任值低或者服务器宕机数量多就只能多给服务器多分配一些票库存了。还有就是Redis宕机也会造成少卖，在网上读了一些博客，这种情况下应该得让运维人员手动停Redis然后恢复Redis服务，用脚本将票库存等热点信息写入Redis。





下面就是单机压测的过程了，其实准备用Kitex微服务、Ncos服务发现来做集群，4核8G服务器我都没有，集群没搭起来，只能单机压测然后推演集群并发性能了--  

压测相关内容参考了字节[官方hertz压测示例项目](https://github.com/cloudwego/hertz-benchmark)


#### 参考视频

https://bytedancecampus1.feishu.cn/minutes/obcn2m3bdq645dm627u6u38w



#### 服务器配置

- CPU: 2核(vCPU)

- Memory:  2G

- CentOS 8.5 64位

- Go: 1.17.13

- Python：3.6

  

#### 项目启动



订单信息消费服务

```
go run ./server/order_server.go
```

放票服务

```
go run ./server/ticketing_server.go
```

模拟客户端的服务

```
./output/bin/client -addr="http://127.0.0.1:8001/ping" -b=1 -c=100 -n=40000 -s=ticketing_server
```


#### 手动压测

ab压测：

```
ab -k -n 50000 -n 500 http://127.0.0.1:8001/ping
```

wrk压测：

```
wrk -t10 -c30 -d 2s -T5s --latency http://127.0.0.1:8001/ping
```



#### 压测脚本执行  

****

（注：脚本自行拉起服务测试，不需要手动启动服务，请确保服务所需端口不被占用）

由于默认压测参数会比较迅速完成一次压测，为了获取最大程度的压测信息，可以手动调整个压测脚本中压测参数 n 大小。

利用脚本启动自定义的客户端client压测：

```
./scripts/benchmark.sh
```

ab压测脚本:

```
./scripts/benchmark_ab.sh
```

wrk压测脚本： (此脚本会在项目目录里生成性能图像)

```
./scripts/benchmark_wrk.sh
```




#### 单机性能报告

本项目单机下性能参考：

测试环境[ ](https://www.cloudwego.io/zh/docs/kitex/overview/#测试环境)

- CPU: 2核
- Memory: 2GB
- CentOS: 8.5 64位
- Go: 1.17.13

![2022-10-17-20-56_qps](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-3e606520-77f6-4d4b-9877-c12b9367d54c/7456ab43-68f7-4b22-b261-3c7fa8dc26fd.png)



Hertz官方性能参考：

测试环境[ ](https://www.cloudwego.io/zh/docs/kitex/overview/#测试环境)

- CPU: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz, 4 cores
- Memory: 8GB
- OS: Debian 5.4.56.bsk.1-amd64 x86_64 GNU/Linux
- Go: 1.15.4

![performance](https://vkceyugu.cdn.bspapp.com/VKCEYUGU-3e606520-77f6-4d4b-9877-c12b9367d54c/3f8d58fa-dee7-4ca9-8b09-1760accb33c8.png)


#### 单机完成度

| 同时抢票人数 | 1e5  | 1e6  | 5e6  | 1e7  | 5e7  |
| ------------ | ---- | ---- | ---- | ---- | ---- |
| 分数         | 60   | 70   | 80   | 90   | 95   |
| 完成         | √    | √    | ⍻    | ×    | ×    |


单机字节镜像的老师给打了72分，如果上集群可能会更好吧

![image](https://user-images.githubusercontent.com/78396698/198861042-1ce62941-4841-4e41-898a-e6ac3368b2b1.png)


