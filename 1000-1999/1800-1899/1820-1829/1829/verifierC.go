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

type book struct {
	m int
	s string
}

func solveC(books []book) string {
	const inf = int(1e9)
	bestBoth := inf
	best1 := inf
	best2 := inf
	for _, b := range books {
		switch b.s {
		case "11":
			if b.m < bestBoth {
				bestBoth = b.m
			}
		case "10":
			if b.m < best1 {
				best1 = b.m
			}
		case "01":
			if b.m < best2 {
				best2 = b.m
			}
		}
	}
	ans := bestBoth
	if best1+best2 < ans {
		ans = best1 + best2
	}
	if ans >= inf {
		return "-1"
	}
	return strconv.Itoa(ans)
}

func genTestsC() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		fmt.Fprintln(&input, n)
		books := make([]book, n)
		for j := 0; j < n; j++ {
			m := rand.Intn(20) + 1
			opt := rand.Intn(4)
			var s string
			switch opt {
			case 0:
				s = "00"
			case 1:
				s = "01"
			case 2:
				s = "10"
			default:
				s = "11"
			}
			books[j] = book{m: m, s: s}
			fmt.Fprintf(&input, "%d %s\n", m, s)
		}
		expected[i] = solveC(books)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsC()
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
