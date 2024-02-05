package main

import (
	"fmt"

	"github.com/github-real-lb/tutor-management-web/util"
)

func main() {
	for i := 0; i < 10; i++ {
		randomDateTime := util.RandomDatetime()
		fmt.Println(randomDateTime)
	}

}
