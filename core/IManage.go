/*
@Time : 2020/5/2 12:08 
@Author : dang
@desz:
*/
package core

type IManager interface {
	NewSes(pk string) (ISession, error)
	Get(sid string) ISession
	Destroy(sid string)
	DestroyByPk(pk string)
	Commit(session ISession)
	Gc()
}

//NewSes(pk string) (*Session, error)
//Get(sid string) *Session