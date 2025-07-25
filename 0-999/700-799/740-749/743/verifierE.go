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

func bitsOnes(x int) int {
	cnt := 0
	for x > 0 {
		cnt++
		x &= x - 1
	}
	return cnt
}

func solveCaseE(n int, seq []int) int {
	pos := make([][]int, 8)
	for i, v := range seq {
		v--
		if v >= 0 && v < 8 {
			pos[v] = append(pos[v], i+1)
		}
	}
	ans1 := 0
	for i := 0; i < 8; i++ {
		if len(pos[i]) > 0 {
			ans1++
		}
	}
	ans := ans1
	if ans1 == 8 {
		INF := n + 1
		dp := make([][9]int, 1<<8)
		for m := range dp {
			for d := range dp[m] {
				dp[m][d] = INF
			}
		}
		dp[0][0] = 0
		for mask := 0; mask < (1 << 8); mask++ {
			t := bitsOnes(mask)
			for d := 0; d <= t; d++ {
				cur := dp[mask][d]
				if cur > n {
					continue
				}
				for i := 0; i < 8; i++ {
					if mask&(1<<i) != 0 {
						continue
					}
					arr := pos[i]
					j := sort.Search(len(arr), func(k int) bool { return arr[k] > cur })
					if j < len(arr) {
						m2 := mask | (1 << i)
						if arr[j] < dp[m2][d] {
							dp[m2][d] = arr[j]
						}
					}
					if len(arr) >= 2 {
						j1 := sort.Search(len(arr), func(k int) bool { return arr[k] > cur })
						if j1 < len(arr) {
							j2 := sort.Search(len(arr), func(k int) bool { return arr[k] > arr[j1] })
							if j2 < len(arr) {
								m2 := mask | (1 << i)
								if arr[j2] < dp[m2][d+1] {
									dp[m2][d+1] = arr[j2]
								}
							}
						}
					}
				}
			}
		}
		full := (1 << 8) - 1
		for d := 0; d <= 8; d++ {
			if dp[full][d] <= n {
				tot := 8 + d
				if tot > ans {
					ans = tot
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
