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
	s string
	t string
}

func expectedCase(tc testCase) string {
	s := tc.s
	t := tc.t
	if len(s) != len(t) {
		return "NO"
	}
	var diff []int
	for i := 0; i < len(s); i++ {
		if s[i] != t[i] {
			diff = append(diff, i)
			if len(diff) > 2 {
				return "NO"
			}
		}
	}
	if len(diff) != 2 {
		return "NO"
	}
	i, j := diff[0], diff[1]
	if s[i] == t[j] && s[j] == t[i] {
		return "YES"
	}
	return "NO"
}

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	s := randString(rng, n)
	t := randString(rng, m)
	for s == t {
		t = randString(rng, m)
	}
	input := fmt.Sprintf("%s\n%s\n", s, t)
	exp := fmt.Sprintf("%s\n", expectedCase(testCase{s, t}))
	return input, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
