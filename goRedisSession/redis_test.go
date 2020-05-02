/*
@Time : 2020/3/28 14:35 
@Author : dang
@desz:
*/
package goRedisSession

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {

	Convey("redis",t,func(){
		r:=NewRedis("127.0.0.1:6379",1,"ses")
		ses:=NewSes("123","0001")
		So(ses.GetPk(),ShouldEqual,"0001")
		r.RSet(ses,"10")
		data,_:=r.RGet(ses.GetSesId())
		So(data.GetPk(),ShouldEqual,"0001")
		b,_:=r.RExists(ses.GetSesId())
		So(b,ShouldEqual,true)
		time.Sleep(time.Second*11)
		b,_=r.RExists(ses.GetSesId())
		So(b,ShouldEqual,false)

		r.RSet(ses,"10")
		data,_=r.RGet(ses.GetSesId())
		So(data.GetPk(),ShouldEqual,"0001")
		r.RDel(ses.GetSesId())
		b,_=r.RExists(ses.GetSesId())
		So(b,ShouldEqual,false)

		err:=r.RDel("2222")
		fmt.Println("del err:",err)

		s,err:=r.RGet("xxxx")
		fmt.Println("s:",s,"err:",err)
		So(err,ShouldEqual,redis.ErrNil)



		//count:=0
		//for{
		//	count+=1;
		//	if count==100{
		//		break
		//	}
		//
		//	//sid,_:= strconv.Itoa(count)
		//	r.RSet(NewSes(strconv.Itoa(count),strconv.Itoa(count)),"3600")
		//}

		rl,i:=r.RSesKey("86")
		fmt.Println(rl,i)
	})

}

func TestRSesKey(t *testing.T) {
	r:=NewRedis("127.0.0.1:6379",1,"ses")


	count:=0
	//for{
	//	count+=1;
	//	if count==100{
	//		break
	//	}
	//
	//	//sid,_:= strconv.Itoa(count)
	//	r.RSet(NewSes(strconv.Itoa(count),strconv.Itoa(count)),"3600")
	//}
	count = 0
	tmp:="0"
	for{
		fmt.Println("wil ",tmp)
		data,_tmp:=r.RSesKey(tmp)
		tmp = _tmp
		count += len(data)

		fmt.Println(data,len(data))
		for _,v :=range data{
			fmt.Println(v)
		}
		if tmp=="0" || tmp==""{
			break
		}
		fmt.Println("tmp",tmp)
		//if count==90{
		//	break
		//}
	}
	fmt.Println("count:",count)
}
