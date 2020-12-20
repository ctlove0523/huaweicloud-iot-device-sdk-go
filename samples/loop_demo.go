package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.Tick(5 * time.Second)
	for ; ; {
		select {
		case <-ticker:
			fmt.Println(time.Now().Second())

		}
	}
}
