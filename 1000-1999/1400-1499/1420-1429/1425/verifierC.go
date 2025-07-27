package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ x, y int }

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func knightDist(n, m int) [][]int {
	dist := make([][]int, n+1)
	for i := range dist {
		dist[i] = make([]int, m+1)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	q := []pair{{1, 1}}
	dist[1][1] = 0
	moves := []pair{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
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
	const mod = 1000000007
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		n := rand.Intn(5) + 3
		m := rand.Intn(5) + 3
		x := rand.Intn(n-2) + 3
		y := rand.Intn(m-2) + 3
		input := fmt.Sprintf("1\n%d %d %d %d\n", x, y, n, m)
		expect := fmt.Sprintf("%d", solve(x, y, n, m))
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", i+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
