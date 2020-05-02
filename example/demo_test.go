/*
@Time : 2020/5/2 14:14 
@Author : dang
@desz:
*/
package example

import (
	"fmt"
	"testing"
)

type Ione interface {
	SetName(name string)
	GetName() string
}

type One struct{
	name string
}

func(this *One) SetName(name string){
	this.name = name
}

func(this One) GetName() string{
	return this.name
}

func test1(o Ione){
	o.SetName("dang")
	o.GetName()
}

func test2(o *Ione){
	(*o).SetName("lili")
}

func TestD(t *testing.T) {
	o:=One{}
	test1(o)
	fmt.Println(o.GetName())

	//te	//fmt.Println(o.GetName())
}
