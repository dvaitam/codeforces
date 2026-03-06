package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	s string
	c string
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func randWord(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('A' + rng.Intn(26))
	}
	return string(b)
}

func genTest(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintln(&sb, t)
	for i := 0; i < t; i++ {
		sLen := rng.Intn(6) + 2
		cLen := rng.Intn(6) + 1
		fmt.Fprintf(&sb, "%s %s\n", randWord(rng, sLen), randWord(rng, cLen))
	}
	return sb.String()
}

func parseInput(input string) ([]testCase, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty generated input")
	}
	var t int
	if _, err := fmt.Sscan(fields[0], &t); err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	if len(fields) != 1+2*t {
		return nil, fmt.Errorf("invalid token count")
	}
	tests := make([]testCase, 0, t)
	p := 1
	for i := 0; i < t; i++ {
		tests = append(tests, testCase{s: fields[p], c: fields[p+1]})
		p += 2
	}
	return tests, nil
}

func isAtMostOneSwap(from, to string) bool {
	if len(from) != len(to) {
		return false
	}
	diff := make([]int, 0, 2)
	for i := 0; i < len(from); i++ {
		if from[i] != to[i] {
			diff = append(diff, i)
			if len(diff) > 2 {
				return false
			}
		}
	}
	if len(diff) == 0 {
		return true
	}
	if len(diff) != 2 {
		return false
	}
	i, j := diff[0], diff[1]
	return from[i] == to[j] && from[j] == to[i]
}

func canMakeSmaller(s, c string) bool {
	if s < c {
		return true
	}
	b := []byte(s)
	n := len(b)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			b[i], b[j] = b[j], b[i]
			if string(b) < c {
				return true
			}
			b[i], b[j] = b[j], b[i]
		}
	}
	return false
}

func validateOutput(tc testCase, out string) error {
	got := strings.TrimSpace(out)
	if got == "---" {
		if canMakeSmaller(tc.s, tc.c) {
			return fmt.Errorf("reported impossible, but a valid string exists")
		}
		return nil
	}
	if !isAtMostOneSwap(tc.s, got) {
		return fmt.Errorf("output is not reachable by at most one swap")
	}
	if !(got < tc.c) {
		return fmt.Errorf("output is not lexicographically smaller than c")
	}
	return nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		pattern := "solbin*"
		if runtime.GOOS == "windows" {
			pattern = "solbin*.exe"
		}
		tmp, err := os.CreateTemp("", pattern)
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			_ = os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { _ = os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	argIdx := 1
	if os.Args[1] == "--" {
		argIdx = 2
	}
	if len(os.Args) <= argIdx {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate, cleanup, err := buildIfGo(os.Args[argIdx])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input := genTest(rng)
		tests, err := parseInput(input)
		if err != nil {
			fmt.Printf("internal parse error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotOut, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		lines := strings.Fields(gotOut)
		if len(lines) != len(tests) {
			fmt.Printf("test %d failed: expected %d output lines, got %d\ninput:\n%s\ngot:\n%s\n", i+1, len(tests), len(lines), input, strings.TrimSpace(gotOut))
			os.Exit(1)
		}
		for j, tc := range tests {
			if err := validateOutput(tc, lines[j]); err != nil {
				fmt.Printf("test %d case %d failed: %v\ninput:\n%s %s\ngot:\n%s\n", i+1, j+1, err, tc.s, tc.c, lines[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
