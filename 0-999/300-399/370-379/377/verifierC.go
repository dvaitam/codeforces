package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseC struct {
	n         int
	strengths []int64
	m         int
	isPick    []bool
	team      []int
}

func generateCase(rng *rand.Rand) (string, testCaseC) {
	n := rng.Intn(6) + 2
	strengths := make([]int64, n)
	for i := range strengths {
		strengths[i] = int64(rng.Intn(50) + 1)
	}
	pickCount := rng.Intn(n/2) + 1
	banCount := rng.Intn(n/2) + 1
	for 2*(pickCount+banCount) > n {
		pickCount = rng.Intn(n/2) + 1
		banCount = rng.Intn(n/2) + 1
	}
	m := 2 * (pickCount + banCount)
	var ops []struct {
		pick bool
		team int
	}
	for t := 1; t <= 2; t++ {
		for i := 0; i < pickCount; i++ {
			ops = append(ops, struct {
				pick bool
				team int
			}{true, t})
		}
		for i := 0; i < banCount; i++ {
			ops = append(ops, struct {
				pick bool
				team int
			}{false, t})
		}
	}
	rng.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i, v := range strengths {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	fmt.Fprintf(&b, "%d\n", m)
	isPick := make([]bool, m)
	team := make([]int, m)
	for i, op := range ops {
		if op.pick {
			b.WriteString("p ")
		} else {
			b.WriteString("b ")
		}
		b.WriteString(fmt.Sprint(op.team))
		b.WriteByte('\n')
		isPick[i] = op.pick
		team[i] = op.team
	}
	tc := testCaseC{n: n, strengths: strengths, m: m, isPick: isPick, team: team}
	return b.String(), tc
}

func expectedAnswer(tc testCaseC) int64 {
	s := make([]int64, len(tc.strengths))
	copy(s, tc.strengths)
	sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })
	m := tc.m
	if m > len(s) {
		m = len(s)
	}
	a := make([]int64, m)
	for i := 0; i < m; i++ {
		a[i] = s[i]
	}
	total := 1 << m
	dp := make([]int64, total)
	for mask := total - 1; mask >= 0; mask-- {
		k := bits.OnesCount(uint(mask))
		if k >= m {
			dp[mask] = 0
			continue
		}
		pick := tc.isPick[k]
		t := tc.team[k]
		if t == 1 {
			best := int64(-1 << 60)
			for i := 0; i < m; i++ {
				if mask&(1<<i) != 0 {
					continue
				}
				val := dp[mask|1<<i]
				if pick {
					val += a[i]
				}
				if val > best {
					best = val
				}
			}
			dp[mask] = best
		} else {
			best := int64(1 << 60)
			for i := 0; i < m; i++ {
				if mask&(1<<i) != 0 {
					continue
				}
				val := dp[mask|1<<i]
				if pick {
					val -= a[i]
				}
				if val < best {
					best = val
				}
			}
			dp[mask] = best
		}
	}
	return dp[0]
}

func runCase(bin string, input string, tc testCaseC) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswer(tc))
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
