package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point [5]int

func dot(a, b Point) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += a[i] * b[i]
	}
	return sum
}

func sub(a, b Point) Point {
	var r Point
	for i := 0; i < 5; i++ {
		r[i] = a[i] - b[i]
	}
	return r
}

func isAcute(p []Point, a, b, c int) bool {
	v1 := sub(p[b], p[a])
	v2 := sub(p[c], p[a])
	return dot(v1, v2) > 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		for j := 0; j < 5; j++ {
			fmt.Fscan(in, &pts[i][j])
		}
	}

	cand := []int{}
	for i := 0; i < n; i++ {
		cand = append(cand, i)
		changed := true
		for changed {
			changed = false
		outer:
			for x := 0; x < len(cand); x++ {
				for y := 0; y < len(cand); y++ {
					if y == x {
						continue
					}
					for z := y + 1; z < len(cand); z++ {
						if z == x {
							continue
						}
						if isAcute(pts, cand[x], cand[y], cand[z]) {
							cand = append(cand[:x], cand[x+1:]...)
							changed = true
							break outer
						}
						if isAcute(pts, cand[y], cand[x], cand[z]) {
							cand = append(cand[:y], cand[y+1:]...)
							changed = true
							break outer
						}
						if isAcute(pts, cand[z], cand[x], cand[y]) {
							cand = append(cand[:z], cand[z+1:]...)
							changed = true
							break outer
						}
					}
				}
			}
		}
	}

	result := []int{}
	for _, idx := range cand {
		good := true
		for j := 0; j < n && good; j++ {
			if j == idx {
				continue
			}
			for k := j + 1; k < n && good; k++ {
				if k == idx {
					continue
				}
				if dot(sub(pts[j], pts[idx]), sub(pts[k], pts[idx])) > 0 {
					good = false
				}
			}
		}
		if good {
			result = append(result, idx)
		}
	}

	fmt.Println(len(result))
	for _, idx := range result {
		fmt.Println(idx + 1)
	}
}
