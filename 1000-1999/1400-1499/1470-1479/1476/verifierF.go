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

type testCase struct {
	input  string
	output string
}

func illuminated(n int, p []int, dir []byte) bool {
	for j := 0; j < n; j++ {
		ok := false
		for i := 0; i < n && !ok; i++ {
			if i == j {
				continue
			}
			if dir[i] == 'L' && i > j && i-p[i] <= j {
				ok = true
			}
			if dir[i] == 'R' && i < j && i+p[i] >= j {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func bruteForce(p []int) (bool, string) {
	n := len(p)
	dir := make([]byte, n)
	for mask := 0; mask < 1<<n; mask++ {
		for i := 0; i < n; i++ {
			if (mask>>i)&1 == 1 {
				dir[i] = 'R'
			} else {
				dir[i] = 'L'
			}
		}
		if illuminated(n, p, dir) {
			return true, string(dir)
		}
	}
	return false, ""
}

func buildCase(p []int) testCase {
	n := len(p)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	ok, dir := bruteForce(p)
	var out string
	if !ok {
		out = "NO\n"
	} else {
		out = fmt.Sprintf("YES\n%s\n", dir)
	}
	return testCase{input: sb.String(), output: out}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Intn(n)
	}
	return buildCase(p)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase([]int{1, 1}))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
