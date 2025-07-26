package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const solutionB = "0-999/900-999/900-909/908/908B.go"

func runCmd(cmd []string, input string) (string, error) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	rand.Seed(42)
	var tests []string
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 2
		grid := make([]string, n)
		sx, sy := rand.Intn(n), rand.Intn(m)
		ex, ey := rand.Intn(n), rand.Intn(m)
		for ex == sx && ey == sy {
			ex, ey = rand.Intn(n), rand.Intn(m)
		}
		for r := 0; r < n; r++ {
			var sb strings.Builder
			for c := 0; c < m; c++ {
				ch := '.'
				if rand.Intn(5) == 0 {
					ch = '#'
				}
				if r == sx && c == sy {
					ch = 'S'
				} else if r == ex && c == ey {
					ch = 'E'
				}
				sb.WriteByte(byte(ch))
			}
			grid[r] = sb.String()
		}
		instrLen := rand.Intn(10) + 1
		var instr strings.Builder
		for j := 0; j < instrLen; j++ {
			instr.WriteByte(byte('0' + rand.Intn(4)))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, row := range grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		sb.WriteString(instr.String())
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	return tests
}

func runCandidate(path string, input string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		return runCmd([]string{"go", "run", path}, input)
	}
	os.Chmod(path, 0755)
	return runCmd([]string{path}, input)
}

func expectedOutput(input string) (string, error) {
	return runCmd([]string{"go", "run", solutionB}, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	tests := generateTests()
	for idx, t := range tests {
		exp, err := expectedOutput(t)
		if err != nil {
			fmt.Printf("failed to run reference solution on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(cand, t)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", idx+1, t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
