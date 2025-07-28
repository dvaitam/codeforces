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

const maxNG = 1000000

var prefixG [maxNG + 1]int64
var triG [2024]int

func init() {
	for i := 1; i <= maxNG; i++ {
		prefixG[i] = prefixG[i-1] + int64(i*i)
	}
	for i := 1; i < len(triG); i++ {
		triG[i] = triG[i-1] + i
	}
}

func solveG(n int) string {
	r := 1
	for triG[r] < n {
		r++
	}
	c := n - triG[r-1]
	var ans int64
	for i := 1; i <= r; i++ {
		j1 := c - (r - i)
		if j1 < 1 {
			j1 = 1
		}
		if j1 > i {
			continue
		}
		j2 := c
		if j2 > i {
			j2 = i
		}
		start := triG[i-1] + j1
		end := triG[i-1] + j2
		ans += prefixG[end] - prefixG[start-1]
	}
	return strconv.FormatInt(ans, 10)
}

func genTestsG() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(maxNG) + 1
		fmt.Fprintln(&input, n)
		expected[i] = solveG(n)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsG()
	lines, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if len(lines) != len(expected) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expected), len(lines))
		os.Exit(1)
	}
	for i, exp := range expected {
		if lines[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
