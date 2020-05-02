/*
@Time : 2020/3/16 23:14 
@Author : dang
@File : manager_test
@Software: GoLand
*/
package goSession

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestManage(t *testing.T) {
	m:=NewManager(3, 8, "*/1 * * * * * *")

	Convey("session manager base", t, func() {


		sess,_ := m.NewSes("dang")
		So(sess.GetSesId(), ShouldNotBeNil)
		//fmt.Println(sess.SesTime())
		So(sess.Time(), ShouldEqual, time.Now().Unix())
		So(sess.Get("name"), ShouldBeNil)
		sess.Set("name", "dangweiwu")
		So(sess.Get("name"), ShouldEqual, "dangweiwu")
		sid := sess.GetSesId()
		m.Destroy(sess.GetSesId())
		So(m.Get(sid), ShouldEqual, nil)
	})

	Convey("cron ", t, func() {
		//InitSesManager(3,6, "*/5 * * * * * *") //刷新3秒
		fmt.Println(time.Now().Unix())

		sess,_:= m.NewSes("dang")
		sid := sess.GetSesId()
		fmt.Println(sess.Time())

		time.Sleep(time.Second * 5)
		fmt.Println(sess.Time())
		sess2 := m.Get(sid)
		sid2 := sess2.GetSesId()
		fmt.Println(sess2.Time())
		So(sid, ShouldNotEqual, sid2)
	})

	Convey("cron ",t,func(){

		sess,_ := m.NewSes("dang")
		sid := sess.GetSesId()
		time.Sleep(time.Second * 10)
		sess2 := m.Get(sid)
		So(sess2, ShouldBeNil)
	})

	Convey("go ru", t,func(c C) {
		var f=func(){
			sess,_:=m.NewSes("dang")
			fmt.Println(sess.GetSesId())
			sess.Set(sess.GetSesId(),"dang")
			//So(sess.Get(sess.SesId()),ShouldEqual,"dang")
		}
		go f()
		go f()
		go f()
		time.Sleep(time.Second*3)
	})

	Convey("username 重复",t,func(c C){
		sess,_:= m.NewSes("dang")
		sess2,_:=m.NewSes("dang")
		So(sess.GetSesId(), ShouldNotEqual,sess2.GetSesId())
	})
}
