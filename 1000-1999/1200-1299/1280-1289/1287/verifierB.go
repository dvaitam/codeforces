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
	input  string
	output string
}

func third(a, b byte) byte {
	if a == b {
		return a
	}
	if a != 'S' && b != 'S' {
		return 'S'
	}
	if a != 'E' && b != 'E' {
		return 'E'
	}
	return 'T'
}

func expected(n, k int, cards []string) int64 {
	index := make(map[string]int, n)
	for i, c := range cards {
		index[c] = i
	}
	var count int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			b := make([]byte, k)
			for t := 0; t < k; t++ {
				b[t] = third(cards[i][t], cards[j][t])
			}
			if idx, ok := index[string(b)]; ok && idx > j {
				count++
			}
		}
	}
	return count
}

func buildTests() []testCase {
	r := rand.New(rand.NewSource(2))
	var tests []testCase
	// deterministic seed ensures stable tests
	// edge cases
	tests = append(tests, testCase{input: "1 1\nS\n", output: "0\n"})
	cards2 := []string{"S", "E", "T"}
	tests = append(tests, testCase{input: "3 1\nS\nE\nT\n", output: fmt.Sprintf("%d\n", expected(3, 1, cards2))})

	for len(tests) < 100 {
		n := r.Intn(6) + 3 // 3..8 to keep runtime small
		k := r.Intn(4) + 1 // 1..4
		cardSet := make(map[string]struct{}, n)
		cards := make([]string, 0, n)
		for len(cards) < n {
			b := make([]byte, k)
			for i := 0; i < k; i++ {
				switch r.Intn(3) {
				case 0:
					b[i] = 'S'
				case 1:
					b[i] = 'E'
				default:
					b[i] = 'T'
				}
			}
			s := string(b)
			if _, ok := cardSet[s]; !ok {
				cardSet[s] = struct{}{}
				cards = append(cards, s)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, c := range cards {
			sb.WriteString(c)
			sb.WriteByte('\n')
		}
		ans := expected(n, k, cards)
		tests = append(tests, testCase{input: sb.String(), output: fmt.Sprintf("%d\n", ans)})
	}
	return tests
}

func run(binary string, in string) (string, string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(in)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	cmd.Env = os.Environ()
	if err := cmd.Start(); err != nil {
		return "", errBuf.String(), err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return outBuf.String(), errBuf.String(), err
	case <-time.After(2 * time.Second):
		cmd.Process.Kill()
		return outBuf.String(), errBuf.String(), fmt.Errorf("timeout")
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := buildTests()
	for i, tc := range tests {
		out, errStr, err := run(binary, tc.input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", i+1, err, errStr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(tc.output) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, tc.output, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
