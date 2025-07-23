package main

import (
	"fmt"
)

func main() {
	var s, v1, v2, t1, t2 int
	if _, err := fmt.Scan(&s, &v1, &v2, &t1, &t2); err != nil {
		return
	}
	time1 := 2*t1 + s*v1
	time2 := 2*t2 + s*v2
	if time1 < time2 {
		fmt.Println("First")
	} else if time1 > time2 {
		fmt.Println("Second")
	} else {
		fmt.Println("Friendship")
	}
}
