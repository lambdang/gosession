/*
@Time : 2020/3/16 23:31 
@Author : dang
@File : session_test
@Software: GoLand
*/
package goSession

import (
	"testing"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"time"
)

func TestSession(t *testing.T) {
	Convey("session test", t, func() {
		id := "123"
		id2 := "456"
		s := NewSes(id,"dang")
		s2 := NewSes(id2,"weiwu")
		//s.Set("name", "username")
		//s2.Set("name", "baba")
		So(s.GetPk(), ShouldEqual, "username")
		So(s2.GetPk(), ShouldEqual, "baba")

		//s.Set(SESSION_USERNAME, "luce")
		//s2.Set("name", "lili")
		//So(s.Get("name"), ShouldEqual, "luce")
		//So(s2.Get("name"), ShouldEqual, "lili")


		So(s.GetSesId(),ShouldEqual,id)
		So(s2.GetSesId(),ShouldEqual,id2)

		fmt.Println(s.Time())
		So(s.Time(),ShouldEqual,time.Now().Unix())
		So(s2.Time(),ShouldEqual,time.Now().Unix())


		s.Del("name")
		So(s.Get("name"),ShouldBeEmpty)
		fmt.Println(s.Get("name"))

	})

}
