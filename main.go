package main

import (
	"fmt"

	"github.com/github-real-lb/tutor-management-web/util"
)

func main() {
	for i := 0; i < 10; i++ {
		randomDiscount := util.RandomFloat64(0.0, 0.99)
		fmt.Println(randomDiscount)
	}

}
