package main

import "fmt"

type Io int

const (
	sunday Io = iota + 2
	monday
	tuesday
)

func main() {
	var io Io
	fmt.Println(io)
}
