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

func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m, r, c int
		fmt.Fscan(in, &n, &m, &r, &c)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &grid[i])
		}
		r--
		c--
		if grid[r][c] == 'B' {
			out.WriteString("0\n")
			continue
		}
		rowBlack := false
		for j := 0; j < m; j++ {
			if grid[r][j] == 'B' {
				rowBlack = true
				break
			}
		}
		colBlack := false
		for i := 0; i < n; i++ {
			if grid[i][c] == 'B' {
				colBlack = true
				break
			}
		}
		if rowBlack || colBlack {
			out.WriteString("1\n")
			continue
		}
		anyBlack := false
		for i := 0; i < n && !anyBlack; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == 'B' {
					anyBlack = true
					break
				}
			}
		}
		if anyBlack {
			out.WriteString("2\n")
		} else {
			out.WriteString("-1\n")
		}
	}
	return strings.TrimSpace(out.String())
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(1))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		m := rng.Intn(50) + 1
		r := rng.Intn(n) + 1
		c := rng.Intn(m) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, r, c))
		for row := 0; row < n; row++ {
			line := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(2) == 0 {
					line[j] = 'B'
				} else {
					line[j] = 'W'
				}
			}
			sb.WriteString(string(line) + "\n")
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveA(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
