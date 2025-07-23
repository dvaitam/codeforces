package main

import "fmt"

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	var x1, y1, x2, y2 int
	if _, err := fmt.Scan(&x1, &y1, &x2, &y2); err != nil {
		return
	}
	var x, y int
	if _, err := fmt.Scan(&x, &y); err != nil {
		return
	}
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	if dx%x != 0 || dy%y != 0 {
		fmt.Println("NO")
		return
	}
	kx := dx / x
	ky := dy / y
	if kx%2 == ky%2 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
