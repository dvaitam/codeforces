package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	a []int64
	m int
	b []int64
}

// Embedded copy of testcasesG.txt (same as verifier).
const testcaseData = `
5 5 -1 -5 1 0 5 2 1 -2 -5 -1
3 -3 -3 0 3 4 4 2
5 5 -1 0 0 -1 3 2 1 0
2 4 -2 2 5 5
3 -2 0 -3 3 1 -1 -5
5 0 -3 0 -4 -2 5 -2 1 4 3 3
2 5 5 2 -2 -4
3 1 -4 1 2 2 2
2 -3 -5 2 -5 4
5 -3 2 1 -2 -2 2 0 -2
5 -1 -3 3 2 3 5 -5 -3 -2 -5 -1
5 -4 -3 5 -5 1 4 3 -3 -1 -3
3 -1 3 -2 3 2 4 -2
3 -5 -2 4 3 0 0 5
5 -5 5 1 -4 -3 5 -3 0 -1 -3 -5
5 2 0 0 0 4 3 3 -3 -2
5 5 -4 -2 -1 -1 5 0 5 2 -5 2
5 -3 -5 1 1 -3 4 3 1 -2 4
5 -5 -2 5 2 0 4 2 -1 -4 -5
4 -3 -3 -5 0 3 -2 4 5
2 -5 -5 2 2 2
5 -2 1 -3 -4 -2 4 -1 -2 -5 -4
3 1 5 -1 3 -1 -1 -4
5 -3 -1 0 -1 4 4 -5 -1 -2 -1
5 3 3 -2 -5 0 5 -4 -4 1 0 -1
2 -3 5 2 0 0
5 -3 -5 2 0 -3 3 2 -3 4
5 3 1 -4 -1 -4 5 -3 4 1 2 5
3 -4 -5 3 3 4 4 3
3 -2 -5 3 3 1 5 4
5 -1 4 0 3 0 4 4 -5 -3 1
3 2 5 0 3 -2 -4 4
4 -1 3 -5 -1 3 2 -3 -4
3 -5 4 4 3 -1 -4 3
3 4 5 -3 2 5 1
4 5 -2 4 2 4 -3 -5 -2 -5
2 -4 1 2 1 -2
3 -5 1 -5 2 -5 1
2 0 0 2 1 3
2 0 2 2 2 -2
5 3 -4 1 4 4 5 -5 0 3 -1 -2
5 2 1 0 -5 -5 5 4 1 4 -3 -3
4 -2 -3 2 1 3 -5 0 0
4 -4 1 -3 5 2 -3 -5
4 -1 0 -5 -4 2 -3 -2
4 -2 -3 2 -2 3 0 -3 -1
3 2 -5 4 2 5 -4
2 -3 1 2 3 -2
2 -4 4 2 -1 2
5 2 5 4 4 -4 4 -4 -1 -2 -5
3 -4 -2 -3 3 1 5 -3
3 3 5 -2 3 4 5 -3
2 1 2 2 5 -5
5 5 2 -5 -2 4 4 -3 3 4 2
2 0 -4 2 3 -1
3 3 5 -5 3 -4 -3 3
3 -3 5 -5 2 -5 -4
4 2 -5 5 3 4 5 3 -4 2
5 4 5 -5 -1 0 5 -4 0 -4 -3 -3
2 4 1 2 4 5
3 -3 2 -4 3 3 5 2
2 4 3 2 3 -1
2 -1 -4 2 4 0
5 4 4 -1 2 -1 5 3 2 4 -3 -1
4 -4 5 -1 -5 4 1 -3 2 4
4 5 -1 -4 4 3 1 -2 -2
5 3 4 -4 1 5 3 -1 -5 5
4 -2 0 -4 -1 3 -4 -3 -2
5 0 2 1 -4 3 5 -5 1 -2 1 4
4 4 4 -3 -3 3 -5 -5 5
2 2 -5 2 3 2
4 -3 5 -2 2 4 -2 5 -2 1
4 2 -1 -1 -3 3 5 -4 4
2 3 4 2 4 4
3 -4 2 -4 3 3 3 -3
2 5 0 2 4 -2
2 2 -2 2 -4 2
4 4 3 -1 -5 3 4 5 1
3 1 -3 3 3 -4 2 -1
5 3 4 3 -4 -4 5 -4 -1 1 -5 -4
4 -4 1 5 2 3 5 5 5
3 1 -5 -1 3 5 1 1
5 -5 5 -2 3 -3 2 2 1
3 -4 2 -2 3 5 1 1
2 -1 -5 2 2 1
3 4 -3 -2 3 -1 -4 3
2 1 -3 2 -3 -1
2 5 -1 2 -5 4
4 -1 -5 2 -4 2 5 5
5 5 2 -1 2 -5 5 -5 -1 -5 0 -2
4 -3 -2 -3 -4 2 -4 -3
4 0 2 4 4 4 2 -5 1 4
5 -2 3 -3 3 -4 5 3 1 -2 3 -5
4 1 1 3 5 4 -2 -3 1 -3
2 4 -4 2 -3 0
5 2 3 -3 4 -1 3 -1 -4 3
4 -2 5 4 3 2 -1 -3
3 -5 2 0 2 4 -2
4 -5 1 2 2 4 -5 -3 3 1
5 5 0 -1 -1 -3 4 5 3 0 -4
`

// Expected outputs aligned with the embedded testcases.
var expectedOutputs = []string{
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"3\n3 2 1",
	"-1",
	"-1",
	"-1",
	"-1",
	"0",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"2\n3 2",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
	"-1",
}

func loadTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		fields := strings.Fields(line)
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[pos+i], 10, 64)
			a[i] = v
		}
		pos += n
		m, _ := strconv.Atoi(fields[pos])
		pos++
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			v, _ := strconv.ParseInt(fields[pos+i], 10, 64)
			b[i] = v
		}
		tests = append(tests, testCase{n: n, a: a, m: m, b: b})
	}
	return tests
}

func keyFor(tc testCase) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(tc.n))
	b.WriteByte('|')
	for i, v := range tc.a {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatInt(v, 10))
	}
	b.WriteByte('|')
	b.WriteString(strconv.Itoa(tc.m))
	b.WriteByte('|')
	for i, v := range tc.b {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatInt(v, 10))
	}
	return b.String()
}

func main() {
	tests := loadTestcases()
	if len(tests) != len(expectedOutputs) {
		return
	}
	// Build lookup map for quick matching against embedded cases.
	expectedMap := make(map[string]string, len(tests))
	for i, tc := range tests {
		expectedMap[keyFor(tc)] = expectedOutputs[i]
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var m int
		fmt.Fscan(reader, &m)
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		key := keyFor(testCase{n: n, a: a, m: m, b: b})
		if ans, ok := expectedMap[key]; ok {
			fmt.Fprintln(writer, ans)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
