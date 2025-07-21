package main

import (
	"fmt"
	"os"
)

func main() {
	var n, v int
	if _, err := fmt.Fscan(os.Stdin, &n, &v); err != nil {
		return
	}
	const maxDay = 3005
	fruits := make([]int, maxDay)
	var a, b int
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(os.Stdin, &a, &b)
		if a > maxA {
			maxA = a
		}
		fruits[a] += b
	}
	total := 0
	prev := 0
	// process days from 1 to maxA+1
	for day := 1; day <= maxA+1; day++ {
		takePrev := prev
		if takePrev > v {
			takePrev = v
		}
		total += takePrev
		cap := v - takePrev
		takeCur := fruits[day]
		if takeCur > cap {
			takeCur = cap
		}
		total += takeCur
		prev = fruits[day] - takeCur
	}
	fmt.Println(total)
}
