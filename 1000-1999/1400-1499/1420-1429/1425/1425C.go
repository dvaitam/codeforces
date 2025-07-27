package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

type pair struct{ x, y int }

// knightDist computes the minimum number of knight moves from (1,1) to
// every cell on an n by m board. unreachable cells get distance -1.
func knightDist(n, m int) [][]int {
	dist := make([][]int, n+1)
	for i := range dist {
		dist[i] = make([]int, m+1)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	moves := []pair{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
	q := []pair{{1, 1}}
	dist[1][1] = 0
	for head := 0; head < len(q); head++ {
		p := q[head]
		d := dist[p.x][p.y]
		for _, mv := range moves {
			nx, ny := p.x+mv.x, p.y+mv.y
			if nx >= 1 && nx <= n && ny >= 1 && ny <= m && dist[nx][ny] == -1 {
				dist[nx][ny] = d + 1
				q = append(q, pair{nx, ny})
			}
		}
	}
	return dist
}

func solve(x, y, n, m int) int {
	dist := knightDist(n, m)
	sum := 0
	for i := x; i <= n; i++ {
		for j := y; j <= m; j++ {
			d := dist[i][j]
			if d < 0 {
				return -1
			}
			sum = (sum + d) % mod
		}
	}
	return sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var X, Y, N, M int
		fmt.Fscan(in, &X, &Y, &N, &M)
		fmt.Fprintln(out, solve(X, Y, N, M))
	}
}
