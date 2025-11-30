package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `100
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
1 4 1 4
2 3 4 4
2 2 3 3
3 4 4 4
4 4 3 3
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 2
2 1
1 2 2 2
1 2 2 2
4 2
4 3 4 2
4 4 1 1
2 2 4 4
3 5
3 3 2
1 2 1 2
2 3 2 3
3 3 1 1
3 3 3 3
2 2 3 3
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 3
3 3 3
3 3 2 3
3 3 3 3
3 3 2 3
5 4
2 5 1 2 4
2 3 1 4
1 3 5 5
2 4 1 5
1 2 5 5
3 3
2 2 3
2 3 1 1
2 3 2 3
2 3 1 1
3 1
1 3 3
2 2 1 2
1 2
1
1 1 1 1
1 1 1 1
2 3
1 1
1 2 2 2
1 2 2 2
2 2 1 2
3 1
1 2 1
1 1 3 3
2 1
2 2
2 2 2 2
1 2
1
1 1 1 1
1 1 1 1
4 5
4 3 4 3
3 3 4 4
1 1 2 3
1 2 3 4
1 1 3 4
2 4 1 3
3 3
3 1 3
3 3 3 3
2 2 3 3
3 3 3 3
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
5 5
2 2 5 1 4
4 5 4 5
1 3 5 5
1 5 4 5
5 5 1 5
4 5 1 5
2 2
2 2
2 2 1 2
1 1 1 2
2 4
1 2
1 2 1 2
2 2 1 2
2 2 2 2
1 1 2 2
2 2
2 2
1 2 1 2
2 2 2 2
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 2
2 1 4 1
4 4 3 4
1 3 3 4
5 3
4 1 2 4 4
2 3 4 4
4 5 3 5
1 1 5 5
2 2
2 1
1 1 1 1
1 2 1 2
5 4
5 2 1 5 3
2 5 2 2
1 4 1 2
1 1 1 5
4 5 2 4
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 5
3 2 3 4
2 3 3 4
4 4 2 3
3 3 1 4
2 2 1 4
3 3 1 3
4 3
2 1 2 1
2 2 4 4
3 3 2 2
3 3 4 4
3 1
2 3 3
2 3 3 3
1 2
1
1 1 1 1
1 1 1 1
4 2
3 4 2 4
3 4 2 4
2 3 4 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 2
3 1 2
1 2 3 3
2 2 1 1
5 2
4 3 1 3 5
3 3 4 4
4 4 2 4
1 1
1
1 1 1 1
2 4
1 1
2 2 2 2
2 2 1 1
2 2 1 2
2 2 1 2
5 2
3 2 5 1 1
4 4 2 3
4 5 3 4
3 4
1 3 1
1 1 1 3
3 3 2 2
2 2 1 2
3 3 2 2
1 2
1
1 1 1 1
1 1 1 1
3 5
2 1 3
3 3 3 3
3 3 3 3
3 3 2 3
2 3 1 2
2 2 1 2
4 1
4 4 4 4
4 4 4 4
5 4
3 5 3 3 4
1 1 5 5
4 5 2 5
2 2 1 4
2 3 2 5
2 4
1 2
1 2 1 1
2 2 1 2
2 2 2 2
1 1 1 1
5 5
2 3 2 1 3
2 2 2 2
3 4 4 5
2 3 5 5
3 3 5 5
3 4 5 5
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
4 3 3 2
4 4 1 1
4 4 3 3
2 2 1 2
1 2 2 3
1 2
1
1 1 1 1
1 1 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
4 4
4 3 2 2
3 4 2 4
3 3 1 2
2 3 1 4
4 4 1 2
4 3
1 1 3 2
1 3 2 3
3 3 3 3
3 4 4 4
4 4
2 3 4 4
4 4 1 4
1 3 3 4
3 4 2 2
4 4 3 3
5 2
2 2 2 5 4
1 1 4 4
3 4 3 4
1 2
1
1 1 1 1
1 1 1 1
5 3
2 5 2 3 3
5 5 2 3
4 4 3 5
1 2 5 5
2 2
2 2
2 2 1 1
2 2 2 2
1 2
1
1 1 1 1
1 1 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 3
1 1 3
2 2 2 2
2 3 3 3
3 3 3 3
4 3
4 1 1 3
3 3 2 3
3 4 2 4
4 4 2 4
2 5
1 1
2 2 2 2
2 2 1 2
1 1 1 1
2 2 1 2
1 2 1 1
5 5
4 5 3 2 2
3 4 2 4
3 4 2 3
2 3 5 5
4 5 1 2
5 5 1 2
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 5
2 1
1 2 2 2
2 2 2 2
1 2 2 2
1 1 2 2
2 2 1 1
5 1
1 1 1 2 5
5 5 4 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 2
3 1 1
2 2 3 3
1 2 1 3
5 1
3 3 3 4 2
2 4 3 5
5 5
3 1 4 1 1
2 3 4 4
4 5 5 5
2 4 4 4
5 5 4 5
1 2 4 5
4 1
1 2 3 3
2 3 1 1
2 2
2 2
2 2 1 2
1 1 2 2
4 3
3 1 2 4
2 2 2 3
2 3 4 4
3 4 1 1
1 4
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
2 3
1 1
1 1 1 2
1 2 1 2
2 2 1 1
2 4
2 2
1 1 2 2
2 2 1 1
2 2 1 1
1 2 1 2
3 4
1 3 3
2 3 1 2
2 3 2 3
3 3 2 3
3 3 3 3
4 4
4 2 2 2
2 2 2 3
1 2 1 2
1 1 3 4
1 3 4 4
5 5
2 1 3 5 4
5 5 1 3
5 5 2 3
5 5 3 5
1 1 2 2
3 4 3 4
2 1
1 1
2 2 2 2
2 4
2 1
2 2 2 2
1 1 2 2
2 2 2 2
2 2 1 2
4 5
2 4 3 2
2 3 1 2
3 3 2 2
4 4 3 3
1 3 4 4
3 3 2 4
2 4
1 1
1 2 2 2
1 2 2 2
1 2 2 2
1 1 2 2
5 2
1 1 5 5 3
3 3 2 2
2 5 1 5
1 1
1
1 1 1 1
4 2
2 2 4 2
3 3 4 4
3 4 1 2
4 3
4 3 4 2
1 2 3 4
1 1 2 2
4 4 2 2
4 5
3 2 2 4
4 4 2 3
1 1 2 2
1 4 1 4
2 4 1 1
2 2 3 4
1 5
1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
1 1 1 1
3 1
1 3 2
2 2 1 3
3 3
3 2 3
3 3 2 3
2 3 2 2
3 3 2 2
1 3
1
1 1 1 1
1 1 1 1
1 1 1 1
2 2
1 1
2 2 2 2
1 2 1 2
3 1
1 2 3
3 3 2 2
5 1
5 3 4 1 4
1 2 4 4
2 5
1 2
1 2 2 2
2 2 1 2
2 2 1 2
2 2 1 2
1 2 1 2
2 3
1 1
2 2 1 2
2 2 1 2
2 2 2 2
4 3
1 2 3 2
1 1 1 1
1 3 3 4
2 3 3 3`

