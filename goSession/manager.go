package goSession

import (
	"github.com/gorhill/cronexpr"
	"github.com/satori/go.uuid"
	"gosessions/core"
	"sync"
	"time"
)

//基于内存的session 简单就是力量
/**
1. sessions保存所有session对象
2. pk保存所有账号 禁止账号重复登录
 */

//var sesMg *Manager

type Manager struct {
	lock      sync.RWMutex
	fleshTime int64
	maxTime   int64
	sessions  map[string]core.ISession
	pk        map[string]string //账号<->token
	cron      string
}

//初始化ses管理器
func NewManager(rsT int64, maxTime int64, cron string) core.IManager {
	m := &Manager{
		fleshTime: rsT,
		maxTime:   maxTime,
		sessions:  make(map[string]core.ISession, 100000),
		pk:        make(map[string]string, 100000),
		cron:      cron,
	}
	//sesMg = m
	go m.Gc()
	return m
}

//ses id 生成器
func sessionId() string {
	return uuid.NewV4().String()
}

//创建sessin
func (this *Manager) NewSes(pk string) (core.ISession, error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	//当有该用户名时 则踢掉当前用户
	if _sid, ok := this.pk[pk]; ok {
		delete(this.sessions, _sid)
		delete(this.pk, pk)
	}

	sid := sessionId()
	ses := NewSes(sid, pk)
	this.pk[pk] = sid
	this.sessions[sid] = ses
	return ses, nil
}

//获取session 具有刷新机制
func (this *Manager) Get(sid string) core.ISession {
	this.lock.Lock()
	defer this.lock.Unlock()

	ses := this.sessions[sid]
	if ses == nil {
		return nil
	}
	now := time.Now().Unix()
	stime := ses.Time()

	if stime+this.maxTime < now { //过期
		if pk := ses.GetPk(); pk != "" {
			delete(this.pk, pk)
		}
		delete(this.sessions, sid)
		return nil
	} else if stime+this.fleshTime < now { //刷新
		newId := sessionId()
		ses.SetSesId(newId)
		delete(this.sessions, sid)
		//ses.Update()
		this.sessions[newId] = ses
		if _pk := ses.GetPk(); _pk != "" {
			this.pk[_pk] = newId
		}
		return ses
	}

	ses.UTime()
	return ses
}

//销毁session
func (this *Manager) Destroy(sid string) {
	//fmt.Println("destroy func")
	this.lock.Lock()
	defer this.lock.Unlock()
	obj := this.sessions[sid]
	userName := obj.GetPk()
	delete(this.pk, userName)
	delete(this.sessions, sid)
}

func (this *Manager) DestroyByPk(pk string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _sid, ok := this.pk[pk]; ok {
		delete(this.pk, pk)
		delete(this.sessions, _sid)
	}
}

func (this *Manager) Commit(se core.ISession) {
	ses := this.sessions[se.GetSesId()]
	ses.UTime()
}

// Gc清理
func (this *Manager) Gc() {
	var nextTime time.Time
	var now time.Time

	expr, err := cronexpr.Parse(this.cron)
	if err != nil {
		panic(err)
	}
	nextTime = expr.Next(time.Now())
	for {
		now = time.Now()
		tInt := time.Now().Unix()
		if nextTime.Before(now) || nextTime.Equal(now) {
			//this.lock.RLock()
			for k, v := range this.sessions {
				_t := v.Time()
				pk := v.GetPk()
				if _t+this.maxTime < tInt || _t+this.maxTime == tInt {
					this.lock.Lock()
					delete(this.sessions, k)
					delete(this.pk, pk)
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
