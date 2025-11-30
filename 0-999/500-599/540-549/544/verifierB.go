package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testcase struct {
	n int
	k int
}

// Embedded testcases from testcasesB.txt.
const testcasesRaw = `1 0
1 1
2 4
3 9
2 2
2 3
8 60
7 13
2 3
1 1
7 38
1 1
5 23
4 3
6 1
1 0
9 1
7 43
4 13
1 0
8 63
9 29
6 14
4 14
5 0
7 35
2 1
5 3
6 32
7 32
4 9
5 18
8 64
7 37
1 1
4 12
7 42
3 5
9 47
2 3
9 13
3 8
7 23
8 3
8 5
5 22
10 75
10 50
3 2
9 29
1 0
9 70
4 12
9 44
10 45
8 34
9 77
1 1
9 16
9 71
4 13
1 1
6 36
9 25
9 52
8 45
7 22
1 1
8 3
4 5
9 74
3 1
9 32
1 0
2 0
8 1
5 7
5 3
10 23
6 18
2 1
3 4
9 21
5 20
5 14
6 31
8 14
1 1
7 21
7 12
5 3
5 23
9 26
10 55
1 0
1 1
3 0
3 7
9 54
9 28`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testcase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d malformed", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		k, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: invalid numbers", idx+1)
		}
		res = append(res, testcase{n: n, k: k})
	}
	return res, nil
}

// Embedded solver logic from 544B.go.
func solve(n, k int) string {
	max := (n*n + 1) / 2
	if max < k {
		return "NO"
	}
	var b strings.Builder
	b.WriteString("YES\n")
	z := 0
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if (z+j)%2 == 0 && k > 0 {
				row[j] = 'L'
				k--
			} else {
				row[j] = 'S'
			}
		}
		z++
		b.Write(row)
		if i+1 < n {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc.n, tc.k)
		input := fmt.Sprintf("%d %d\n", tc.n, tc.k)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, errb.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\nexpected:\n%s\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
