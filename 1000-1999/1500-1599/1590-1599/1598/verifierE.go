package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type query struct{ x, y int }

type testCase struct {
	input    string
	expected string
}

func solveE(n, m int, ops []query) string {
	grid := make([][]bool, n+2)
	rlen := make([][]int, n+2)
	dlen := make([][]int, n+2)
	for i := 0; i <= n+1; i++ {
		grid[i] = make([]bool, m+2)
		rlen[i] = make([]int, m+2)
		dlen[i] = make([]int, m+2)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			grid[i][j] = true
		}
	}
	var total int64
	for i := n; i >= 1; i-- {
		for j := m; j >= 1; j-- {
			r := 1
			if j < m {
				r += dlen[i][j+1]
			}
			d := 1
			if i < n {
				d += rlen[i+1][j]
			}
			rlen[i][j] = r
			dlen[i][j] = d
			total += int64(r + d - 1)
		}
	}
	var out strings.Builder
	for _, op := range ops {
		x, y := op.x, op.y
		prev := grid[x][y]
		grid[x][y] = !grid[x][y]
		queue := []query{{x, y}}
		recompute := func(i, j int, prev bool) {
			oldR := rlen[i][j]
			oldD := dlen[i][j]
			var oldC int64
			if prev {
				oldC = int64(oldR + oldD - 1)
			}
			var newR, newD int
			if grid[i][j] {
				newR = 1
				if j < m && grid[i][j+1] {
					newR += dlen[i][j+1]
				}
				newD = 1
				if i < n && grid[i+1][j] {
					newD += rlen[i+1][j]
				}
			}
			rlen[i][j] = newR
			dlen[i][j] = newD
			var newC int64
			if grid[i][j] {
				newC = int64(newR + newD - 1)
			}
			if oldC != newC {
				total += newC - oldC
			}
		}
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			pr := grid[p.x][p.y]
			if p.x == x && p.y == y {
				pr = prev
			}
			oldR := rlen[p.x][p.y]
			oldD := dlen[p.x][p.y]
			recompute(p.x, p.y, pr)
			if rlen[p.x][p.y] != oldR || dlen[p.x][p.y] != oldD {
				if p.x > 1 {
					queue = append(queue, query{p.x - 1, p.y})
				}
				if p.y > 1 {
					queue = append(queue, query{p.x, p.y - 1})
				}
			}
		}
		out.WriteString(fmt.Sprintln(total))
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		q := rng.Intn(5) + 1
		ops := make([]query, q)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
		for j := 0; j < q; j++ {
			x := rng.Intn(n) + 1
			y := rng.Intn(m) + 1
			ops[j] = query{x, y}
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		exp := solveE(n, m, ops)
		cases[i] = testCase{input: sb.String(), expected: exp}
	}
	return cases
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
