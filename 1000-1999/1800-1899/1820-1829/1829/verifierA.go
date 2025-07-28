package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const targetA = "codeforces"

func solveA(s string) string {
	diff := 0
	for i := 0; i < len(targetA) && i < len(s); i++ {
		if s[i] != targetA[i] {
			diff++
		}
	}
	return fmt.Sprint(diff)
}

func genTestsA() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		b := make([]byte, len(targetA))
		for j := range b {
			b[j] = byte('a' + rand.Intn(26))
		}
		s := string(b)
		fmt.Fprintln(&input, s)
		expected[i] = solveA(s)
	}
	return expected, input.String()
}

func runBinary(path, in string) ([]string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(&out)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsA()
	lines, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if len(lines) != len(expected) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expected), len(lines))
		os.Exit(1)
	}
	rand.Seed(1) // regenerate inputs for error messages
	for i, exp := range expected {
		if lines[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
	_ = time.Now()
}
