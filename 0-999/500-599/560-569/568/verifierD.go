package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type line struct {
	a, b, c int
}

type testCase struct {
	n, k     int
	lines    []line
	expected string
}

func intersection(l1, l2 line) (float64, float64, bool) {
	det := float64(l1.a*l2.b - l2.a*l1.b)
	if math.Abs(det) < 1e-9 {
		return 0, 0, false
	}
	x := float64(l1.b*l2.c-l2.b*l1.c) / det
	y := float64(l1.c*l2.a-l2.c*l1.a) / float64(l1.b*l2.a-l1.a*l2.b)
	return x, y, true
}

func buildCandidates(lines []line) []int {
	n := len(lines)
	all := make([]int, 0)
	for i := 0; i < n; i++ {
		all = append(all, 1<<i) // sign on road i only
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			x, y, ok := intersection(lines[i], lines[j])
			if !ok {
				continue
			}
			mask := (1 << i) | (1 << j)
			for t := 0; t < n; t++ {
				if t == i || t == j {
					continue
				}
				if math.Abs(float64(lines[t].a)*x+float64(lines[t].b)*y+float64(lines[t].c)) < 1e-6 {
					mask |= 1 << t
				}
			}
			all = append(all, mask)
		}
	}
	// remove duplicates
	uniq := make(map[int]bool)
	res := make([]int, 0, len(all))
	for _, m := range all {
		if !uniq[m] {
			uniq[m] = true
			res = append(res, m)
		}
	}
	return res
}

func checkPossible(n, k int, lines []line) bool {
	cand := buildCandidates(lines)
	target := (1 << n) - 1
	var dfs func(idx, used int, mask int) bool
	dfs = func(idx, used int, mask int) bool {
		if mask == target && used <= k {
			return true
		}
		if used >= k || idx == len(cand) {
			return false
		}
		if dfs(idx+1, used, mask) {
			return true
		}
		return dfs(idx+1, used+1, mask|cand[idx])
	}
	return dfs(0, 0, 0)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	k := rng.Intn(3) + 1
	lines := make([]line, n)
	for i := 0; i < n; i++ {
		for {
			a := rng.Intn(5) - 2
			b := rng.Intn(5) - 2
			c := rng.Intn(5) - 2
			if a != 0 || b != 0 {
				lines[i] = line{a, b, c}
				break
			}
		}
	}
	ans := "NO\n"
	if checkPossible(n, k, lines) {
		ans = "YES"
	}
	return testCase{n: n, k: k, lines: lines, expected: ans}
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for _, ln := range tc.lines {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", ln.a, ln.b, ln.c))
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
	outFirst := strings.Fields(out.String())
	if len(outFirst) == 0 {
		return fmt.Errorf("no output")
	}
	got := strings.ToUpper(outFirst[0])
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
