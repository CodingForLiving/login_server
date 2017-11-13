package main

import (
	"log"
	"encoding/json"
	"github.com/garyburd/redigo/redis"	
)

type Db struct {
	redisPool *redis.Pool
	mysql MysqlClient
}

func (this *Db) Init() {
	pool := &redis.Pool{
		MaxIdle: 100,
		MaxActive: 100, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisAddr)
			if err != nil {
				log.Println("连接redis失败", err)
			}
			if config.RedisAuth != "" {
				c.Do("auth", config.RedisAuth)
			}

			c.Do("select", config.RedisIndex) 
			return c, err
		},
	}

	this.redisPool = pool

	this.mysql = MysqlClient{}
	this.mysql.init(config.MysqlStr)

	go this.mysql.loop()
}

func (this *Db) OnWxLogin(account *Account) string{
	str, err := json.Marshal(account)

	c := this.redisPool.Get()
	defer c.Close()

	// 插入redis
	n, err := c.Do("HSETNX", "Accounts", account.Openid, string(str))
	if err != nil {
		log.Println("账号存储redis报错:", err)
		return "账号存储redis报错:"
	}

	// 账号已经存在了
	if n.(int64) == 0 {
		log.Println("账号已存在，不需要新建")
		return "success"
	}

	//申请id->
	n, err = c.Do("incr", "IncrAccountId")
	if err != nil {
		log.Println("账号申请accountid出错", err)
		return "账号申请accountid出错"
	}

	account.Id = int(n.(int64))

	// 通知mysql协程入库
	this.mysql.c <- account
	return "success"
}

func (this *Db) OnRegister(account *Account) string {
	log.Println("OnRegister")
	str, err := json.Marshal(account)
	
	c := this.redisPool.Get()
	defer c.Close()

	// 插入redis
	n, err := c.Do("HSETNX", "Accounts", account.Account, string(str))
	if err != nil {
		log.Println("账号存储redis报错:", err)
		return "账号存储redis报错:"
	}

	// 账号已经存在了
	if n.(int64) == 0 {
		log.Println("账号已存在，不需要新建")
		return "success"
	}

	//申请id->
	n, err = c.Do("incr", "IncrAccountId")
	if err != nil {
		log.Println("账号申请accountid出错", err)
		return "账号申请accountid出错"
	}

	account.Id = int(n.(int64))

	// 通知mysql协程入库
	this.mysql.c <- account
	return "success"
}

func (this *Db) OnLogin(account *Account) string{
	c := this.redisPool.Get()
	defer c.Close()

	// 插入redis
	ret, err := c.Do("hget", "Accounts", account.Account)
	if err != nil {
		log.Println("查询redis账号数据出错:", err)
		return "账号存储redis报错:"
	}

	//str, err1 := json.Marshal(account)
	// 账号已经存在了

	str := ret.([]uint8)
	m := map[string]interface{}{}
	json.Unmarshal(str, &m)
	//json.Unmarshal(data, v)
	//申请id->
	return "success"
}

func (this *Db) Close() {
	this.redisPool.Close()
	this.mysql.NotifyClose()
}