package main

import (
	"fmt"
	"regexp"
)

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
