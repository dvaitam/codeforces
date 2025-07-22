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

func solveB(s, b, e string) int {
	ns, nb, ne := len(s), len(b), len(e)
	nn := nb
	if ne > nn {
		nn = ne
	}
	fb := make([]bool, ns)
	fe := make([]bool, ns)
	for i := 0; i+nb <= ns; i++ {
		if s[i:i+nb] == b {
			fb[i] = true
		}
	}
	for i := 0; i+ne <= ns; i++ {
		if s[i:i+ne] == e {
			fe[i] = true
		}
	}
	p := make([]int, ns+1)
	q := make([]int, ns+1)
	h := make([]int, ns)
	for i := 0; i < ns; i++ {
		for j := i - 1; j >= 0; j-- {
			if s[i] == s[j] {
				q[j+1] = p[j] + 1
			} else {
				q[j+1] = 0
			}
			if q[j+1] > h[i] {
				h[i] = q[j+1]
			}
		}
		p, q = q, p
	}
	w := 0
	for i := ns - 1; i >= 0; i-- {
		if !fb[i] {
			continue
		}
		for j := i + nn; j <= ns; j++ {
			if fe[j-ne] && h[j-1] < j-i {
				w++
			}
		}
	}
	return w
}

func randStr(rng *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('a' + rng.Intn(4)))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 5
	s := randStr(rng, n)
	b := randStr(rng, rng.Intn(3)+1)
	e := randStr(rng, rng.Intn(3)+1)
	input := fmt.Sprintf("%s %s %s\n", s, b, e)
	exp := fmt.Sprintf("%d", solveB(s, b, e))
	return input, exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edge := []struct{ s, b, e string }{
		{"abc", "a", "c"},
		{"aaaa", "aa", "aa"},
		{"ababab", "ab", "ab"},
		{"abcde", "x", "y"},
		{"zzz", "z", "z"},
	}
	for i, tc := range edge {
		input := fmt.Sprintf("%s %s %s\n", tc.s, tc.b, tc.e)
		if err := runCase(exe, input, fmt.Sprintf("%d", solveB(tc.s, tc.b, tc.e))); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
