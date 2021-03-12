package runtimes

import (
	"math/rand"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

var POOL *redis.Pool

/**
* 包的执行函数
 */
func init() {
	//使用redis连接池
	redisHost := beego.AppConfig.String("redishost")
	POOL = PoolInitRedis(redisHost, "")
	POOL.Stats()

	//定义日志文件
	logs.SetLogger(logs.AdapterFile, `{"filename":"run.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

}

// redis连接池
func PoolInitRedis(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2, //空闲数
		IdleTimeout: 120 * time.Second,
		MaxActive:   1000, //最大数
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

/**
* 生成一个指定最小值和最大值的随机整数
* min		最小值
* max		最大值
 */
func getRandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}
