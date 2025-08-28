package main

import (
	"fmt"
	"task9/unpack"
)

func main() {
	str := "a4bc2d5e"
	res, _ := unpack.Unpack(str)
	fmt.Println(res)
}
