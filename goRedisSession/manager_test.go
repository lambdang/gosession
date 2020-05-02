/*
@Time : 2020/3/28 14:32 
@Author : dang
@desz:
*/
package goRedisSession

import(
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManager(t *testing.T) {
	//基于redis的测试
	m:=NewManage("127.0.0.1:6379",1,"ses",5,60*5,"0 0 2 * * * *")
	Convey("test ses manager",t,func(){
		//m:=SesMg()
		ses,_:=m.NewSes("dang")
		So(ses.GetPk(),ShouldEqual,"dang")

		nses:=m.Get(ses.GetSesId())
		So(nses.GetSesId(),ShouldEqual,ses.GetSesId())

		nses.Set("token","1234")
		m.Commit(nses)
		ses2:=m.Get(ses.GetSesId())
		So(ses2.Get("token").(string),ShouldEqual,"1234")


		nses.Set("token",123)
		m.Commit(nses)
		ses3:=m.Get(ses.GetSesId())
		So(ses3.Get("token").(float64),ShouldEqual,123)

		//测试过期

	})
}
