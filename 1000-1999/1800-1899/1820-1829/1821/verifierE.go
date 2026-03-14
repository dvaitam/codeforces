package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isRegular(s string) bool {
	bal := 0
	for _, ch := range s {
		if ch == '(' {
			bal++
		} else {
			bal--
		}
		if bal < 0 {
			return false
		}
	}
	return bal == 0
}

func minCostSeq(s string, memo map[string]int) int {
	if s == "" {
		return 0
	}
	if v, ok := memo[s]; ok {
		return v
	}
	n := len(s)
	best := math.MaxInt32
	for i := 0; i < n-1; i++ {
		if s[i] == '(' && s[i+1] == ')' {
			t := s[:i] + s[i+2:]
			c := n - (i + 2) + minCostSeq(t, memo)
			if c < best {
				best = c
			}
		}
	}
	memo[s] = best
	return best
}

func reachable(s string, k int) map[string]struct{} {
	type state struct {
		str   string
		moves int
	}
	q := []state{{s, 0}}
	vis := map[string]int{s: 0}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.moves == k {
			continue
		}
		n := len(cur.str)
		for i := 0; i < n; i++ {
			for j := 0; j <= n; j++ {
				if j == i || j == i+1 {
					continue
				}
				t := cur.str[:i] + cur.str[i+1:]
				jj := j
				if j > i {
					jj = j - 1
				}
				t = t[:jj] + string(cur.str[i]) + t[jj:]
				if v, ok := vis[t]; !ok || v > cur.moves+1 {
					vis[t] = cur.moves + 1
					q = append(q, state{t, cur.moves + 1})
				}
			}
		}
	}
	res := make(map[string]struct{})
	for str := range vis {
		if isRegular(str) {
			res[str] = struct{}{}
		}
	}
	return res
}

func solveE(k int, s string) int {
	best := math.MaxInt32
	for str := range reachable(s, k) {
		c := minCostSeq(str, map[string]int{})
		if c < best {
			best = c
		}
	}
	return best
}

func genRegular(rng *rand.Rand, n int) string {
	open := n / 2
	cls := open
	bal := 0
	var sb strings.Builder
	for open+cls > 0 {
		if open > 0 && (bal == 0 || rng.Intn(open+cls) < open) {
			sb.WriteByte('(')
			open--
			bal++
		} else {
			sb.WriteByte(')')
			cls--
			bal--
		}
	}
	return sb.String()
}

func genCaseE(rng *rand.Rand) (int, string) {
	n := rng.Intn(4)*2 + 2 // 2, 4, 6, or 8
	s := genRegular(rng, n)
	k := rng.Intn(4) // 0..3
	return k, s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/candidate")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	const numTests = 100
	type testCase struct {
		k int
		s string
	}
	cases := make([]testCase, numTests)
	expected := make([]int, numTests)

	// Build all test cases and compute expected answers.
	var inputBuf strings.Builder
	fmt.Fprintln(&inputBuf, numTests)
	for i := 0; i < numTests; i++ {
		k, s := genCaseE(rng)
		cases[i] = testCase{k, s}
		expected[i] = solveE(k, s)
		fmt.Fprintln(&inputBuf, k)
		fmt.Fprintln(&inputBuf, s)
	}

	// Run candidate once with all test cases.
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(inputBuf.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\nstderr: %s\n", err, stderr.String())
		os.Exit(1)
	}

	// Parse and compare outputs.
	lines := strings.Fields(strings.TrimSpace(stdout.String()))
	if len(lines) != numTests {
		fmt.Fprintf(os.Stderr, "expected %d output values, got %d\n", numTests, len(lines))
		os.Exit(1)
	}
	for i := 0; i < numTests; i++ {
		got, err := strconv.Atoi(lines[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output %q: %v\n", i+1, lines[i], err)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: k=%d s=%q expected %d got %d\n",
				i+1, cases[i].k, cases[i].s, expected[i], got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed")
}
