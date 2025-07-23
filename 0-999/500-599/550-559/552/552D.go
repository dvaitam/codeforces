package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	if n < 3 {
		fmt.Println(0)
		return
	}
	total := int64(n) * int64(n-1) * int64(n-2) / 6
	var deg int64
	for i := 0; i < n; i++ {
		mp := make(map[[2]int]int)
		xi, yi := points[i][0], points[i][1]
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			dx := points[j][0] - xi
			dy := points[j][1] - yi
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			mp[[2]int{dx, dy}]++
		}
		for _, c := range mp {
			if c >= 2 {
				deg += int64(c * (c - 1) / 2)
			}
		}
	}
	deg /= 3
	fmt.Println(total - deg)
}
