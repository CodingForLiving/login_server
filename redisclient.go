package main

import (
	"log"
    "github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	conn redis.Conn
}

func (this *RedisClient) init() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println("Connect to redis error", err)
		return
	}
	this.conn = c
}

func (this *RedisClient) reconnect(){
}

func (this *RedisClient) setAccount(field string, value string) (interface{}, error) {
	n, err := this.conn.Do("HSETNX", "Accounts", field, value)
	return n, err
}

func (this *RedisClient) newAccountID() (interface{}, error) {
	n, err := this.conn.Do("incr", "IncrAccountId")
	return n, err
}