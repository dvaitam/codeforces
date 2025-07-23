package main

import "fmt"

func main() {
	var t, s, x int64
	if _, err := fmt.Scan(&t, &s, &x); err != nil {
		return
	}
	if x == t || (x >= t+s && ((x-t)%s == 0 || (x-t-1)%s == 0)) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
