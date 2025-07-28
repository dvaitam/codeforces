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
}

func solveCase(tc testCase) string {
	n := len(tc.s)
	i := 0
	for i < n {
		if tc.s[i] == 'W' {
			i++
			continue
		}
		j := i
		hasR := false
		hasB := false
		for j < n && tc.s[j] != 'W' {
			if tc.s[j] == 'R' {
				hasR = true
			}
			if tc.s[j] == 'B' {
				hasB = true
			}
			j++
		}
		if !(hasR && hasB) {
			return "NO"
		}
		i = j
	}
	return "YES"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	letters := []byte{'R', 'B', 'W'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(3)]
	}
	tc := testCase{s: string(b)}
	input := fmt.Sprintf("1\n%d\n%s\n", n, tc.s)
	output := solveCase(tc) + "\n"
	return input, output
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
