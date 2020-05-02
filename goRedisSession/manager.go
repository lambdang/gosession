/*
@Time : 2020/3/20 15:44 
@Author : dang
@File : Manager.go
@desz: redis实现的session
//具有过期刷新token功能
//基于redis

*/
package goRedisSession

import (
	"github.com/gorhill/cronexpr"
	"github.com/satori/go.uuid"
	"gosessions/core"
	"strconv"
	"sync"
	"time"
)

//var _sessionManager *Manager

type Manager struct {
	redis     *RedisSession
	lock      sync.RWMutex
	fleshTime int64
	maxTime   int64
	pk        map[string]string //唯一键:token pk:redis token
	gcCron    string            //定时清理
}


/**
host：redis host
db：redis db
prefix：redis 前缀
fleshTime: 刷新时间
maxtime：最长有效期
gcCron: 设定清理时间
 */
func NewManage(host string, db int, prefix string, fleshTime int64, maxTime int64, gcCron string) core.IManager {

	_sessionManager := &Manager{
		redis:     NewRedis(host, db, prefix),
		fleshTime: fleshTime,
		maxTime:   maxTime,
		gcCron:    gcCron,
		pk:  make(map[string]string, 100000),
	}
	_sessionManager.initPk() //内存缓存所有pk跟redis键
	go _sessionManager.Gc()
	return _sessionManager
}

func sessionId() string {
	return uuid.NewV4().String()
}


//创建session
func (this *Manager) NewSes(pk string) (core.ISession, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	//当有该用户时 则踢掉当前用户
	if _sid, ok := this.pk[pk]; ok {
		this.redis.RDel(_sid)
		delete(this.pk, pk)
	}

	sid := sessionId()
	ses := NewSes(sessionId(), pk)
	this.pk[pk] = sid
	this.redis.RSet(ses, strconv.Itoa(int(this.maxTime)))
	return ses, nil
}

//获取get 具有刷新机制
//获取session 具有刷新机制 过了刷新期会更新sid
func (this *Manager) Get(sid string) core.ISession {
	this.lock.Lock()
	defer this.lock.Unlock()
	ses, _ := this.redis.RGet(sid)
	if ses == nil {
		return nil
	}
	now := time.Now().Unix()
	t := ses.Time()
	if t+this.maxTime < now { //过期
		this.redis.RDel(ses.GetSesId())
		delete(this.pk, ses.GetPk())
		return nil
	} else if t+this.fleshTime < now { //过期刷新
		oldId := ses.GetSesId()
		this.redis.RDel(oldId)
		nsid := sessionId()
		ses.SetSesId(nsid)
		this.redis.RSet(ses, strconv.Itoa(int(this.maxTime))) //更新redis
		this.pk[ses.GetPk()] = nsid                  //更新username
		return ses
	}
	return ses
}

//提交
func (this *Manager) Commit(ses core.ISession) {
	ses.UTime()
	this.redis.RSet(ses, strconv.Itoa(int(this.maxTime)))
}

//销毁session
func (this *Manager) Destroy(sid string) {
	//fmt.Println("destroy func")
	this.lock.Lock()
	defer this.lock.Unlock()
	ses, _ := this.redis.RGet(sid)
	if ses == nil {
		return
	}
	this.redis.RDel(sid)
	pk := ses.GetPk()
	delete(this.pk, pk)
}

//销毁2
func (this *Manager) DestroyByPk(pk string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if sid, ok := this.pk[pk]; ok {
		this.redis.RDel(sid)
		delete(this.pk, pk)
	}
}

// Gc清理 内存清理
func (this *Manager) Gc() {
	var nextTime time.Time
	var now time.Time

	expr, err := cronexpr.Parse(this.gcCron)
	if err != nil {
		panic(err)
	}
	nextTime = expr.Next(time.Now())
	for {
		now = time.Now()
		if nextTime.Before(now) || nextTime.Equal(now) {
			//this.lock.RLock()
			for k, v := range this.pk {
				if r, _ := this.redis.RExists(v); !r {
					this.lock.Lock()
					delete(this.pk, k)
					this.lock.Unlock()
				}
			}

			//this.lock.RUnlock()
			//定时任务结束
			nextTime = expr.Next(time.Now())
		}
		select {
		case <-time.NewTimer(30 * 60).C: //半个小时检测一次超时
			//case <-time.NewTimer(time.Second).C: //每秒检测 测试用
		}
	}
}

//从redis中初始化session id 到 username
func (this *Manager) initPk() {
	count := 0
	tmp := "0"
	for {
		data, _tmp := this.redis.RSesKey(tmp)
		tmp = _tmp
		count += len(data)
		for _, i := range data {
			if ses, _ := this.redis.RGet(i); ses != nil {
				this.pk[ses.GetPk()] = ses.GetSesId()
			}
		}
		if tmp == "0" || tmp == "" {
			break
		}
	}

}
