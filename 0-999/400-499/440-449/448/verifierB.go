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

func expectedB(s, t string) string {
	var freqS, freqT [26]int
	for _, ch := range s {
		freqS[ch-'a']++
	}
	for _, ch := range t {
		freqT[ch-'a']++
	}
	for i := 0; i < 26; i++ {
		if freqS[i] < freqT[i] {
			return "need tree"
		}
	}
	j := 0
	for i := 0; i < len(s) && j < len(t); i++ {
		if s[i] == t[j] {
			j++
		}
	}
	if j == len(t) {
		if len(s) != len(t) {
			return "automaton"
		}
		return "array"
	}
	if len(s) == len(t) {
		return "array"
	}
	return "both"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		sb1 := make([]rune, n)
		sb2 := make([]rune, m)
		for j := 0; j < n; j++ {
			sb1[j] = letters[rng.Intn(len(letters))]
		}
		for j := 0; j < m; j++ {
			sb2[j] = letters[rng.Intn(len(letters))]
		}
		s := string(sb1)
		t := string(sb2)
		input := fmt.Sprintf("%s\n%s\n", s, t)
		exp := expectedB(s, t)
		cases[i] = testCase{input: input, expected: exp}
	}
	return cases
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
