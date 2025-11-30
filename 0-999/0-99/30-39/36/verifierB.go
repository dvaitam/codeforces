package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"2 3\n.*\n.*",
	"3 2\n*..\n*.*\n*.*",
	"3 3\n..*\n...\n.*.",
	"3 3\n..*\n*.*\n..*",
	"3 1\n*..\n*.*\n*.*",
	"3 3\n**.\n*.*\n*.*",
	"3 1\n*..\n***\n.*.",
	"3 3\n*..\n...\n.**",
	"3 2\n*.*\n..*\n.**",
	"2 3\n*.\n**",
	"2 1\n..\n.*",
	"2 3\n..\n.*",
	"2 2\n.*\n..",
	"3 2\n...\n*.*\n***",
	"3 2\n..*\n***\n.*.",
	"3 3\n.*.\n..*\n...",
	"3 3\n*.*\n..*\n**.",
	"3 1\n..*\n..*\n*.*",
	"3 1\n...\n*..\n*.*",
	"2 1\n*.\n*.",
	"3 2\n*.*\n**.\n..*",
	"2 2\n*.\n*.",
	"3 3\n**.\n...\n...",
	"2 2\n*.\n*.",
	"3 1\n.*.\n*.*\n..*",
	"2 3\n*.\n..",
	"3 2\n*.*\n*..\n*..",
	"2 2\n..\n..",
	"3 1\n.*.\n...\n***",
	"3 3\n**.\n.*.\n..*",
	"3 2\n***\n..*\n*.*",
	"2 3\n**\n*.",
	"2 2\n..\n*.",
	"3 1\n*.*\n.**\n.*.",
	"3 3\n*.*\n...\n...",
	"3 1\n*..\n..*\n***",
	"2 1\n*.\n..",
	"2 1\n**\n.*",
	"2 1\n..\n*.",
	"2 2\n*.\n..",
	"3 1\n***\n**.\n**.",
	"3 2\n.*.\n.*.\n..*",
	"3 2\n*..\n.*.\n.**",
	"2 1\n**\n*.",
	"3 2\n*..\n.*.\n.**",
	"3 3\n*.*\n..*\n..*",
	"3 1\n.**\n**.\n..*",
	"2 2\n..\n.*",
	"2 2\n**\n.*",
	"3 1\n*..\n**.\n**.",
	"3 3\n..*\n...\n*..",
	"3 1\n.**\n**.\n***",
	"2 1\n*.\n*.",
	"2 1\n.*\n..",
	"3 3\n***\n*..\n***",
	"3 3\n***\n.**\n..*",
	"2 2\n**\n.*",
	"2 2\n.*\n**",
	"3 3\n.*.\n*.*\n.**",
	"2 1\n.*\n**",
	"3 1\n***\n.*.\n*.*",
	"2 2\n.*\n..",
	"2 3\n.*\n**",
	"2 3\n..\n**",
	"2 1\n**\n..",
	"3 3\n**.\n**.\n***",
	"2 2\n..\n**",
	"3 1\n*..\n.*.\n**.",
	"2 1\n..\n.*",
	"2 1\n.*\n**",
	"2 1\n..\n..",
	"2 2\n**\n*.",
	"2 1\n.*\n.*",
	"3 2\n...\n***\n.**",
	"3 2\n...\n...\n***",
	"3 1\n*..\n..*\n..*",
	"3 1\n.**\n...\n..*",
	"2 2\n..\n.*",
	"2 3\n*.\n**",
	"3 3\n**.\n...\n**.",
	"3 3\n*.*\n.*.\n*..",
	"2 1\n.*\n..",
	"2 3\n**\n.*",
	"2 1\n.*\n.*",
	"3 1\n*.*\n*.*\n.*.",
	"3 2\n.**\n.*.\n*..",
	"3 2\n**.\n.*.\n*..",
	"3 2\n*.*\n*..\n.**",
	"3 2\n***\n*..\n.**",
	"3 3\n.**\n***\n*..",
	"2 1\n**\n.*",
	"3 3\n***\n**.\n..*",
	"2 3\n..\n.*",
	"2 1\n*.\n.*",
	"2 3\n..\n..",
	"3 1\n***\n..*\n.**",
	"3 3\n.*.\n**.\n.**",
	"2 2\n*.\n**",
	"2 1\n*.\n.*",
	"3 3\n***\n***\n.*.",
}

func solveCase(n, k int, model []string) []string {
	cur := make([][]byte, n)
	for i := 0; i < n; i++ {
		cur[i] = []byte(model[i])
	}
	size := n
	for step := 2; step <= k; step++ {
		nextSize := size * n
		next := make([][]byte, nextSize)
		for i := range next {
			next[i] = make([]byte, nextSize)
			for j := range next[i] {
				next[i][j] = '.'
			}
		}
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if cur[i][j] == '*' {
					for di := 0; di < n; di++ {
						for dj := 0; dj < n; dj++ {
							next[i*n+di][j*n+dj] = '*'
						}
					}
				} else {
					for di := 0; di < n; di++ {
						for dj := 0; dj < n; dj++ {
							next[i*n+di][j*n+dj] = model[di][dj]
						}
					}
				}
			}
		}
		cur = next
		size = nextSize
	}
	res := make([]string, size)
	for i := 0; i < size; i++ {
		res[i] = string(cur[i])
	}
	return res
}

func parseCase(raw string) (int, int, []string, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 {
		return 0, 0, nil, fmt.Errorf("invalid case")
	}
	var n, k int
	if _, err := fmt.Sscan(lines[0], &n, &k); err != nil {
		return 0, 0, nil, err
	}
	if len(lines)-1 != n {
		return 0, 0, nil, fmt.Errorf("expected %d rows, got %d", n, len(lines)-1)
	}
	model := make([]string, n)
	copy(model, lines[1:])
	return n, k, model, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, raw := range rawTestcases {
		n, k, model, err := parseCase(raw)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expectedLines := solveCase(n, k, model)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for _, row := range model {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimRight(got, "\n"), "\n")
		if len(outLines) != len(expectedLines) {
			fmt.Printf("case %d failed: expected %d lines got %d\n", idx+1, len(expectedLines), len(outLines))
			os.Exit(1)
		}
		for i := range expectedLines {
			if outLines[i] != expectedLines[i] {
				fmt.Printf("case %d failed at line %d\nexpected: %s\ngot: %s\n", idx+1, i+1, expectedLines[i], outLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
