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

const ALPHA = 26
const INF = 1000000000

func solve(str []string) (int, []string) {
	n := len(str)
	m := len(str[0])
	cost := make([][ALPHA][ALPHA]int, n)
	for i := 0; i < n; i++ {
		var even [ALPHA]int
		var odd [ALPHA]int
		row := str[i]
		for j := 0; j < m; j++ {
			c := int(row[j] - 'a')
			if j%2 == 0 {
				for x := 0; x < ALPHA; x++ {
					if x != c {
						even[x]++
					}
				}
			} else {
				for x := 0; x < ALPHA; x++ {
					if x != c {
						odd[x]++
					}
				}
			}
		}
		for x := 0; x < ALPHA; x++ {
			for y := 0; y < ALPHA; y++ {
				cost[i][x][y] = even[x] + odd[y]
			}
		}
	}
	var dpPrev [ALPHA][ALPHA]int
	var dpCurr [ALPHA][ALPHA]int
	pre := make([][ALPHA][ALPHA][2]int, n+1)
	for x := 0; x < ALPHA; x++ {
		for y := 0; y < ALPHA; y++ {
			if x != y {
				dpPrev[x][y] = 0
			} else {
				dpPrev[x][y] = INF
			}
		}
	}
	for i := 0; i < n; i++ {
		best := INF
		ba, bb := 0, 0
		for x := 0; x < ALPHA; x++ {
			for y := 0; y < ALPHA; y++ {
				if dpPrev[x][y] < best {
					best = dpPrev[x][y]
					ba, bb = x, y
				}
			}
		}
		for x := 0; x < ALPHA; x++ {
			for y := 0; y < ALPHA; y++ {
				dpCurr[x][y] = INF
				if x == y {
					continue
				}
				if x != ba && y != bb {
					dpCurr[x][y] = best + cost[i][x][y]
					pre[i+1][x][y][0] = ba
					pre[i+1][x][y][1] = bb
				} else {
					curBest := INF
					pa, pb := 0, 0
					for u := 0; u < ALPHA; u++ {
						if u == x {
							continue
						}
						for v := 0; v < ALPHA; v++ {
							if v == y {
								continue
							}
							if dpPrev[u][v] < curBest {
								curBest = dpPrev[u][v]
								pa, pb = u, v
							}
						}
					}
					dpCurr[x][y] = curBest + cost[i][x][y]
					pre[i+1][x][y][0] = pa
					pre[i+1][x][y][1] = pb
				}
			}
		}
		for x := 0; x < ALPHA; x++ {
			for y := 0; y < ALPHA; y++ {
				dpPrev[x][y] = dpCurr[x][y]
			}
		}
	}
	ansCost := INF
	a, b := 0, 0
	for x := 0; x < ALPHA; x++ {
		for y := 0; y < ALPHA; y++ {
			if dpPrev[x][y] < ansCost {
				ansCost = dpPrev[x][y]
				a, b = x, y
			}
		}
	}
	ans := make([]string, n)
	for i := n - 1; i >= 0; i-- {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if j%2 == 0 {
				row[j] = byte('a' + a)
			} else {
				row[j] = byte('a' + b)
			}
		}
		ans[i] = string(row)
		pa := pre[i+1][a][b][0]
		pb := pre[i+1][a][b][1]
		a, b = pa, pb
	}
	return ansCost, ans
}

func valid(output []string) bool {
	n := len(output)
	m := len(output[0])
	for i := 0; i < n; i++ {
		colors := make(map[byte]struct{})
		for j := 0; j < m; j++ {
			c := output[i][j]
			colors[c] = struct{}{}
			if j > 0 && output[i][j-1] == c {
				return false
			}
			if i > 0 && output[i-1][j] == c {
				return false
			}
		}
		if len(colors) > 2 {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, int, []string) {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	strs := make([]string, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			b[j] = byte('a' + rng.Intn(ALPHA))
		}
		strs[i] = string(b)
		sb.WriteString(strs[i] + "\n")
	}
	cost, _ := solve(strs)
	input := sb.String()
	return input, cost, strs
}

func runCase(bin, input string, expCost int, orig []string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) < 1+len(orig) {
		return fmt.Errorf("expected %d lines, got %d", 1+len(orig), len(lines))
	}
	var gotCost int
	if _, err := fmt.Sscan(lines[0], &gotCost); err != nil {
		return fmt.Errorf("failed to parse cost: %v", err)
	}
	grid := lines[1 : 1+len(orig)]
	if gotCost != expCost {
		return fmt.Errorf("expected cost %d got %d", expCost, gotCost)
	}
	if !valid(grid) {
		return fmt.Errorf("output grid invalid")
	}
	diff := 0
	for i := 0; i < len(orig); i++ {
		if len(grid[i]) != len(orig[0]) {
			return fmt.Errorf("row %d length mismatch", i)
		}
		for j := 0; j < len(orig[0]); j++ {
			if grid[i][j] != orig[i][j] {
				diff++
			}
		}
	}
	if diff != gotCost {
		return fmt.Errorf("reported cost %d but actual %d", gotCost, diff)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, cost, orig := generateCase(rng)
		if err := runCase(bin, input, cost, orig); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
