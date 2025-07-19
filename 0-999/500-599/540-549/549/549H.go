package main

import (
	"fmt"
)

func main() {
	var a, b, c, d int64
	if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
		return
	}
	expr1 := abs64(a + b + c + d)
	expr2 := abs64(a - b + c - d)
	expr3 := abs64(a + b - c - d)
	expr4 := abs64(a - b - c + d)
	t := expr1
	if expr2 > t {
		t = expr2
	}
	if expr3 > t {
		t = expr3
	}
	if expr4 > t {
		t = expr4
	}
	if t == 0 {
		fmt.Print(0)
	} else {
		numerator := abs64(a*d - b*c)
		res := float64(numerator) / float64(t)
		fmt.Printf("%.10f", res)
	}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
