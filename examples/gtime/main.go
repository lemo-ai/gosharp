package main

import (
	"fmt"
	"time"

	"github.com/lemo-ai/gosharp/gtime"
)

func main() {
	now := gtime.Now()
	fmt.Println("Now:", now.String())
	fmt.Println("Timestamp:", now.Timestamp())

	t, err := gtime.StrToTime("2024-01-02 15:04:05")
	if err != nil {
		panic(err)
	}
	fmt.Println("Parsed:", t.String())
	fmt.Println("After 1h:", t.Add(time.Hour).String())

	loc, err := gtime.SetTimeZone("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	fmt.Println("Package location:", loc)
	fmt.Println("gtime.Location():", gtime.Location())
}
