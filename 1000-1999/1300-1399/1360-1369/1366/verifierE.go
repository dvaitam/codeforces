package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded reference solution (1366E.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	b := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &b[i])
	}

	idx := n - 1
	ans := int64(1)
	for i := m - 1; i >= 0; i-- {
		last := -1
		for idx >= 0 && a[idx] >= b[i] {
			if a[idx] == b[i] && last == -1 {
				last = idx
			}
			idx--
		}
		if last == -1 {
			fmt.Fprintln(writer, 0)
			return
		}
		if i == 0 {
			if idx != -1 {
				fmt.Fprintln(writer, 0)
				return
			}
		} else {
			ans = ans * int64(last-idx) % mod
		}
	}
	fmt.Fprintln(writer, ans%mod)
}
`

const testcasesRaw = `2 2 4 13 4 5
1 1 1 4
5 3 2 8 17 18 12 3 4 4
3 1 1 9 9 2
2 2 10 12 1 5
3 3 13 17 8 2 3 6
3 1 18 10 1 3
5 3 17 7 14 14 20 3 6 9
2 1 10 9 1
1 1 15 3
5 5 16 11 5 7 3 4 5 8 10 11
3 2 19 11 18 2 4
1 1 8 3
5 5 8 4 11 6 10 4 4 4 6 6
3 3 11 1 11 3 5 6
6 4 20 3 10 20 7 15 3 4 6 9
5 2 11 19 1 12 2 4 5
3 2 10 19 4 4 5
4 2 4 2 2 2 2 6
6 2 20 2 18 16 19 8 3 3
1 1 14 2
4 2 8 15 14 16 1 2
4 4 8 14 7 16 2 2 2 4
3 1 17 7 8 4
3 1 11 2 11 5
1 1 2 4
4 1 14 7 19 6 3
3 3 16 11 14 5 6 8
3 2 16 3 9 2 2
4 2 9 2 6 15 5 8
6 4 13 7 1 7 6 1 5 7 7 10
4 2 18 2 7 6 5 7
5 4 17 15 1 3 2 5 5 8 12
3 3 5 2 12 1 5 5
3 2 3 3 18 4 7
2 2 13 8 4 7
1 1 4 5
3 3 14 14 15 1 2 4
4 4 4 18 6 12 2 3 4 6
4 3 9 18 1 6 1 3 3
5 1 16 20 16 17 3 5
2 2 10 12 2 3
6 1 2 20 11 18 15 19 3
5 4 20 20 15 13 5 3 7 9 11
2 2 3 20 2 6
2 2 12 7 5 7
6 5 3 3 13 6 11 12 3 4 6 6 10
1 1 12 1
2 1 19 16 5
1 1 6 4
6 6 8 20 10 13 20 8 4 5 7 9 10 12
5 5 15 13 17 13 11 3 6 9 13 13
3 1 18 15 18 5
3 2 13 20 1 2 6
1 1 12 5
3 1 9 16 8 4
5 1 5 8 3 10 5 1
2 2 19 20 5 9
3 1 17 8 6 1
2 1 11 3 1
3 1 10 20 6 2
6 4 5 3 18 12 1 18 2 6 9 10
2 2 16 17 1 4
2 1 9 17 4
5 3 13 11 6 13 2 4 4 6
6 1 10 5 3 6 4 20 1
2 1 18 1 4
6 5 6 15 13 11 6 17 5 6 6 9 13
3 2 14 20 2 2 4
5 5 10 16 11 3 8 3 3 3 7 9
5 4 12 3 6 2 16 5 9 13 14
1 1 3 3
5 2 10 4 17 17 1 1 2
5 3 20 1 2 1 16 2 3 5
2 2 10 13 5 8
1 1 14 4
1 1 18 3
2 2 11 6 2 2
4 4 9 11 14 14 3 7 8 10
3 3 3 1 13 3 7 9
6 6 5 20 11 3 6 20 1 2 4 5 8 8
6 3 1 15 19 7 6 3 2 6 6
1 1 10 5
5 2 3 6 1 6 7 3 3
5 3 11 20 12 6 6 2 2 2
2 1 9 16 1
2 1 13 16 5
5 4 18 11 12 14 13 1 1 2 3
3 1 10 8 12 5
2 2 13 12 5 5
2 1 12 6 5
2 1 1 4 2
3 3 11 18 16 1 3 6
2 1 13 20 4
3 2 1 15 5 4 4
6 6 14 7 4 11 10 20 5 5 5 8 9 13
1 1 20 2
5 3 20 1 1 13 14 2 6 7
5 1 5 6 20 10 8 4
6 3 18 12 6 12 18 19 3 5 9
1 1 13 3
5 1 19 19 15 9 20 15
6 3 16 8 1 16 4 5 1 4 6
1 1 11 3
4 1 5 2 14 5 4
5 1 6 17 16 18 14 2
5 1 4 1 2 4 8 3
3 1 1 1 1 2
6 1 12 5 2 8 6 7 6
4 1 6 6 15 6 5
4 1 5 3 19 1 3
6 1 9 6 16 7 18 5 5
3 1 18 16 16 5
2 1 6 5 5
4 1 8 12 9 5 5
3 1 11 12 8 8
3 1 5 6 11 5
4 1 9 11 13 11 13
3 1 10 5 20 16
2 1 6 5 5
4 1 6 7 6 12 13
6 1 5 1 2 2 5 6 4
1 1 1 1
3 1 14 8 17 15
6 1 5 20 1 1 14 6 5
1 1 2 2
3 1 12 13 4 13
1 1 16 7
2 1 2 1 1
2 1 2 2 1
4 1 13 1 2 10 9
5 1 12 1 9 7 12 16
2 1 2 2 1
1 1 18 7
3 1 17 18 1 7
3 1 12 8 13 8
3 1 10 13 7 16
3 1 3 4 9 2
5 1 3 5 5 6 6 7
5 1 1 7 7 9 9 10
4 1 11 13 15 16 16
3 1 9 12 15 4
6 1 15 16 17 18 18 19 16
5 1 16 17 18 19 20 20
4 1 13 15 17 18 20
6 1 15 1 4 6 7 8 9
5 1 10 11 14 14 15 16
4 1 7 9 9 11 11
6 1 5 10 11 12 13 16 17
2 1 12 15 15
5 1 5 6 7 8 9 14
6 1 6 10 12 13 14 20 20
4 1 5 13 13 14 15
4 1 11 12 15 16 16
4 1 16 18 19 20 20
5 1 11 12 13 14 14 15
3 1 7 12 14 20
4 1 7 11 17 19 19
3 1 7 9 14 15
4 1 1 3 15 18 18
5 1 3 5 16 19 20 20
5 1 6 8 12 17 18 20
2 1 1 10 10
3 1 4 9 11 20
3 1 4 8 8 8
4 1 9 12 14 18 19
3 1 14 15 17 17
3 1 16 17 18 18
5 1 1 2 8 10 12 16
5 1 3 9 11 17 18 19
3 1 6 7 8 12
5 1 4 5 9 11 12 17
1 1 1 1`

var _ = solutionSource

const mod int64 = 998244353

type testCase struct {
	n int
	m int
	a []int
	b []int
}

func solveCase(n, m int, a, b []int) int64 {
	idx := n - 1
	ans := int64(1)
	for i := m - 1; i >= 0; i-- {
		last := -1
		for idx >= 0 && a[idx] >= b[i] {
			if a[idx] == b[i] && last == -1 {
				last = idx
			}
			idx--
		}
		if last == -1 {
			return 0
		}
		if i == 0 {
			if idx != -1 {
				return 0
			}
		} else {
			ans = ans * int64(last-idx) % mod
		}
	}
	return ans % mod
}

func parseTests() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	tests := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid test %d", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad n on line %d: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad m on line %d: %v", idx+1, err)
		}
		if len(fields) != 2+n+m {
			return nil, fmt.Errorf("line %d wrong token count", idx+1)
		}
		a := make([]int, n)
		b := make([]int, m)
		pos := 2
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[pos])
			a[i] = v
			pos++
		}
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(fields[pos])
			b[i] = v
			pos++
		}
		tests = append(tests, testCase{n: n, m: m, a: a, b: b})
	}
	return tests, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		expected := solveCase(tc.n, tc.m, tc.a, tc.b)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", tc.a[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < tc.m; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", tc.b[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx+1, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
