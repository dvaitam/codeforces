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
	input    string
	expected string
}

func solveCase(k, a, b int, s string) string {
	n := len(s)
	if n < a*k || n > b*k {
		return "No solution"
	}
	var sb strings.Builder
	ix := 0
	for k > 0 {
		seg := (n - ix) / k
		sb.WriteString(s[ix : ix+seg])
		if k > 1 {
			sb.WriteByte('\n')
		}
		ix += seg
		k--
	}
	return sb.String()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randString(rng *rand.Rand, n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomCase(rng *rand.Rand) testCase {
	k := rng.Intn(5) + 1
	a := rng.Intn(5) + 1
	b := a + rng.Intn(5)
	n := rng.Intn(b*k-a*k+1) + a*k
	s := randString(rng, n)
	input := fmt.Sprintf("%d %d %d\n%s\n", k, a, b, s)
	return testCase{input: input, expected: solveCase(k, a, b, s)}
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
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1 1 1\na\n", expected: "a"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