type testCase struct {
	input    string
	expected string
}

func g(arr []int, i, j int) int {
	if i > j {
		return 0
	}
	required := make(map[int]struct{})
	for p := i; p <= j; p++ {
		required[arr[p-1]] = struct{}{}
	}
	for x := j; x >= 1; x-- {
		if _, ok := required[arr[x-1]]; ok {
			delete(required, arr[x-1])
			if len(required) == 0 {
				return x
			}
		}
	}
	return 0
}

func solveCase(n, q int, arr []int, queries [][4]int) string {
	var sb strings.Builder
	for qi, qu := range queries {
		l, r, x, y := qu[0], qu[1], qu[2], qu[3]
		ans := 0
		for i := l; i <= r; i++ {
			for j := x; j <= y; j++ {
				if i <= j {
					ans += g(arr, i, j)
				}
			}
		}
		sb.WriteString(strconv.Itoa(ans))
		if qi+1 < len(queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if pos+1 >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n/q", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseIdx+1, err)
		}
		pos++
		q, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: bad q: %w", caseIdx+1, err)
		}
		pos++
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if pos >= len(fields) {
				return nil, fmt.Errorf("case %d: missing array value", caseIdx+1)
			}
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d: bad array value: %w", caseIdx+1, err)
			}
			arr[i] = v
			pos++
		}
		queries := make([][4]int, q)
		for i := 0; i < q; i++ {
			if pos+3 >= len(fields) {
				return nil, fmt.Errorf("case %d: missing query", caseIdx+1)
			}
			l, _ := strconv.Atoi(fields[pos])
			r, _ := strconv.Atoi(fields[pos+1])
			x, _ := strconv.Atoi(fields[pos+2])
			y, _ := strconv.Atoi(fields[pos+3])
			pos += 4
			queries[i] = [4]int{l, r, x, y}
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range queries {
			fmt.Fprintf(&sb, "%d %d %d %d\n", qu[0], qu[1], qu[2], qu[3])
		}
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: solveCase(n, q, arr, queries),
		})
	}
	return cases, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runExe(os.Args[1], tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
