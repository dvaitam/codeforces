package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type seg struct{ r, l, rr int }

func solveB(n, m int, grid []string) int {
	weeds := []seg{}
	for i := 0; i < n; i++ {
		l, r := m+1, 0
		for j := 0; j < m; j++ {
			if grid[i][j] == 'W' {
				if j+1 < l {
					l = j + 1
				}
				if j+1 > r {
					r = j + 1
				}
			}
		}
		if r > 0 {
			weeds = append(weeds, seg{i + 1, l, r})
		}
	}
	if len(weeds) == 0 {
		return 0
	}
	k := len(weeds)
	const INF = int(1e9)
	dp0 := make([]int, k)
	dp1 := make([]int, k)
	f := weeds[0]
	base := f.r - 1
	dp0[0] = base + (2*f.rr - f.l - 1)
	dp1[0] = base + (f.rr - 1)
	for t := 1; t < k; t++ {
		c := weeds[t]
		p := weeds[t-1]
		gap := c.r - p.r
		nd0, nd1 := INF, INF
		start := p.l
		costL := min(abs(start-c.l)+2*(c.rr-c.l), abs(start-c.rr)+(c.rr-c.l))
		costR := min(abs(start-c.rr)+2*(c.rr-c.l), abs(start-c.l)+(c.rr-c.l))
		nd0 = min(nd0, dp0[t-1]+gap+costL)
		nd1 = min(nd1, dp0[t-1]+gap+costR)
		start = p.rr
		costL = min(abs(start-c.l)+2*(c.rr-c.l), abs(start-c.rr)+(c.rr-c.l))
		costR = min(abs(start-c.rr)+2*(c.rr-c.l), abs(start-c.l)+(c.rr-c.l))
		nd0 = min(nd0, dp1[t-1]+gap+costL)
		nd1 = min(nd1, dp1[t-1]+gap+costR)
		dp0[t], dp1[t] = nd0, nd1
	}
	res := dp0[k-1]
	if dp1[k-1] < res {
		res = dp1[k-1]
	}
	return res
}

func genCase() (string, int) {
	n := rand.Intn(10) + 1
	m := rand.Intn(10) + 1
	grid := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rand.Intn(2) == 0 {
				row[j] = '.'
			} else {
				row[j] = 'W'
			}
		}
		grid[i] = string(row)
		sb.WriteString(grid[i])
		if i+1 < n {
			sb.WriteByte('\n')
		}
	}
	expect := solveB(n, m, grid)
	return sb.String(), expect
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		in, expect := genCase()
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			fmt.Println(in)
			return
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
