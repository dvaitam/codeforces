package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	m    int
	rows []string
}

// solve is the embedded logic from 16A.go.
func solve(n, m int, rows []string) string {
	prev := byte(0)
	ok := true
	for i := 0; i < n; i++ {
		if len(rows[i]) != m {
			ok = false
			break
		}
		c := rows[i][0]
		for j := 1; j < m; j++ {
			if rows[i][j] != c {
				ok = false
				break
			}
		}
		if !ok {
			break
		}
		if i > 0 && c == prev {
			ok = false
			break
		}
		prev = c
	}
	if ok {
		return "YES"
	}
	return "NO"
}

// Embedded copy of testcasesA.txt (one case per line: n m rows...).
const testcaseData = `
4 4 0487 6475 9382 4219
3 5 92411 57815 65938
4 4 8408 0160 9753 5139
2 2 28 71
1 3 871
3 5 41858 39894 71965
5 2 42 32 09 47 11
2 2 01 86
5 3 833 969 477 515 917
5 3 330 413 525 601 230
5 5 91013 99161 51090 32173 00869
1 3 131
3 3 620 870 916
2 3 579 230
2 2 58 41
5 4 2076 9845 6428 0715 0842
2 4 5945 9924
4 4 1093 5233 7696 0696
1 2 71
3 2 78 78 90
1 4 5470
4 2 81 20 66 50
2 1 0 8
5 1 3 1 9 3 4
3 2 17 61 04
4 1 4 2 8 5
1 2 40
1 1 3
3 5 55909 77658 23694
1 2 24
3 3 515 900 422
5 3 568 241 730 428 146
3 3 611 877 551
4 1 7 6 0 4
3 2 29 61 11
2 2 06 01
4 5 84779 36153 49263 51108
4 2 17 64 30 39
2 1 3 7
4 3 821 972 966 875
4 4 3893 0555 0824 9269
3 4 1180 1320 4075
2 2 75 86
5 5 09189 16348 96769 93002 48945
1 4 4466
4 1 2 2 3 4
3 1 0 7 6
2 4 9125 6097
4 4 0172 0099 2518 5367
1 1 9
4 5 51942 64183 06753 75100
4 3 089 933 188 684
1 2 69
4 1 1 6 1 1
4 2 07 66 07 54
1 3 115
1 3 520
2 3 192 303
1 1 4
3 1 9 3 2
2 4 1754 2035
3 4 4485 2911 8942
4 2 23 58 33 24
3 4 0290 6112 6486
2 5 64513 75806
4 1 6 5 7 3
3 3 712 141 892
4 4 2662 3758 2571 7340
4 5 70341 48926 71738 64014
1 1 4
4 5 96714 54391 01448 51832
1 4 4482
5 5 38168 64475 92211 66972 84576
2 4 7857 0742
4 1 9 3 0 5
4 4 0811 6050 1904 4329
3 2 16 75 62
3 4 6272 8523 2756
4 4 6337 3906 0312 5023
5 3 918 456 708 869 774
4 2 54 00 02 50
3 1 2 1 6
2 5 68373 59191
3 3 875 408 035
1 2 85
2 2 44 48
4 3 753 048 107 770
4 4 7711 1312 6379 1686
1 2 37
2 2 45 56
1 5 49834
4 5 97844 30191 26334 08860
1 4 4195
2 5 43318 45357
3 5 22088 55902 62281
2 2 79 33
2 2 65 99
2 4 1908 9577
3 1 3 8 2
4 4 8514 2963 5460 3053
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad m: %v", i+1, err)
		}
		rows := fields[2:]
		if len(rows) != n {
			return nil, fmt.Errorf("line %d: expected %d rows got %d", i+1, n, len(rows))
		}
		tests = append(tests, testCase{n: n, m: m, rows: rows})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, r := range tc.rows {
		input.WriteString(r)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc.n, tc.m, tc.rows)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
