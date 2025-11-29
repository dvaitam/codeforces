package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test data from testcasesD.txt.
const testData = `2 4 4
2 2 4
4 5 4 0 4
1 3
3 4 1 1
6 3 4 4 3 3 5
7 1 1 5 1 4 3 5
1 5
7 0 1 4 0 2 0 2
4 4 5 3 5
7 3 3 5 4 3 1 2
1 0
2 3 1
3 5 3 5
7 2 3 4 3 4 2 4
5 3 4 1 2 5
1 2
5 5 5 1 5 2
5 4 4 0 5 5
2 5 4
3 2 0 0
4 5 3 0 2
7 0 3 1 0 2 3 3
7 0 0 4 4 0 3 5
5 2 4 2 4 1
1 2
1 0
1 4
5 0 1 3 2 4
3 1 5 0
7 2 2 2 1 3 3 3
7 4 3 5 4 5 4 0
5 4 2 3 5 5
6 1 2 3 2 4 2
5 2 0 3 4 2
1 3
5 4 5 1 0 5
6 2 3 2 5 2 4
6 2 5 3 0 4 0
6 0 2 2 5 3 2
5 4 2 1 2 1
3 2 4 2
3 3 0 0
5 5 5 1 2 4
2 5 2
2 2 1
6 3 5 5 0 0 4
3 2 5 1
4 1 0 2 5
6 1 4 3 2 1 0
1 4
2 2 4
2 2 2
7 5 0 4 2 4 1 3
3 4 2 3
3 5 3 2
4 4 3 0 3
2 1 0
4 4 4 3 4
6 1 0 5 3 5 5
5 2 4 2 1 0
7 4 2 0 1 0 0 5
5 1 3 4 0 0
4 5 0 1 4
3 1 5 0
5 4 3 0 4 0
3 1 2 4
4 0 2 1 1
1 4
7 0 1 1 2 1 0 3
6 4 3 0 2 1 2
5 4 4 3 0 3
3 0 0 1
1 0
1 0
4 0 5 0 4
5 3 2 1 2 0
3 3 5 3
5 2 2 2 1 2
4 0 1 4 0
6 5 3 0 4 1 0
3 3 4 5
7 4 3 5 0 4 3 0
3 5 3 5
3 3 5 3
4 0 1 1 4
3 5 4 0
7 3 1 3 1 0 2 2
5 2 0 3 5 0
7 5 5 4 3 5 0 5
3 4 4 0
7 4 5 0 3 1 1 3
1 4
1 4
1 5
4 1 0 2 0
1 0
6 3 5 2 4 2 0
1 4
5 4 5 1 0 4`

func mex(arr []int) int {
	present := make(map[int]bool)
	for _, v := range arr {
		present[v] = true
	}
	m := 0
	for {
		if !present[m] {
			return m
		}
		m++
	}
}

func isSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
}

// Embedded solver logic from 1375D.go: produce the sequence of operations.
func solveOps(a []int) []int {
	n := len(a)
	ops := make([]int, 0, 2*n)

	isSorted := func() bool {
		for i := 1; i < n; i++ {
			if a[i] < a[i-1] {
				return false
			}
		}
		return true
	}

	for !isSorted() {
		present := make([]bool, n+1)
		for _, v := range a {
			if v >= 0 && v <= n {
				present[v] = true
			}
		}
		mex := 0
		for mex <= n && present[mex] {
			mex++
		}
		if mex < n {
			a[mex] = mex
			ops = append(ops, mex+1)
		} else {
			idx := -1
			for i := 0; i < n; i++ {
				if a[i] != i {
					idx = i
					break
				}
			}
			if idx == -1 {
				break
			}
			a[idx] = mex
			ops = append(ops, idx+1)
		}
	}
	return ops
}

func parseTests() ([]([]int), error) {
	lines := strings.Split(testData, "\n")
	tests := make([][]int, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n in line %q", line)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %q has wrong count", line)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(fields[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid value in line %q", line)
			}
			arr[i] = val
		}
		tests = append(tests, arr)
	}
	return tests, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for idx, arr := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(outBuf.String()))
		if len(outFields) == 0 {
			fmt.Printf("Test %d: empty output\n", idx+1)
			os.Exit(1)
		}
		m, err := strconv.Atoi(outFields[0])
		if err != nil {
			fmt.Printf("Test %d: invalid operations count\n", idx+1)
			os.Exit(1)
		}
		if len(outFields)-1 != m {
			fmt.Printf("Test %d: expected %d positions got %d\n", idx+1, m, len(outFields)-1)
			os.Exit(1)
		}
		ops := make([]int, m)
		for i := 0; i < m; i++ {
			pos, err := strconv.Atoi(outFields[i+1])
			if err != nil || pos < 1 || pos > len(arr) {
				fmt.Printf("Test %d: invalid position %q\n", idx+1, outFields[i+1])
				os.Exit(1)
			}
			ops[i] = pos - 1
		}
		expectedOps := solveOps(append([]int(nil), arr...))
		if len(expectedOps) != len(ops) {
			fmt.Printf("Test %d failed: expected %d operations got %d\n", idx+1, len(expectedOps), len(ops))
			os.Exit(1)
		}
		for i := range ops {
			if expectedOps[i]-1 != ops[i] {
				fmt.Printf("Test %d failed: expected position %d got %d at step %d\n", idx+1, expectedOps[i], ops[i]+1, i+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
