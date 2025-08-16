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

func lowerBoundGt(v []int, x int) int {
	return sort.Search(len(v), func(i int) bool { return v[i] > x })
}

func advance(pos []int, last, length int) (int, bool) {
	if length == 0 {
		return last, true
	}
	start := lowerBoundGt(pos, last)
	idx := start + length - 1
	if idx < len(pos) {
		return pos[idx], true
	}
	return 0, false
}

func solveCaseE(n int, seq []int) int {
	pos := make([][]int, 8)
	for i, v := range seq {
		v--
		if v >= 0 && v < 8 {
			pos[v] = append(pos[v], i+1)
		}
	}
	minCnt := n
	for i := 0; i < 8; i++ {
		if len(pos[i]) < minCnt {
			minCnt = len(pos[i])
		}
	}
	ans := 0
	allMask := (1 << 8) - 1
	for base := 0; base <= minCnt; base++ {
		INF := n + 1
		dp := make([][9]int, 1<<8)
		for i := range dp {
			for j := 0; j <= 8; j++ {
				dp[i][j] = INF
			}
		}
		dp[0][0] = 0
		for mask := 0; mask < (1 << 8); mask++ {
			for p := 0; p <= 8; p++ {
				last := dp[mask][p]
				if last > n {
					continue
				}
				for k := 0; k < 8; k++ {
					if mask>>k&1 == 1 {
						continue
					}
					if npos, ok := advance(pos[k], last, base); ok {
						nm := mask | (1 << k)
						if npos < dp[nm][p] {
							dp[nm][p] = npos
						}
					}
					if p < 8 {
						if npos, ok := advance(pos[k], last, base+1); ok {
							nm := mask | (1 << k)
							if npos < dp[nm][p+1] {
								dp[nm][p+1] = npos
							}
						}
					}
				}
			}
		}
		for p := 0; p <= 8; p++ {
			if dp[allMask][p] <= n {
				if val := 8*base + p; val > ans {
					ans = val
				}
			}
		}
	}
	return ans
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	seq := make([]int, n)
	for i := range seq {
		seq[i] = rng.Intn(8) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range seq {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := fmt.Sprintf("%d", solveCaseE(n, seq))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
