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

type query struct{ a, b, c, d int }

func bruteCount(grid [][]byte, q query) int {
	a, b, c, d := q.a, q.b, q.c, q.d
	count := 0
	for x1 := a; x1 <= c; x1++ {
		for y1 := b; y1 <= d; y1++ {
			for x2 := x1; x2 <= c; x2++ {
				for y2 := y1; y2 <= d; y2++ {
					ok := true
					for x := x1; x <= x2 && ok; x++ {
						for y := y1; y <= y2; y++ {
							if grid[x][y] == '1' {
								ok = false
								break
							}
						}
					}
					if ok {
						count++
					}
				}
			}
		}
	}
	return count
}

func solveCase(n, m int, grid [][]byte, qs []query) []int {
	ans := make([]int, len(qs))
	for i, q := range qs {
		ans[i] = bruteCount(grid, q)
	}
	return ans
}

func runCase(bin string, n, m, qCount int, grid [][]byte, qs []query) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, qCount))
	for i := 0; i < n; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	for _, qq := range qs {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", qq.a+1, qq.b+1, qq.c+1, qq.d+1))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := solveCase(n, m, grid, qs)
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		var val int
		fmt.Sscan(f, &val)
		if val != expected[i] {
			return fmt.Errorf("expected %v got %v", expected, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		qCount := rng.Intn(3) + 1
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			row := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(2) == 0 {
					row[j] = '0'
				} else {
					row[j] = '1'
				}
			}
			grid[i] = row
		}
		qs := make([]query, qCount)
		for i := 0; i < qCount; i++ {
			a := rng.Intn(n)
			c := rng.Intn(n)
			if a > c {
				a, c = c, a
			}
			b := rng.Intn(m)
			d := rng.Intn(m)
			if b > d {
				b, d = d, b
			}
			qs[i] = query{a, b, c, d}
		}
		if err := runCase(bin, n, m, qCount, grid, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
