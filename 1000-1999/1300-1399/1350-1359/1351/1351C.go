package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		x, y := 0, 0
		visited := make(map[[4]int]bool)
		time := 0
		for i := 0; i < len(s); i++ {
			nx, ny := x, y
			switch s[i] {
			case 'S':
				ny--
			case 'N':
				ny++
			case 'W':
				nx--
			case 'E':
				nx++
			}
			key := [4]int{x, y, nx, ny}
			if nx < x || (nx == x && ny < y) {
				key = [4]int{nx, ny, x, y}
			}
			if visited[key] {
				time += 1
			} else {
				visited[key] = true
				time += 5
			}
			x, y = nx, ny
		}
		fmt.Fprintln(out, time)
	}
}
