package main

import "fmt"

func main() {
	var a, b int
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	if a == 0 && b == 0 {
		fmt.Println("NO")
		return
	}
	if a-b < 0 {
		a, b = b, a
	}
	if a-b <= 1 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
