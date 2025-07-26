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

func solveQuery(s string, l, r int, prefix [][26]int) string {
	if l == r {
		return "Yes"
	}
	if s[l-1] != s[r-1] {
		return "Yes"
	}
	distinct := 0
	for c := 0; c < 26; c++ {
		if prefix[r][c]-prefix[l-1][c] > 0 {
			distinct++
		}
	}
	if distinct >= 3 {
		return "Yes"
	}
	return "No"
}

func expectedB(s string, queries [][2]int) []string {
	n := len(s)
	prefix := make([][26]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i]
		prefix[i+1][s[i]-'a']++
	}
	res := make([]string, len(queries))
	for i, q := range queries {
		res[i] = solveQuery(s, q[0], q[1], prefix)
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(20) + 1
		var sb strings.Builder
		sbytes := make([]byte, n)
		for i := 0; i < n; i++ {
			sbytes[i] = byte('a' + rng.Intn(26))
		}
		s := string(sbytes)
		q := rng.Intn(20) + 1
		queries := make([][2]int, q)
		fmt.Fprintf(&sb, "%s\n%d\n", s, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
			fmt.Fprintf(&sb, "%d %d\n", l, r)
		}
		input := sb.String()
		expected := expectedB(s, queries)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != q {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", t, q, len(lines), input)
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			if strings.TrimSpace(lines[i]) != expected[i] {
				fmt.Fprintf(os.Stderr, "case %d failed on query %d: expected %s got %s\ninput:\n%s", t, i+1, expected[i], lines[i], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
