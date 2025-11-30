package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type query struct {
	t int
	a int64
}

type testCase struct {
	n       int
	q       int
	arr     []int64
	queries []query
}

// Embedded testcases from testcasesE.txt.
const testcaseData = `
1 2 3 1 9 2 7
2 3 7 8 2 9 1 3 2 3
1 1 8 2 5
3 4 1 10 5 2 8 1 1 2 4 2 5
2 4 2 7 1 1 1 6 1 5 2 1
1 1 10 1 3
2 4 2 10 1 1 1 2 1 2 1 4
1 4 7 1 2 2 8 1 9 1 3
3 4 9 3 2 1 9 1 2 2 6 2 9
4 3 4 10 2 3 2 6 2 6 2 1
1 1 3 2 4
2 2 2 8 2 3 2 5
2 1 3 6 2 4
1 4 3 2 6 1 2 1 6 2 2
3 4 9 6 1 2 4 2 1 2 9 2 3
3 1 5 2 6 2 1
3 2 3 5 6 1 3 2 4
1 4 5 1 5 1 1 1 6 1 7
2 2 1 9 1 8 2 10
3 3 6 2 1 2 6 1 8 1 7
1 2 10 2 7 2 3
2 4 7 8 1 5 1 3 1 5 1 2
2 3 3 7 2 9 2 6 2 8
4 1 5 10 6 5 2 5
4 1 8 9 9 4 2 5
4 1 1 7 9 2 2 1
1 4 7 2 9 2 9 2 9 2 10
2 1 10 6 1 5
4 1 4 2 3 10 2 3
4 3 8 3 8 5 2 6 1 1 2 9
3 4 9 10 2 2 3 1 3 1 6 1 1
1 4 8 1 10 1 3 2 10 1 10
3 1 5 6 3 2 6
1 3 2 1 10 1 5 2 10
3 2 1 8 7 2 4 1 7
3 4 7 8 7 1 7 1 7 1 5 1 3
4 3 6 5 8 6 1 3 2 3 2 4
4 2 2 9 10 10 2 6 1 2
3 2 5 6 2 1 7 2 6
3 3 1 4 5 1 3 1 8 1 2
3 1 1 7 2 2 8
4 2 3 1 2 3 1 7 1 5
1 1 4 2 7
2 4 5 6 1 7 2 3 1 2 1 9
3 4 2 6 9 2 5 2 1 1 6 1 8
2 2 2 3 2 7 1 6
2 1 6 3 1 4
3 4 6 10 8 2 5 1 10 1 5 1 9
1 3 8 1 2 1 8 2 7
3 3 4 6 8 1 9 1 8 1 1
1 1 3 1 1
3 3 4 2 1 1 3 2 7 2 7
2 4 1 10 1 10 1 8 1 7 1 4
1 2 6 2 5 2 6
4 4 6 5 9 5 2 9 2 7 2 7 2 1
2 4 5 7 2 3 1 8 1 3 1 2
3 2 10 3 8 1 8 1 7
2 1 1 9 1 2
1 2 3 1 6 2 2
1 4 8 1 6 1 10 1 10 2 10
3 1 3 7 2 1 3
2 3 5 6 1 8 2 9 1 2
4 1 1 4 6 7 1 6
4 1 6 2 8 10 2 5
2 4 1 6 2 3 1 9 2 8 1 10
2 2 9 10 2 8 1 7
3 4 5 1 7 1 8 1 3 2 6 2 10
2 4 8 1 2 3 2 1 2 3 1 2
4 4 9 3 7 6 1 9 2 10 1 2 1 4
1 2 8 1 1 2 6
3 3 10 5 8 1 5 2 7 1 8
1 4 6 1 10 2 5 1 3 2 4
1 3 7 1 5 2 10 1 7
1 1 3 2 7
1 1 10 1 9
1 2 6 2 9 1 3
2 1 8 5 1 9
1 4 2 1 1 1 3 2 4 2 10
1 2 6 1 10 2 5
2 3 1 10 2 6 1 6 2 6
1 3 10 1 6 1 1 2 1
4 2 7 10 8 5 1 3 2 9
1 1 1 1 3
4 4 10 3 4 9 1 1 1 7 2 5 1 7
1 3 7 1 9 2 10 1 1
4 4 9 1 4 6 2 3 1 4 1 4 1 5
3 3 5 3 3 2 5 1 7 1 10
4 2 7 4 7 7 1 9 2 10
3 1 10 9 4 1 9
4 4 4 10 2 9 1 3 1 5 1 6 1 3
2 4 5 1 2 3 1 9 1 8 2 7
1 2 10 2 4 2 2
3 4 6 5 1 2 2 2 5 1 6 1 4
2 2 7 2 2 3 2 7
1 1 9 2 2
3 2 9 5 2 1 9 2 9
3 2 3 2 5 2 3 2 5
2 2 2 8 1 6 2 5
3 1 3 8 6 1 2
4 4 8 6 10 5 2 3 2 1 2 2 2 10
`

func parseTestcases() ([]testCase, error) {
	data := strings.TrimSpace(testcaseData)
	if data == "" {
		return nil, fmt.Errorf("no test data")
	}
	lines := strings.Split(data, "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("case %d missing n/q", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		q, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("case %d bad q: %v", i+1, err)
		}
		expect := 2 + n + 2*q
		if len(fields) != expect {
			return nil, fmt.Errorf("case %d expected %d values, got %d", i+1, expect, len(fields))
		}
		arr := make([]int64, n)
		pos := 2
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(fields[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad arr value %d: %v", i+1, j+1, err)
			}
			arr[j] = v
			pos++
		}
		qs := make([]query, q)
		for j := 0; j < q; j++ {
			t, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("case %d bad query type %d: %v", i+1, j+1, err)
			}
			a, err := strconv.ParseInt(fields[pos+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("case %d bad query value %d: %v", i+1, j+1, err)
			}
			qs[j] = query{t: t, a: a}
			pos += 2
		}
		res = append(res, testCase{n: n, q: q, arr: arr, queries: qs})
	}
	return res, nil
}

// calcUnsuitable mirrors 1500E.go.
func calcUnsuitable(arr []int64) int64 {
	n := len(arr)
	if n <= 1 {
		return 0
	}
	pre := make([]int64, n+1)
	for i, v := range arr {
		pre[i+1] = pre[i] + v
	}
	total := pre[n]
	base := total - 2*arr[0]
	if n == 2 {
		if base < 0 {
			return 0
		}
		return base
	}
	gaps := int64(0)
	for k := 1; k <= n-2; k++ {
		gap := pre[k+1] + pre[n-k] - total
		if gap > 0 {
			gaps += gap
		}
	}
	res := base - gaps
	if res < 0 {
		return 0
	}
	return res
}

func solve(tc testCase) string {
	arr := append([]int64(nil), tc.arr...)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	answers := make([]string, 0, tc.q+1)
	answers = append(answers, strconv.FormatInt(calcUnsuitable(arr), 10))
	for _, qu := range tc.queries {
		if qu.t == 1 { // add
			arr = append(arr, qu.a)
			sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		} else { // remove
			idx := sort.Search(len(arr), func(i int) bool { return arr[i] >= qu.a })
			if idx < len(arr) && arr[idx] == qu.a {
				arr = append(arr[:idx], arr[idx+1:]...)
			}
		}
		answers = append(answers, strconv.FormatInt(calcUnsuitable(arr), 10))
	}
	return strings.Join(answers, "\n")
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.queries {
		fmt.Fprintf(&sb, "%d %d\n", qu.t, qu.a)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	output, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("runtime error: %v\n%s", err, string(ee.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
