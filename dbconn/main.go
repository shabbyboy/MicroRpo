package main

import (
	"fmt"
	"os"
)

func SubAll(channel string, message string) {
	fmt.Println("all:", channel, message)
}

func SubOne(channel string, message string) {
	fmt.Println("one:", channel, message)
}

func main() {
	fmt.Println(os.Getwd())
}
