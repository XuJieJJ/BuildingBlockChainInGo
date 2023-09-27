package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

//反序列化

type Student struct {
	Name string
	Age  int
}

func main() {
	s1 := Student{
		Name: "demo1",
		Age:  20,
	}
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(s1)
	if err != nil {
		fmt.Println("编码失败：", err)
		return
	}

	//初始化解码器
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	var s2 Student
	fmt.Println("解码之前s2=", s2)

	decoder.Decode(&s2)
	fmt.Println("解码之后=", s2)

}
