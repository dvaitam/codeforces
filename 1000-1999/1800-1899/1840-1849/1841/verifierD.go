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

type Segment struct {
	l int
	r int
}

type testCaseD struct {
	input    string
	expected string
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func minRemovals(segs []Segment) int {
	sort.Slice(segs, func(i, j int) bool {
		if segs[i].r == segs[j].r {
			return segs[i].l < segs[j].l
		}
		return segs[i].r < segs[j].r
	})
	n := len(segs)
	rvals := make([]int, n)
	for i := 0; i < n; i++ {
		rvals[i] = segs[i].r
	}
	dp := make([]int, n+1)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = prefix[i-1]
		li, ri := segs[i-1].l, segs[i-1].r
		for j := i - 1; j > 0; j-- {
			lj, rj := segs[j-1].l, segs[j-1].r
			if rj < li {
				break
			}
			if ri < lj {
				continue
			}
			start := li
			if lj < start {
				start = lj
			}
			p := sort.SearchInts(rvals, start) - 1
			cand := 2
			if p >= 0 {
				cand = prefix[p+1] + 2
			}
			if cand > dp[i] {
				dp[i] = cand
			}
		}
		if dp[i] > prefix[i-1] {
			prefix[i] = dp[i]
		} else {
			prefix[i] = prefix[i-1]
		}
	}
	return n - prefix[n]
}

func generateCaseD(rng *rand.Rand) testCaseD {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for j := 0; j < t; j++ {
		n := rng.Intn(6) + 2
		segs := make([]Segment, n)
		in.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			l := rng.Intn(20)
			r := l + rng.Intn(20)
			segs[i] = Segment{l, r}
			in.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		out.WriteString(fmt.Sprintf("%d\n", minRemovals(segs)))
	}
	return testCaseD{input: in.String(), expected: out.String()}
}

func runCaseD(bin string, tc testCaseD) error {
	got, err := runCandidate(bin, tc.input)
	if err != nil {
		return err
	}
	got = strings.TrimSpace(got)
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
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
	cases := []testCaseD{generateCaseD(rng)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCaseD(rng))
	}
	for i, tc := range cases {
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
