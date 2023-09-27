package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// 序列化学习 gob
type People struct {
	Name string
	Age  int
}

func main() {
	p := People{
		Name: "Jacky one",
		Age:  18,
	}
	/**
	定义一个字节容器
	*/
	buf := bytes.Buffer{}
	//初始化编译器
	encoder := gob.NewEncoder(&buf)

	//编码操作
	err := encoder.Encode(p)
	if err != nil {
		fmt.Println("编码失败，错误原因:", err)
		return
	}
	fmt.Println(string(buf.Bytes()))

}
