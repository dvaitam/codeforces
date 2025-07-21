package main

import (
	"bufio"
	"fmt"
	"os"
)

type move struct {
	c byte
	k int
}

func solve(a, b, n, m int, chips [][2]int, moves []move) (int, int) {
	chipSet := make(map[[2]int]struct{})
	for _, p := range chips {
		chipSet[p] = struct{}{}
	}
	r1, r2 := 1, a
	c1, c2 := 1, b
	alice, bob := 0, 0
	turn := 0
	for _, mv := range moves {
		var removed [][2]int
		switch mv.c {
		case 'U':
			for i := r1; i < r1+mv.k; i++ {
				for j := c1; j <= c2; j++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			r1 += mv.k
		case 'D':
			for i := r2 - mv.k + 1; i <= r2; i++ {
				for j := c1; j <= c2; j++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			r2 -= mv.k
		case 'L':
			for j := c1; j < c1+mv.k; j++ {
				for i := r1; i <= r2; i++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			c1 += mv.k
		case 'R':
			for j := c2 - mv.k + 1; j <= c2; j++ {
				for i := r1; i <= r2; i++ {
					p := [2]int{i, j}
					if _, ok := chipSet[p]; ok {
						removed = append(removed, p)
					}
				}
			}
			c2 -= mv.k
		}
		if turn%2 == 0 {
			alice += len(removed)
		} else {
			bob += len(removed)
		}
		for _, p := range removed {
			delete(chipSet, p)
		}
		turn++
	}
	return alice, bob
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, n, m int
		fmt.Fscan(in, &a, &b, &n, &m)
		chips := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &chips[i][0], &chips[i][1])
		}
		moves := make([]move, m)
		for i := 0; i < m; i++ {
			var c string
			fmt.Fscan(in, &c, &moves[i].k)
			moves[i].c = c[0]
		}
		x, y := solve(a, b, n, m, chips, moves)
		fmt.Fprintf(out, "%d %d\n", x, y)
	}
}
