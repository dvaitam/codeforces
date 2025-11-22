package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

const letters = "ab"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(2)]
	}
	return string(b)
}

// parseCandidate parses the output produced by a solution. It expects the
// number of operations followed by that many pairs of integers.
func parseCandidate(out string) ([][2]int, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, fmt.Errorf("unexpected end of output")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, fmt.Errorf("invalid integer %q", scanner.Text())
		}
		return v, nil
	}

	k, err := nextInt()
	if err != nil {
		return nil, err
	}
	if k < 0 || k > 200000 {
		return nil, fmt.Errorf("invalid operation count %d", k)
	}
	ops := make([][2]int, 0, k)
	for i := 0; i < k; i++ {
		a, err := nextInt()
		if err != nil {
			return nil, err
		}
		b, err := nextInt()
		if err != nil {
			return nil, err
		}
		ops = append(ops, [2]int{a, b})
	}
	return ops, nil
}

// applyOps simulates the prefix swap operations.
func applyOps(s, t string, ops [][2]int) (string, string, error) {
	for idx, op := range ops {
		a, b := op[0], op[1]
		if a < 0 || a > len(s) || b < 0 || b > len(t) {
			return s, t, fmt.Errorf("operation %d uses invalid prefix lengths %d and %d", idx+1, a, b)
		}
		u := s[:a]
		v := t[:b]
		s = v + s[a:]
		t = u + t[b:]
	}
	return s, t, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierD.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	rand.Seed(1)
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		s := randString(n)
		tStr := randString(m)
		if !strings.ContainsRune(s, 'a') {
			s = "a" + s[1:]
		}
		if !strings.ContainsRune(tStr, 'b') {
			tStr = "b" + tStr[1:]
		}
		var b bytes.Buffer
		fmt.Fprintln(&b, s)
		fmt.Fprintln(&b, tStr)
		input := b.String()

		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}

		ops, err := parseCandidate(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output: %v\n", t+1, err)
			os.Exit(1)
		}
		finalS, finalT, err := applyOps(s, tStr, ops)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid operations: %v\n", t+1, err)
			os.Exit(1)
		}
		if finalS != finalT {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sfound %d operations\nfinal s: %s\nfinal t: %s\n", t+1, input, len(ops), finalS, finalT)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
