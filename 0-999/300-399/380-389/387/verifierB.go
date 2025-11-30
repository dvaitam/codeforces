package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// referenceSolutionSource embeds the original 387B solution for traceability.
const referenceSolutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }
   i, j, matched := 0, 0, 0
   for i < n && j < m {
       if b[j] >= a[i] {
           matched++
           i++
           j++
       } else {
           j++
       }
   }
   // need to come up with problems for unmatched requirements
   fmt.Println(n - matched)
}
`

var _ = referenceSolutionSource

const rawTestcases = `100
1 1 9 15
1 7 9 4 5 12 17 21 23 24
3 3 2 7 19 9 18 22
3 5 2 16 19 12 14 15 21 23
5 5 1 3 7 10 13 5 9 11 12 23
1 3 2 2 6 9
2 6 10 12 4 5 10 13 16 18
2 2 6 10 17 24
1 3 13 10 11 14
1 1 18 16
4 6 4 11 16 18 2 10 11 14 16 23
2 3 13 19 3 7 21
2 2 1 13 4 13
5 7 7 10 15 16 17 3 6 8 9 12 14 19
2 4 3 4 1 15 17 23
2 2 13 16 7 9
1 7 7 4 5 7 12 13 15 20
5 5 4 5 11 13 16 11 14 16 17 22
4 7 1 7 8 18 2 5 11 17 22 23 24
3 7 5 10 13 2 3 16 17 21 23 24
1 2 5 2 10
1 7 15 5 6 11 12 15 17 21
4 4 3 17 18 19 7 10 14 24
5 7 8 13 14 16 19 1 6 9 10 17 22 24
3 3 9 10 16 2 13 14
2 7 5 8 2 5 10 11 14 16 22
4 4 2 5 12 14 2 13 15 20
1 4 5 1 2 5 20
3 3 7 12 18 4 13 16
1 5 15 4 11 20 21 22
5 6 4 5 9 10 13 2 7 12 13 15 24
4 6 2 3 16 18 1 7 9 17 19 21
2 2 17 19 14 17
3 3 5 14 18 3 4 14
1 1 14 5
1 7 15 1 9 11 14 16 22 24
1 3 3 4 12 23
1 3 12 1 6 8
3 3 1 5 7 7 22 24
1 6 1 1 5 8 10 12 20
2 5 4 16 1 5 9 12 23
2 4 11 16 10 11 18 21
2 6 3 4 5 6 10 13 18 19
2 3 11 17 6 8 10
3 6 1 2 5 3 5 10 13 14 23
5 6 2 5 10 12 14 2 8 12 15 17 21
4 7 1 11 14 15 3 4 6 7 10 12 16
3 3 5 15 18 6 13 14
4 5 5 8 11 15 3 12 15 16 21
2 4 1 15 1 7 15 20
3 3 5 10 18 14 16 23
1 6 16 1 4 8 9 13 18
3 3 1 9 13 17 19 23
4 7 4 9 10 12 2 3 7 9 20 21 22
3 7 4 11 17 3 6 8 10 14 17 20
2 6 7 17 4 13 14 18 21 24
3 5 5 12 15 4 6 13 19 23
4 5 10 12 16 18 7 14 16 23 24
5 6 2 3 10 15 16 1 2 7 16 20 24
3 6 1 13 17 1 3 12 13 22 24
1 1 1 9
3 4 5 7 10 4 14 15 23
3 6 6 11 14 5 14 15 21 22 24
5 6 5 6 7 15 17 7 8 13 14 16 24
4 5 2 13 18 19 3 6 8 12 21
1 6 6 3 8 10 17 20 23
3 5 2 14 15 17 18 21 22 23
4 7 7 9 16 19 2 6 9 11 12 21 22
1 3 1 3 5 14
2 6 13 18 4 7 8 11 15 20
5 5 5 11 15 18 19 1 2 7 12 17
1 2 17 7 12
2 4 10 19 9 13 16 17
3 4 2 3 10 1 15 16 24
4 4 14 15 16 17 3 4 5 8
4 5 3 13 14 15 2 5 6 8 16
3 5 4 11 14 7 10 18 20 23
3 6 9 15 17 1 4 8 9 20 21
2 7 8 14 1 7 10 14 17 18 22
1 1 13 21
3 3 8 12 19 18 22 23
3 4 3 8 17 8 10 11 22
3 6 5 6 10 1 11 12 17 18 19
1 7 5 3 5 6 7 13 17 23
4 5 5 8 13 19 5 12 19 20 21
4 4 1 12 16 17 1 8 10 15
5 7 6 11 12 16 18 3 5 7 9 11 13 20
3 6 2 7 19 8 9 11 15 23 24
3 4 1 10 12 2 18 19 24
2 4 1 16 1 2 8 21
1 1 8 21
3 3 2 12 14 5 7 15
4 5 6 10 11 12 1 9 13 14 24
5 7 2 4 7 14 15 1 5 6 17 20 21 22
1 3 8 1 6 8
2 7 6 18 3 4 5 14 15 20 23
5 7 1 2 9 11 13 2 3 5 10 12 16 21
4 5 6 12 13 17 1 9 11 13 16
3 7 10 16 18 2 9 15 18 19 23 24
4 4 2 12 13 16 1 2 9 24`

var testcases = mustParseTestcases(rawTestcases)

type testcase struct {
	n int
	m int
	a []int
	b []int
}

func mustParseTestcases(data string) []testcase {
	r := strings.NewReader(strings.TrimSpace(data))
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		panic(err)
	}
	out := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(r, &n, &m); err != nil {
			panic(err)
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(r, &a[j])
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(r, &b[j])
		}
		out = append(out, testcase{n: n, m: m, a: a, b: b})
	}
	return out
}

func solve(a, b []int) int {
	n, m := len(a), len(b)
	i, j, matched := 0, 0, 0
	for i < n && j < m {
		if b[j] >= a[i] {
			matched++
			i++
			j++
		} else {
			j++
		}
	}
	return n - matched
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

	for idx, tc := range testcases {
		expected := strconv.Itoa(solve(tc.a, tc.b))

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
