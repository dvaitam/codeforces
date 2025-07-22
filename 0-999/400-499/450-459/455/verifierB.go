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

func runCandidate(bin string, input string) (string, error) {
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

type node struct {
	next [26]*node
	win  bool
	lose bool
}

func dfs(u *node) {
	has := false
	u.win = false
	u.lose = false
	for i := 0; i < 26; i++ {
		v := u.next[i]
		if v == nil {
			continue
		}
		has = true
		dfs(v)
		if !v.win {
			u.win = true
		}
		if !v.lose {
			u.lose = true
		}
	}
	if !has {
		u.lose = true
	}
}

func expected(n int, k int64, words []string) string {
	root := &node{}
	for _, s := range words {
		cur := root
		for i := 0; i < len(s); i++ {
			c := s[i] - 'a'
			if cur.next[c] == nil {
				cur.next[c] = &node{}
			}
			cur = cur.next[c]
		}
	}
	dfs(root)
	if !root.win {
		return "Second"
	}
	if root.lose {
		return "First"
	}
	if k%2 == 1 {
		return "First"
	}
	return "Second"
}

func buildCase(n int, k int64, words []string) (string, string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, w := range words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return sb.String(), expected(n, k, words)
}

func randString(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := int64(rng.Intn(5) + 1)
	words := make([]string, n)
	for i := range words {
		l := rng.Intn(4) + 1
		words[i] = randString(rng, l)
	}
	return buildCase(n, k, words)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct{ in, want string }
	// simple case
	{
		in, want := buildCase(1, 1, []string{"a"})
		cases = append(cases, struct{ in, want string }{in, want})
	}
	for i := 0; i < 100; i++ {
		in, want := generateCase(rng)
		cases = append(cases, struct{ in, want string }{in, want})
	}

	for idx, tc := range cases {
		got, err := runCandidate(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.want {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", idx+1, tc.want, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
