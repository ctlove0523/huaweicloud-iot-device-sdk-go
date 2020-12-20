package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()
	//now = now.Add(8 * time.Hour)

	fmt.Println(now.Format("20060102T150405Z"))
	fmt.Println(now.Location())
	fmt.Println(now.Zone())
	fmt.Println(now.Unix())
}
