package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Testcases embedded from testcasesB.txt (one case per line).
const rawTestcases = `3 6 2 3 4 3 7 5 4 6 2 1 3
1 6 1
3 8 7 3 5 5 6 4 5 6 1 2 3
1 7 1
3 1 7 2 5 3 2 5 7 3 3 2 1
4 7 8 7 9 9 3 6 8 4 7 3 5 8 0 1 7 4 1 2 3
4 0 9 0 6 1 4 3 5 2 1 2 5 0 3 0 0 2 1 3 4
2 1 9 9 3 2 1
4 3 0 8 9 6 1 8 2 3 3 6 6 7 0 6 3 3 2 1 4
1 9 1
3 5 5 7 2 9 8 1 4 1 2 1 3
3 0 2 9 2 6 3 9 5 3 3 1 2
1 7 1
1 8 1
4 7 2 7 8 5 2 6 4 6 1 9 8 5 3 7 3 1 4 2 3
1 7 1
1 8 1
4 7 3 6 3 4 7 7 2 3 7 4 5 7 9 2 8 3 2 4 1
1 0 1
2 5 0 5 9 2 1
4 6 3 6 7 3 2 6 0 5 8 3 5 2 8 8 9 3 2 1 4
2 9 7 0 3 1 2
3 8 3 0 0 0 2 3 7 9 2 3 1
4 4 5 7 3 2 5 8 5 7 6 7 5 7 1 2 3 3 2 4 1
4 8 3 5 0 0 6 3 1 6 9 9 6 8 3 0 2 1 3 2 4
1 1 1
4 0 3 6 9 5 4 1 0 3 9 0 9 8 5 2 7 4 2 1 3
3 3 6 6 8 5 3 9 4 1 2 3 1
3 9 9 7 1 8 7 6 5 6 2 1 3
1 6 1
2 2 0 1 6 2 1
1 4 1
3 6 1 5 4 5 9 3 8 8 3 2 1
1 2 1
4 0 4 5 2 5 1 9 8 6 1 4 9 1 8 1 9 3 1 4 2
2 5 1 8 3 2 1
2 6 3 9 7 1 2
2 4 5 8 3 2 1
4 5 2 4 8 0 1 5 4 5 0 1 9 0 4 1 2 4 2 1 3
4 1 6 8 8 6 1 1 9 7 5 2 4 3 6 4 0 3 2 4 1
1 5 1
3 8 4 0 3 8 2 4 9 0 1 3 2
3 9 3 4 5 8 7 7 1 6 3 1 2
1 3 1
4 6 0 3 1 4 1 2 2 7 0 5 8 6 1 8 3 1 3 2 4
4 6 7 0 2 8 5 3 3 2 3 4 2 2 8 4 8 2 3 4 1
4 2 3 8 2 9 2 3 8 4 0 7 3 0 6 2 6 3 2 1 4
4 1 2 1 8 0 1 0 3 6 7 4 6 5 6 4 1 4 1 3 2
2 0 0 8 1 1 2
2 7 2 9 8 1 2
1 9 1
2 5 5 5 4 1 2
2 0 7 3 9 1 2
3 7 4 7 2 7 4 5 5 9 1 3 2
4 8 4 9 3 3 6 3 4 8 2 8 2 9 8 8 8 1 3 2 4
4 4 7 8 0 3 4 7 1 0 3 9 1 4 0 3 2 1 4 3 2
3 0 0 0 9 2 4 3 6 7 1 2 3
3 1 9 6 1 0 2 0 9 8 1 3 2
2 6 3 5 2 2 1
4 2 9 9 8 6 9 3 0 0 6 9 0 2 8 1 6 4 1 2 3
1 5 1
2 5 4 7 7 1 2
3 1 7 2 1 8 4 8 6 9 1 3 2
2 9 6 2 1 2 1
2 7 4 1 6 1 2
3 1 5 8 8 1 0 0 9 0 2 3 1
4 6 7 0 5 2 1 6 9 3 0 3 6 0 0 9 6 1 4 3 2
4 6 8 0 7 8 2 8 2 3 9 2 1 9 7 4 4 4 2 1 3
1 0 1
1 0 1
3 2 2 8 2 1 5 0 5 2 3 1 2
1 5 1
2 4 7 5 7 1 2
1 1 1
3 7 6 6 8 7 8 6 8 1 2 3 1
4 8 5 2 1 2 2 0 6 9 1 1 8 4 9 4 5 1 3 2 4
4 8 5 3 8 0 6 5 5 8 0 3 9 1 7 1 5 2 1 3 4
1 8 1
1 0 1
3 0 5 2 5 4 7 7 6 0 2 1 3
4 5 2 0 7 8 2 6 4 9 7 3 0 7 9 9 2 1 3 2 4
4 9 2 8 6 6 3 5 6 4 3 3 4 3 8 9 5 3 2 1 4
2 5 3 2 4 2 1
1 7 1
4 3 2 9 5 9 3 1 8 7 3 4 6 4 8 3 5 3 1 2 4
1 3 1
2 9 2 8 7 2 1
3 8 4 7 6 5 6 0 1 0 3 1 2
3 9 8 7 4 8 6 7 5 3 2 1 3
1 1 1
1 9 1
3 9 1 9 3 0 0 1 4 1 1 2 3
2 9 9 8 3 1 2
3 2 2 9 3 2 2 3 4 7 3 1 2
3 3 1 2 6 2 3 9 0 5 1 2 3
2 4 1 9 4 2 1
1 2 1
2 9 0 6 2 2 1
4 8 0 8 5 3 2 7 5 7 3 4 3 1 5 7 4 1 2 4 3
2 4 8 6 1 2 1`

func parseTestcases() []string {
	lines := strings.Split(strings.TrimSpace(rawTestcases), "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}

func solveCase(n int, matrix [][]int64, order []int) string {
	results := make([]int64, n)
	active := make([]int, 0, n)
	for k := n - 1; k >= 0; k-- {
		pivot := order[k]
		active = append(active, pivot)
		pivotRow := matrix[pivot]
		for i := 0; i < n; i++ {
			rowI := matrix[i]
			distIP := rowI[pivot]
			for j := 0; j < n; j++ {
				sum := distIP + pivotRow[j]
				if sum < rowI[j] {
					rowI[j] = sum
				}
			}
		}
		var cur int64
		for _, i := range active {
			rowI := matrix[i]
			for _, j := range active {
				cur += rowI[j]
			}
		}
		results[k] = cur
	}
	var sb strings.Builder
	for i, v := range results {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
}

func parseCase(line string) (int, [][]int64, []int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return 0, nil, nil, fmt.Errorf("empty case")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil || n <= 0 {
		return 0, nil, nil, fmt.Errorf("invalid n")
	}
	if len(fields) != 1+n*n+n {
		return 0, nil, nil, fmt.Errorf("expected %d numbers got %d", 1+n*n+n, len(fields))
	}
	matrix := make([][]int64, n)
	idx := 1
	for i := 0; i < n; i++ {
		matrix[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			val, _ := strconv.ParseInt(fields[idx], 10, 64)
			matrix[i][j] = val
			idx++
		}
	}
	order := make([]int, n)
	for i := 0; i < n; i++ {
		v, _ := strconv.Atoi(fields[idx+i])
		order[i] = v - 1
	}
	return n, matrix, order, nil
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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for idx, line := range testcases {
		n, matrix, order, err := parseCase(line)
		if err != nil {
			fmt.Printf("case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		expected := solveCase(n, matrix, order)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatInt(matrix[i][j], 10))
			}
			sb.WriteByte('\n')
		}
		for i, v := range order {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v + 1))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
