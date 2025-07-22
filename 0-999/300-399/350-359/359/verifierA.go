package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// solveCase computes the minimal operations using the algorithm from 359A.go
func solveCase(n, m int, grid [][]int) int {
	type rect struct {
		m1, m2 int
		r1, r2 int
		c1, c2 int
		bits   []uint64
	}
	nm := n * m
	words := (nm + 63) >> 6
	full := make([]uint64, words)
	for idx := 0; idx < nm; idx++ {
		full[idx>>6] |= 1 << (uint(idx) & 63)
	}
	rects := make([][]rect, 4)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != 1 {
				continue
			}
			// corner (1,1)
			{
				r := rect{m1: i + 1, m2: j + 1, r1: 0, r2: i, c1: 0, c2: j}
				r.bits = make([]uint64, words)
				for x := r.r1; x <= r.r2; x++ {
					for y := r.c1; y <= r.c2; y++ {
						idx := x*m + y
						r.bits[idx>>6] |= 1 << (uint(idx) & 63)
					}
				}
				rects[0] = append(rects[0], r)
			}
			// corner (1,m)
			{
				r := rect{m1: i + 1, m2: m - j, r1: 0, r2: i, c1: j, c2: m - 1}
				r.bits = make([]uint64, words)
				for x := r.r1; x <= r.r2; x++ {
					for y := r.c1; y <= r.c2; y++ {
						idx := x*m + y
						r.bits[idx>>6] |= 1 << (uint(idx) & 63)
					}
				}
				rects[1] = append(rects[1], r)
			}
			// corner (n,1)
			{
				r := rect{m1: n - i, m2: j + 1, r1: i, r2: n - 1, c1: 0, c2: j}
				r.bits = make([]uint64, words)
				for x := r.r1; x <= r.r2; x++ {
					for y := r.c1; y <= r.c2; y++ {
						idx := x*m + y
						r.bits[idx>>6] |= 1 << (uint(idx) & 63)
					}
				}
				rects[2] = append(rects[2], r)
			}
			// corner (n,m)
			{
				r := rect{m1: n - i, m2: m - j, r1: i, r2: n - 1, c1: j, c2: m - 1}
				r.bits = make([]uint64, words)
				for x := r.r1; x <= r.r2; x++ {
					for y := r.c1; y <= r.c2; y++ {
						idx := x*m + y
						r.bits[idx>>6] |= 1 << (uint(idx) & 63)
					}
				}
				rects[3] = append(rects[3], r)
			}
		}
	}
	for c := 0; c < 4; c++ {
		arr := rects[c]
		sort.Slice(arr, func(i, j int) bool {
			if arr[i].m1 != arr[j].m1 {
				return arr[i].m1 > arr[j].m1
			}
			return arr[i].m2 > arr[j].m2
		})
		filtered := make([]rect, 0, len(arr))
		best2 := -1
		for _, r := range arr {
			if r.m2 > best2 {
				filtered = append(filtered, r)
				best2 = r.m2
			}
		}
		rects[c] = filtered
	}
	// k=1
	for c := 0; c < 4; c++ {
		for _, r := range rects[c] {
			ok := true
			for w := 0; w < words; w++ {
				if r.bits[w] != full[w] {
					ok = false
					break
				}
			}
			if ok {
				return 1
			}
		}
	}
	// k=2
	for c1 := 0; c1 < 4; c1++ {
		for c2 := c1 + 1; c2 < 4; c2++ {
			for _, r1 := range rects[c1] {
				for _, r2 := range rects[c2] {
					ok := true
					for w := 0; w < words; w++ {
						if (r1.bits[w] | r2.bits[w]) != full[w] {
							ok = false
							break
						}
					}
					if ok {
						return 2
					}
				}
			}
		}
	}
	// k=3
	for c1 := 0; c1 < 4; c1++ {
		for c2 := c1 + 1; c2 < 4; c2++ {
			for c3 := c2 + 1; c3 < 4; c3++ {
				for _, r1 := range rects[c1] {
					for _, r2 := range rects[c2] {
						for _, r3 := range rects[c3] {
							ok := true
							for w := 0; w < words; w++ {
								if (r1.bits[w] | r2.bits[w] | r3.bits[w]) != full[w] {
									ok = false
									break
								}
							}
							if ok {
								return 3
							}
						}
					}
				}
			}
		}
	}
	return 4
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 3
	m := rng.Intn(8) + 3
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
	}
	// ensure at least one good cell not on a corner
	x := rng.Intn(n-2) + 1
	y := rng.Intn(m-2) + 1
	grid[x][y] = 1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == 0 || j == 0 || i == n-1 || j == m-1 {
				continue
			}
			if i == x && j == y {
				continue
			}
			if rng.Intn(2) == 1 {
				grid[i][j] = 1
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	exp := fmt.Sprintf("%d", solveCase(n, m, grid))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
