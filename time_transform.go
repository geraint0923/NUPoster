package main

import (
	"fmt"
	"time"
)

func TransformTime(t string) int64 {
	var year, month, day, hour, minute int
	year = 1970
	month = 1
	day = 1
	hour = 0
	minute = 0
	_, _ = fmt.Sscanf(t, "%d-%d-%dT%d:%d", &year, &month, &day, &hour, &minute)
	//	fmt.Printf("%d %d %d %d %d\n", year, month, day, hour, minute)
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC).UnixNano()
}

/*
func main() {
	tt := "2014-11-07T15:10"
	fmt.Println(tt)
	fmt.Println(strconv.FormatInt(TransformTime(tt), 10))
}
*/
