package main

import (
	"fmt"
)

func main() {
	var s, t string
	var n int64
	if _, err := fmt.Scan(&s, &t); err != nil {
		return
	}
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	pos := map[byte]int{'^': 0, '>': 1, 'v': 2, '<': 3}
	start := pos[s[0]]
	end := pos[t[0]]
	step := int(n % 4)
	cw := (start + step) % 4
	ccw := (start - step) % 4
	if ccw < 0 {
		ccw += 4
	}
	if cw == end && ccw == end {
		fmt.Println("undefined")
	} else if cw == end {
		fmt.Println("cw")
	} else if ccw == end {
		fmt.Println("ccw")
	} else {
		fmt.Println("undefined")
	}
}
