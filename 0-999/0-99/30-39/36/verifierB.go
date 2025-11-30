package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Testcases embedded from testcasesB.txt (count + cases).
const rawTestcases = `100
2 3
.*
.*
3 2
*..
*.*
*.*
3 3
..*
...
.*.
3 3
..*
*.*
..*
3 1
*..
*.*
*.*
3 3
**.
*.*
*.*
3 1
*..
***
.*.
3 3
*..
...
.**
3 2
*.*
..*
.**
2 3
*.
**
2 1
..
.*
2 3
..
.*
2 2
.*
..
3 2
...
*.*
***
3 2
..*
***
.*.
3 3
.*.
..*
...
3 3
*.*
..*
**.
3 1
..*
..*
*.*
3 1
...
*..
*.*
2 1
*.
*.
3 2
*.*
**.
..*
2 2
*.
*.
3 3
**.
...
...
2 2
*.
*.
3 1
.*.
*.*
..*
2 3
*.
..
3 2
*.*
*..
*..
2 2
..
..
3 1
.*.
...
***
3 3
**.
.*.
..*
3 2
***
..*
*.*
2 3
**
*.
2 2
..
*.
3 1
*.*
.**
.*.
3 3
*.*
...
...
3 1
*..
..*
***
2 1
*.
..
2 1
**
.*
2 1
..
*.
2 2
*.
..
3 1
***
**.
**.
3 2
.*.
.*.
..*
3 2
*..
.*.
.**
2 1
**
*.
3 2
*..
.*.
.**
3 3
*.*
..*
..*
3 1
.**
**.
..*
2 2
..
.*
2 2
**
.*
3 1
*..
**.
**.
3 3
..*
...
*..
3 1
.**
**.
***
2 1
*.
*.
2 1
.*
..
3 3
***
*..
***
3 3
***
.**
..*
2 2
**
.*
2 2
.*
**
3 3
.*.
*.*
.**
2 1
.*
**
3 1
***
.*.
*.*
2 2
.*
..
2 3
.*
**
2 3
..
**
2 1
**
..
3 3
**.
**.
***
2 2
..
**
3 1
*..
.*.
**.
2 1
..
.*
2 1
.*
**
2 1
..
..
2 2
**
*.
2 1
.*
.*
3 2
...
***
.**
3 2
...
...
***
3 1
*..
..*
..*
3 1
.**
...
..*
2 2
..
.*
2 3
*.
**
3 3
**.
...
**.
3 3
*.*
.*.
*..
2 1
.*
..
2 3
**
.*
2 1
.*
.*
3 1
*.*
*.*
.*.
3 2
.**
.*.
*..
3 2
**.
.*.
*..
3 2
*.*
*..
.**
3 2
***
*..
.**
3 3
.**
***
*..
2 1
**
.*
3 3
***
**.
..*
2 3
..
.*
2 1
*.
.*
2 3
..
..
3 1
***
..*
.**
3 3
.*.
**.
.**
2 2
*.
**
2 1
*.
.*
3 3
***
***
.*.`

type testCase struct {
	n, k int
	model []string
}

func loadTestcases() ([]testCase, error) {
	reader := bufio.NewReader(strings.NewReader(rawTestcases))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		var n, k int
		if _, err := fmt.Fscan(reader, &n, &k); err != nil {
			return nil, fmt.Errorf("case %d header: %w", i+1, err)
		}
		model := make([]string, n)
		for r := 0; r < n; r++ {
			if _, err := fmt.Fscan(reader, &model[r]); err != nil {
				return nil, fmt.Errorf("case %d row %d: %w", i+1, r+1, err)
			}
		}
		cases = append(cases, testCase{n: n, k: k, model: model})
	}
	return cases, nil
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
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Printf("failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range testcases {
		n, k, model := tc.n, tc.k, tc.model
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
