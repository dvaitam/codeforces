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

type pair struct{ r, c int }

func expected(T [][]byte, n int) string {
	maxD := 2 * n
	pos := make([][]pair, maxD+2)
	for d := 2; d <= maxD; d++ {
		for r := 1; r <= n; r++ {
			c := d - r
			if c >= 1 && c <= n {
				pos[d] = append(pos[d], pair{r, c})
			}
		}
	}
	succ := make([][][]int, maxD+2)
	for d := 2; d < maxD; d++ {
		m := len(pos[d])
		succ[d] = make([][]int, m)
		nextMap := make(map[int]map[int]int)
		for j, p := range pos[d+1] {
			if nextMap[p.r] == nil {
				nextMap[p.r] = make(map[int]int)
			}
			nextMap[p.r][p.c] = j
		}
		for i, p := range pos[d] {
			if p.r+1 <= n {
				if idx, ok := nextMap[p.r+1][p.c]; ok {
					succ[d][i] = append(succ[d][i], idx)
				}
			}
			if p.c+1 <= n {
				if idx, ok := nextMap[p.r][p.c+1]; ok {
					succ[d][i] = append(succ[d][i], idx)
				}
			}
		}
	}
	if n == 1 {
		if T[1][1] == 'a' {
			return "FIRST"
		} else if T[1][1] == 'b' {
			return "SECOND"
		} else {
			return "DRAW"
		}
	}

	dp := make([]map[int]int, maxD+2)
	for d := 0; d <= maxD+1; d++ {
		dp[d] = make(map[int]int)
	}
	var solve func(d, mask int) int
	solve = func(d, mask int) int {
		if d > maxD-1 {
			return 0
		}
		if v, ok := dp[d][mask]; ok {
			return v
		}
		turnFirst := (d % 2) == 0
		best := -1000000000
		if !turnFirst {
			best = 1000000000
		}
		for ch := 'a'; ch <= 'z'; ch++ {
			mask2 := 0
			for i := 0; i < len(pos[d]); i++ {
				if mask&(1<<i) != 0 {
					p := pos[d][i]
					if T[p.r][p.c] == byte(ch) {
						mask2 |= 1 << i
					}
				}
			}
			if mask2 == 0 {
				continue
			}
			nextMask := 0
			for i := 0; i < len(pos[d]); i++ {
				if mask2&(1<<i) != 0 {
					for _, j := range succ[d][i] {
						nextMask |= 1 << j
					}
				}
			}
			score := 0
			if ch == 'a' {
				score = 1
			} else if ch == 'b' {
				score = -1
			}
			val := score + solve(d+1, nextMask)
			if turnFirst {
				if val > best {
					best = val
				}
			} else {
				if val < best {
					best = val
				}
			}
		}
		dp[d][mask] = best
		return best
	}
	startScore := 0
	if T[1][1] == 'a' {
		startScore = 1
	} else if T[1][1] == 'b' {
		startScore = -1
	}
	initialMask := 0
	for _, j := range succ[2][0] {
		initialMask |= 1 << j
	}
	res := startScore + solve(3, initialMask)
	if res > 0 {
		return "FIRST"
	} else if res < 0 {
		return "SECOND"
	} else {
		return "DRAW"
	}
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		T := make([][]byte, n+1)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for r := 1; r <= n; r++ {
			row := make([]byte, n+1)
			for c := 1; c <= n; c++ {
				row[c] = byte('a' + rng.Intn(3))
			}
			T[r] = row
			for c := 1; c <= n; c++ {
				input.WriteByte(row[c])
			}
			if r < n {
				input.WriteByte('\n')
			} else {
				input.WriteByte('\n')
			}
		}
		expect := expected(T, n)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
