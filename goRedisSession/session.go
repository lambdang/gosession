/*
@Time : 2020/3/28 11:00 
@Author : dang
@desz: resis session
*/
package goRedisSession

import (
	"gosessions/core"
	"time"
)

type Session struct {
	Id      string                 `json:"id"` //token
	Pk      string                 `json:"pkName"`
	ActTime int64                  `json:"actTime"`
	Data    map[string]interface{} `json:"data"`
}

func NewSes(id, pk string) core.ISession {
	return &Session{
		id,
		pk,
		time.Now().Unix(),
		map[string]interface{}{},
	}
}


func (this Session) Time() int64 {
	return this.ActTime
}

func (this Session) GetPk() string {
	return this.Pk
}

func (this Session) GetSesId() string {
	return this.Id
}
func (this *Session) SetSesId(id string) {
	this.Id = id
}

func (this *Session) Get(key string) interface{} {
	return this.Data[key]
}

func (this *Session) Set(key string, value interface{}) {
	this.Data[key] = value
}
func (this *Session) Del(key string) {
	delete(this.Data, key)
}

func (this *Session) UTime() {
	this.ActTime = time.Now().Unix()
}
