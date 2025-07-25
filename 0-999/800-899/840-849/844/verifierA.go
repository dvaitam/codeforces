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

type testCaseA struct {
	s string
	k int
}

func expectedA(tc testCaseA) string {
	if len(tc.s) < tc.k {
		return "impossible"
	}
	seen := make(map[rune]struct{})
	for _, ch := range tc.s {
		seen[ch] = struct{}{}
	}
	distinct := len(seen)
	if distinct >= tc.k {
		return "0"
	}
	return fmt.Sprintf("%d", tc.k-distinct)
}

func runCase(bin string, tc testCaseA) error {
	input := fmt.Sprintf("%s\n%d\n", tc.s, tc.k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expectedA(tc)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func genCase(rng *rand.Rand) testCaseA {
	letters := "abcdefghijklmnopqrstuvwxyz"
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	k := rng.Intn(26) + 1
	return testCaseA{string(b), k}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseA, 100)
	for i := 0; i < 100; i++ {
		cases[i] = genCase(rng)
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
