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

type cell struct{ x, y int }

func solve(n int, banned []cell) int {
	bannedRow := make([]bool, n+1)
	bannedCol := make([]bool, n+1)
	for _, c := range banned {
		bannedRow[c.x] = true
		bannedCol[c.y] = true
	}
	rowClear := make([]bool, n+1)
	colClear := make([]bool, n+1)
	for i := 2; i <= n-1; i++ {
		rowClear[i] = !bannedRow[i]
		colClear[i] = !bannedCol[i]
	}
	total := 0
	seen := make([]bool, n+1)
	for c := 2; c <= n-1; c++ {
		if seen[c] {
			continue
		}
		p := n + 1 - c
		seen[c] = true
		if p >= 2 && p <= n-1 {
			seen[p] = true
		}
		type node struct {
			kind byte
			idx  int
		}
		var nodes []node
		if colClear[c] {
			nodes = append(nodes, node{'T', c}, node{'B', c})
		}
		if rowClear[c] {
			nodes = append(nodes, node{'L', c}, node{'R', c})
		}
		if p != c && p >= 2 && p <= n-1 {
			if colClear[p] {
				nodes = append(nodes, node{'T', p}, node{'B', p})
			}
			if rowClear[p] {
				nodes = append(nodes, node{'L', p}, node{'R', p})
			}
		}
		K := len(nodes)
		if K == 0 {
			continue
		}
		conflict := make([][]bool, K)
		for i := range conflict {
			conflict[i] = make([]bool, K)
		}
		for i := 0; i < K; i++ {
			for j := i + 1; j < K; j++ {
				a, b := nodes[i], nodes[j]
				c1, c2 := a.idx, b.idx
				f := false
				if (a.kind == 'T' || a.kind == 'B') && (b.kind == 'T' || b.kind == 'B') && c1 == c2 {
					f = true
				}
				if (a.kind == 'L' || a.kind == 'R') && (b.kind == 'L' || b.kind == 'R') && c1 == c2 {
					f = true
				}
				if (a.kind == 'T' || a.kind == 'B') && (b.kind == 'L' || b.kind == 'R') {
					if a.kind == 'T' && b.kind == 'L' && c1 == c2 {
						f = true
					}
					if a.kind == 'B' && b.kind == 'R' && c1 == c2 {
						f = true
					}
					if a.kind == 'T' && b.kind == 'R' && c1+c2 == n+1 {
						f = true
					}
					if a.kind == 'B' && b.kind == 'L' && c1+c2 == n+1 {
						f = true
					}
				}
				if !f && (a.kind == 'L' || a.kind == 'R') && (b.kind == 'T' || b.kind == 'B') {
					if b.kind == 'T' && a.kind == 'L' && c1 == c2 {
						f = true
					}
					if b.kind == 'B' && a.kind == 'R' && c1 == c2 {
						f = true
					}
					if b.kind == 'T' && a.kind == 'R' && c1+c2 == n+1 {
						f = true
					}
					if b.kind == 'B' && a.kind == 'L' && c1+c2 == n+1 {
						f = true
					}
				}
				conflict[i][j] = f
				conflict[j][i] = f
			}
		}
		best := 0
		for mask := 0; mask < (1 << K); mask++ {
			cnt := bitsOn(mask)
			if cnt <= best {
				continue
			}
			ok := true
			for i := 0; i < K && ok; i++ {
				if mask&(1<<i) == 0 {
					continue
				}
				for j := i + 1; j < K; j++ {
					if mask&(1<<j) != 0 && conflict[i][j] {
						ok = false
						break
					}
				}
			}
			if ok {
				best = cnt
			}
		}
		total += best
	}
	return total
}

func bitsOn(x int) int {
	cnt := 0
	for x > 0 {
		cnt += x & 1
		x >>= 1
	}
	return cnt
}

func runCase(bin string, n int, banned []cell) error {
	input := fmt.Sprintf("%d %d\n", n, len(banned))
	for _, c := range banned {
		input += fmt.Sprintf("%d %d\n", c.x, c.y)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
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
	exp := solve(n, banned)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCase(rng *rand.Rand) (int, []cell) {
	n := rng.Intn(8) + 2 // 2..9
	m := rng.Intn(n)
	used := make(map[cell]struct{})
	banned := make([]cell, 0, m)
	for len(banned) < m {
		c := cell{rng.Intn(n) + 1, rng.Intn(n) + 1}
		if _, ok := used[c]; ok {
			continue
		}
		used[c] = struct{}{}
		banned = append(banned, c)
	}
	return n, banned
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct {
		n int
		b []cell
	}{{2, nil}, {3, []cell{{2, 2}}}}
	for i := 0; i < 100; i++ {
		n, b := genCase(rng)
		cases = append(cases, struct {
			n int
			b []cell
		}{n, b})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.n, tc.b); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
