package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// runBinary executes the given binary with provided input and returns its output or error.
func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// expectedB computes the correct output for problem B.
func expectedB(input string) (string, error) {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return "", err
	}
	titles := make([]string, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(in, &titles[i]); err != nil {
			return "", err
		}
	}
	for length := 1; ; length++ {
		subs := make(map[string]struct{})
		for _, t := range titles {
			if len(t) < length {
				continue
			}
			for j := 0; j+length <= len(t); j++ {
				subs[t[j:j+length]] = struct{}{}
			}
		}
		buf := make([]byte, length)
		var found string
		var dfs func(int) bool
		dfs = func(pos int) bool {
			if pos == length {
				s := string(buf)
				if _, ok := subs[s]; !ok {
					found = s
					return true
				}
				return false
			}
			for c := byte('a'); c <= 'z'; c++ {
				buf[pos] = c
				if dfs(pos + 1) {
					return true
				}
			}
			return false
		}
		if dfs(0) {
			return fmt.Sprintf("%s\n", found), nil
		}
	}
}

// generateCase creates a deterministic random test case for problem B.
func generateCase(rng *rand.Rand) string {
	n := rng.Intn(30) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		l := rng.Intn(20) + 1
		var b strings.Builder
		for j := 0; j < l; j++ {
			b.WriteByte(byte('a' + rng.Intn(26)))
		}
		fmt.Fprintf(&sb, "%s\n", b.String())
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))

	var cases []string
	// fixed edge cases
	cases = append(cases, "1\na\n")
	cases = append(cases, "2\na\nb\n")
	cases = append(cases, "3\nabc\ndef\nxyz\n")
	cases = append(cases, "1\nabcdefghijklmnopqrstuvwxyz\n")

	for len(cases) < 100 {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		expect, err := expectedB(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected: %s got: %s", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
