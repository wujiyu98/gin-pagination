package main

import (
	"fmt"
	pagination "github.com/wujiyu98/gin-pagination"
)

func main() {
	p := pagination.Init(1,10, 5,100, "/")
	fmt.Println(p.GetList())
	fmt.Println(p.BsPage())
}