/*
@Time : 2020/5/2 12:10 
@Author : dang
@desz:
*/
package core

type ISession interface {
	Time() int64
	GetPk() string
	GetSesId() string
	SetSesId(id string)
	Get(key string) interface{}
	Set(key string, value interface{})
	Del(key string)
	UTime()
}
