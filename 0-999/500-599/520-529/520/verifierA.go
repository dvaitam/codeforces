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
	n int
	s string
}

func generateCase(rng *rand.Rand) (string, testCaseA) {
	n := rng.Intn(100) + 1
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	s := sb.String()
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, testCaseA{n: n, s: s}
}

func expected(tc testCaseA) string {
	seen := [26]bool{}
	for i := 0; i < tc.n; i++ {
		ch := tc.s[i]
		if ch >= 'a' && ch <= 'z' {
			seen[ch-'a'] = true
		} else if ch >= 'A' && ch <= 'Z' {
			seen[ch-'A'] = true
		}
	}
	for i := 0; i < 26; i++ {
		if !seen[i] {
			return "NO"
		}
	}
	return "YES"
}

func runCase(bin string, input string, tc testCaseA) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := expected(tc)
	if got != want {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
