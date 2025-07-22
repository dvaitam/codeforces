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

func solve(n, m int, grid [][]int) int {
	lo, hi := 0, int(1e9)
	rows := make([][]uint64, n)
	W := (m + 63) >> 6
	for i := range rows {
		rows[i] = make([]uint64, W)
	}
	check := func(x int) bool {
		good := make([]int, 0, n)
		for i := 0; i < n; i++ {
			cnt := 0
			for j := 0; j < m; j++ {
				if grid[i][j] >= x {
					w := j >> 6
					b := uint(j & 63)
					rows[i][w] |= 1 << b
					cnt++
				}
			}
			if cnt >= 2 {
				good = append(good, i)
			}
		}
		if len(good) < 2 {
			for i := 0; i < n; i++ {
				for k := range rows[i] {
					rows[i][k] = 0
				}
			}
			return false
		}
		for ii := 0; ii < len(good); ii++ {
			i := good[ii]
			for jj := ii + 1; jj < len(good); jj++ {
				k := good[jj]
				found := 0
				for w := 0; w < W; w++ {
					common := rows[i][w] & rows[k][w]
					if common == 0 {
						continue
					}
					found += bitsOn64(common)
					if found >= 2 {
						for t := 0; t < n; t++ {
							for u := range rows[t] {
								rows[t][u] = 0
							}
						}
						return true
					}
				}
			}
		}
		for i := 0; i < n; i++ {
			for k := range rows[i] {
				rows[i][k] = 0
			}
		}
		return false
	}
	for lo < hi {
		mid := (lo + hi + 1) >> 1
		if check(mid) {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func bitsOn64(x uint64) int { return int(strings.Count(fmt.Sprintf("%b", x), "1")) }

func runCase(bin string, n, m int, grid [][]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d ", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solve(n, m, grid)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, int, [][]int) {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = rng.Intn(100)
		}
	}
	return n, m, grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		n, m int
		g    [][]int
	}
	g1 := [][]int{{1, 2}, {3, 4}}
	cases = append(cases, struct {
		n, m int
		g    [][]int
	}{2, 2, g1})
	for i := 0; i < 100; i++ {
		n, m, g := genCase(rng)
		cases = append(cases, struct {
			n, m int
			g    [][]int
		}{n, m, g})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.n, tc.m, tc.g); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
