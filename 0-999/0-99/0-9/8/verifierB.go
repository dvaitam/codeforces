package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pt struct{ x, y int }

func solve(input string) string {
	s := strings.TrimSpace(input)
	positions := make([]pt, 0, len(s)+1)
	positions = append(positions, pt{0, 0})
	visited := map[pt]int{positions[0]: 0}
	x, y := 0, 0
	for i, ch := range s {
		switch ch {
		case 'L':
			x--
		case 'R':
			x++
		case 'U':
			y++
		case 'D':
			y--
		}
		cur := pt{x, y}
		if _, ok := visited[cur]; ok {
			return "BUG"
		}
		positions = append(positions, cur)
		visited[cur] = i + 1
	}
	dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i, p := range positions {
		for _, d := range dirs {
			np := pt{p.x + d.x, p.y + d.y}
			if j, ok := visited[np]; ok {
				if j != i-1 && j != i+1 {
					return "BUG"
				}
			}
		}
	}
	return "OK"
}

type test struct {
	input    string
	expected string
}

func generateTests() []test {
	rand.Seed(99)
	var tests []test
	fixed := []string{"LR", "UD", "LLRR", "URDL"}
	for _, f := range fixed {
		tests = append(tests, test{f + "\n", solve(f)})
	}
	moves := []byte{'L', 'R', 'U', 'D'}
	for len(tests) < 100 {
		n := rand.Intn(20) + 1
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteByte(moves[rand.Intn(4)])
		}
		inp := b.String() + "\n"
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
