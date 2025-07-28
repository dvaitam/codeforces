package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Course struct {
	x int64
	y int64
}

type Query struct {
	l int
	r int
}

type TestCase struct {
	n       int
	courses []Course
	q       int
	queries []Query
}

func genTests() []TestCase {
	rand.Seed(time.Now().UnixNano())
	const T = 100
	tests := make([]TestCase, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(6) + 1
		courses := make([]Course, n)
		for j := range courses {
			x := rand.Intn(16)
			y := x + rand.Intn(16-x)
			courses[j] = Course{int64(x), int64(y)}
		}
		q := rand.Intn(6) + 1
		queries := make([]Query, q)
		for j := range queries {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			queries[j] = Query{l, r}
		}
		tests[i] = TestCase{n: n, courses: courses, q: q, queries: queries}
	}
	return tests
}

func buildInput(tests []TestCase) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&buf, tc.n)
		for _, c := range tc.courses {
			fmt.Fprintf(&buf, "%d %d\n", c.x, c.y)
		}
		fmt.Fprintln(&buf, tc.q)
		for _, qu := range tc.queries {
			fmt.Fprintf(&buf, "%d %d\n", qu.l, qu.r)
		}
	}
	return buf.Bytes()
}

func runBinary(path string, input []byte) ([]byte, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	return cmd.CombinedOutput()
}

func canSet(x, y int64, b uint) bool {
	start := x
	if ((start >> b) & 1) == 0 {
		start = (start>>(b+1))<<(b+1) + (1 << b)
	}
	return start <= y
}

func expected(tc TestCase) string {
	res := make([]int, tc.q)
	for idx, q := range tc.queries {
		val := 0
		for b := 0; b < 30; b++ {
			ok := false
			for i := q.l - 1; i <= q.r-1; i++ {
				if canSet(tc.courses[i].x, tc.courses[i].y, uint(b)) {
					ok = true
					break
				}
			}
			if ok {
				val |= 1 << b
			}
		}
		res[idx] = val
	}
	parts := make([]string, len(res))
	for i, v := range res {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, "\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests := genTests()
	input := buildInput(tests)

	expectedOutputs := make([]string, len(tests))
	for i, tc := range tests {
		expectedOutputs[i] = expected(tc)
	}

	out, err := runBinary(binary, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}

	gotLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// Flatten expected outputs lines
	var expLines []string
	for _, exp := range expectedOutputs {
		expLines = append(expLines, strings.Split(exp, "\n")...)
	}
	if len(gotLines) != len(expLines) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expLines), len(gotLines))
		os.Exit(1)
	}
	for i := range expLines {
		if strings.TrimSpace(gotLines[i]) != strings.TrimSpace(expLines[i]) {
			fmt.Fprintf(os.Stderr, "mismatch on output line %d expected %s got %s\n", i+1, expLines[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
