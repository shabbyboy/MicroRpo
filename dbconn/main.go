package main

import (
	"fmt"
	"regexp"
)

/*
此main 函数 主要是开发时 测试用，并无其他用处，可以删除
 */

func SubAll(channel string, message string) {
	fmt.Println("all:", channel, message)
}

func SubOne(channel string, message string) {
	fmt.Println("one:", channel, message)
}

func main() {
	match,_ := regexp.MatchString("aa.*.bb","aa.cc.bb")
	fmt.Println(match)
}
