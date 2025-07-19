package main

import (
	"fmt"
)

func main() {
	var k, l int
	if _, err := fmt.Scan(&k, &l); err != nil {
		return
	}
	w := 0
	for k != 0 && l%k == 0 {
		l /= k
		w++
	}
	if w == 0 || l > 1 {
		fmt.Println("NO")
	} else {
		fmt.Println("YES")
		fmt.Println(w - 1)
	}
}
