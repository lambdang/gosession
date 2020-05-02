/*
@Time : 2020/3/21 15:08 
@Author : dang
@File : redisCtr
@Software: GoLand
*/
package goRedisSession

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"gosessions/core"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RedisSession struct {
	redisPool *redis.Pool
	prefix    string
}

func NewRedis(host string, db int, prefix string) *RedisSession {
	redisPool := &redis.Pool{
		MaxIdle:     256, //空闲等待 256
		MaxActive:   0,   //tcp最大连接数 无限
		IdleTimeout: time.Duration(120),
		Wait:        true,

		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", host,
				//redis.DialPassword(password),
				redis.DialDatabase(db),
				redis.DialConnectTimeout(1000*time.Millisecond),
				redis.DialReadTimeout(1000*time.Millisecond),
				redis.DialWriteTimeout(1000*time.Millisecond))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
	return &RedisSession{redisPool, prefix}
}

//redis 操作封装
func (this *RedisSession) RSet(ses core.ISession, ttl string) error {
	con := this.redisPool.Get()
	if con == nil {
		return errors.New("redis连接失败！")
	}
	defer con.Close()
	data, err := json.Marshal(ses)
	if err != nil {
		return err
	}
	_, err = con.Do("SET", fmt.Sprintf("%s:%s", this.prefix, ses.GetSesId()), data, "EX", ttl)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisSession) RGet(sid string) (core.ISession, error) {
	con := this.redisPool.Get()
	if con == nil {
		return nil, errors.New("redis连接失败！")
	}
	defer con.Close()

	data, err := redis.Bytes(con.Do("GET", fmt.Sprintf("%s:%s", this.prefix, sid)))
	if err != nil {
		return nil, err
	}
	ses := &Session{}
	err = json.Unmarshal(data, ses)
	if err != nil {
		return nil, err
	}
	return ses, nil
}

func (this *RedisSession) RExists(sid string) (bool, error) {
	con := this.redisPool.Get()
	if con == nil {
		return false, nil
	}
	defer con.Close()
	return redis.Bool(con.Do("EXISTS", fmt.Sprintf("%s:%s", this.prefix, sid)))

}

func (this *RedisSession) RDel(sid string) error {
	con := this.redisPool.Get()
	if con == nil {
		return nil
	}
	defer con.Close()
	_, err := con.Do("DEL", fmt.Sprintf("%s:%s", this.prefix, sid))
	if err != nil {
		return err
	}
	return nil
}


//获取sesionkey
func(this *RedisSession) RSesKey(index string) ([]string,string){
	con:=this.redisPool.Get()
	if con == nil{
		return nil,""
	}
	defer con.Close()
	r,err:=redis.Values(con.Do("SCAN",index,"MATCH",fmt.Sprintf("%s*",this.prefix)))
	if err!=nil{
		return nil,"0"
	}

	if len(r)!=0{
		fmt.Println("index",string(r[0].([]byte)))
		index:=string(r[0].([]byte))
		list:=r[1].([]interface{})
		rl:=[]string{}
		for _,v:=range list{
			rl = append(rl,string(v.([]byte)))
		}
		return rl,index

	}else{
		return nil,"0"
	}
}