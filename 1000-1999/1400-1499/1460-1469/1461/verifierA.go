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
	n int
	k int
}

func genCase(r *rand.Rand) ([]testCase, string) {
	t := r.Intn(10) + 1
	cases := make([]testCase, t)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := r.Intn(1000) + 1
		k := r.Intn(n) + 1
		cases[i] = testCase{n: n, k: k}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	}
	return cases, sb.String()
}

func run(bin, input string) (string, error) {
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

func maxPalLen(s string) int {
	n := len(s)
	best := 0
	for center := 0; center < n; center++ {
		l, r := center, center
		for l >= 0 && r < n && s[l] == s[r] {
			if r-l+1 > best {
				best = r - l + 1
			}
			l--
			r++
		}
		l, r = center, center+1
		for l >= 0 && r < n && s[l] == s[r] {
			if r-l+1 > best {
				best = r - l + 1
			}
			l--
			r++
		}
	}
	return best
}

func validate(cases []testCase, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != len(cases) {
		return fmt.Errorf("expected %d lines got %d", len(cases), len(lines))
	}
	for i, tc := range cases {
		s := strings.TrimSpace(lines[i])
		if len(s) != tc.n {
			return fmt.Errorf("case %d: expected length %d got %d", i+1, tc.n, len(s))
		}
		for _, ch := range s {
			if ch != 'a' && ch != 'b' && ch != 'c' {
				return fmt.Errorf("case %d: invalid char %c", i+1, ch)
			}
		}
		if maxPalLen(s) > tc.k {
			return fmt.Errorf("case %d: longest palindrome > k", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cases, input := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validate(cases, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
