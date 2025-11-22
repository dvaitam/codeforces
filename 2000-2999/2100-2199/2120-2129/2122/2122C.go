package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type point struct {
	x, y int
	id   int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &pts[i].x, &pts[i].y)
			pts[i].id = i
		}

		idx := make([]int, n)
		for i := range idx {
			idx[i] = i
		}

		lowX := make([]bool, n)
		sort.Slice(idx, func(i, j int) bool {
			return pts[idx[i]].x < pts[idx[j]].x
		})
		for i := 0; i < n/2; i++ {
			lowX[idx[i]] = true
		}

		lowY := make([]bool, n)
		sort.Slice(idx, func(i, j int) bool {
			return pts[idx[i]].y < pts[idx[j]].y
		})
		for i := 0; i < n/2; i++ {
			lowY[idx[i]] = true
		}

		A, B, C, D := make([]int, 0), make([]int, 0), make([]int, 0), make([]int, 0)
		for i := 0; i < n; i++ {
			if lowX[i] && lowY[i] {
				A = append(A, i)
			} else if lowX[i] && !lowY[i] {
				B = append(B, i)
			} else if !lowX[i] && lowY[i] {
				C = append(C, i)
			} else {
				D = append(D, i)
			}
		}

		for i := 0; i < len(A); i++ {
			fmt.Fprintf(out, "%d %d\n", A[i]+1, D[i]+1)
		}
		for i := 0; i < len(B); i++ {
			fmt.Fprintf(out, "%d %d\n", B[i]+1, C[i]+1)
		}
	}
}
