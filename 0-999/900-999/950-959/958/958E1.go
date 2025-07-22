package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct{ X, Y int64 }

var (
	ships []Point
	bases []Point
	N     int
	match []int
	used  []bool
)

func cross(a, b, c Point) int64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

func segmentsIntersect(a, b, c, d Point) bool {
	return cross(a, b, c)*cross(a, b, d) < 0 && cross(c, d, a)*cross(c, d, b) < 0
}

func search(i int) bool {
	if i == N {
		return true
	}
	for j := 0; j < N; j++ {
		if used[j] {
			continue
		}
		ok := true
		for k := 0; k < i; k++ {
			if segmentsIntersect(ships[i], bases[j], ships[k], bases[match[k]]) {
				ok = false
				break
			}
		}
		if ok {
			match[i] = j
			used[j] = true
			if search(i + 1) {
				return true
			}
			used[j] = false
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var R, B int
	if _, err := fmt.Fscan(in, &R, &B); err != nil {
		return
	}
	ships = make([]Point, R)
	for i := 0; i < R; i++ {
		fmt.Fscan(in, &ships[i].X, &ships[i].Y)
	}
	bases = make([]Point, B)
	for i := 0; i < B; i++ {
		fmt.Fscan(in, &bases[i].X, &bases[i].Y)
	}

	if R != B {
		fmt.Fprintln(out, "No")
		return
	}
	N = R
	match = make([]int, N)
	used = make([]bool, N)

	if search(0) {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}
