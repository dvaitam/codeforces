package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesData = `
3 1 2 3
9 1 1 2 4 2 0 4 0 7
1 1
7 1 1 3 3 5 6 0
19 1 1 1 4 5 2 5 0 6 7 8 2 0 1 11 11 0 2 6
4 1 0 2 0
11 1 0 2 3 5 1 2 6 4 8 0
6 0 0 3 0 5 6
17 0 2 3 1 2 6 1 1 8 2 4 3 12 4 10 8 8
17 1 0 3 4 1 0 2 8 0 5 1 3 10 12 14 7 14
17 0 2 2 1 3 6 0 8 1 8 5 0 13 12 2 3 13
20 1 2 3 2 3 4 5 1 2 5 0 2 11 2 0 10 6 1 13 20
2 1 1
2 0 1
3 1 0 3
12 1 2 2 4 3 3 2 0 7 4 3 6
3 1 0 0
1 1
1 0
13 0 1 3 4 5 6 7 7 1 0 8 6 12
9 0 2 0 0 2 2 1 7 0
5 1 0 0 3 2
6 0 0 1 1 5 5
8 1 0 3 3 0 1 3 4
11 0 0 2 1 1 3 7 5 9 2 6
19 0 0 0 3 5 6 2 2 0 0 5 8 0 0 1 3 4 4 12
1 1
14 1 2 1 1 2 4 3 8 6 1 2 6 9 10
12 0 1 3 1 3 3 3 6 3 10 7 6
19 0 1 2 4 2 4 0 7 3 4 0 9 5 13 12 3 17 1 4
13 0 1 3 3 0 5 7 7 2 2 5 7 6
6 1 2 0 2 2 1
12 1 2 0 1 2 1 5 4 6 10 0 4
18 1 0 3 2 3 5 3 5 2 2 1 9 5 2 0 13 12 14
3 0 1 1
6 0 0 1 1 0 4
5 1 1 2 2 5
4 1 1 1 2
17 1 2 1 3 5 5 4 3 6 5 6 7 8 4 13 13 3
5 0 0 0 1 4
20 0 1 1 3 0 6 0 0 1 5 11 7 5 1 14 11 8 15 7 5
18 0 2 1 0 5 6 2 8 6 9 10 3 7 11 12 8 0 18
5 1 0 3 4 0
13 0 2 3 2 1 4 7 0 7 9 1 4 13
17 1 2 3 3 2 1 3 8 5 2 4 9 13 14 4 14 2
3 1 1 3
18 0 1 3 1 0 2 2 5 1 2 10 0 13 10 4 6 0 1
13 1 2 0 3 4 2 5 1 0 3 3 7 13
10 1 0 0 3 2 4 5 1 1 4
19 1 0 2 3 1 1 4 3 7 10 5 4 6 9 4 3 12 11 16
16 0 2 2 2 3 2 5 6 4 1 7 4 1 7 4 11
8 0 1 3 1 0 5 6 6
15 1 2 1 3 4 2 7 3 5 8 10 0 1 7 10
13 0 1 0 4 0 3 1 4 3 5 2 1 2
12 0 0 0 0 3 4 0 2 1 9 9 3
3 1 0 0
17 1 0 1 1 1 3 6 7 0 6 3 6 12 0 8 0 11
12 1 2 3 1 4 4 1 4 1 1 4 0
5 0 1 1 4 5
11 0 1 0 4 0 5 7 1 8 7 7
8 1 0 2 0 3 1 3 3
5 1 0 1 1 1
17 1 1 2 0 3 2 4 6 3 7 8 5 0 13 12 12 15
3 0 1 2
17 0 0 0 0 2 4 0 2 7 8 3 5 5 2 10 11 16 12
15 0 1 2 4 3 5 4 5 6 10 6 1 5 8 0
15 0 2 3 0 1 1 4 1 8 9 3 11 11 10 9
3 0 1 1
7 1 0 1 2 5 1 1 2
13 0 0 0 0 1 1 1 3 6 5 4 7 10
9 1 0 2 3 4 6 6 5 4
1 0
15 1 2 3 5 1 2 7 7 3 11 1 2 5 4 6
3 1 1 2
1 1
10 0 2 3 0 0 3 5 7 4 2
2 1 0
2 1 0
7 0 1 0 1 3 1 3 3
12 1 2 4 5 1 5 4 10 7 2 12 12
8 1 1 0 3 5 4 4 7
11 1 1 0 2 3 4 3 5 3 8 7
10 1 2 4 0 5 4 5 6 2 7
6 1 2 2 0 5 3
2 1 1
11 1 2 3 4 3 2 8 0 2 7 6
2 0 0
14 0 1 3 1 3 5 6 8 7 3 2 6 3 2
14 1 1 3 1 3 4 5 6 6 10 4 7 13 8
16 1 0 1 2 2 4 6 8 9 6 9 11 7 12 9 0
12 0 1 0 0 2 4 4 2 1 6 10 7
5 1 1 0 0 0
6 0 0 1 2 0 3
15 1 2 2 1 2 5 3 8 9 5 1 4 3 2 12
12 0 1 1 2 5 4 5 1 8 5 0 4
16 0 0 1 1 1 2 5 6 8 2 7 6 0 10 7 14
`

type testCase struct {
	n   int
	arr []int
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// solve replicates the logic of 1364C.go.
func solve(tc testCase) string {
	n := tc.n
	ans := make([]int, n)
	sl := make([]int, n)
	ct := 0
	y := 0
	ok := true
	x := 0
	for i := 0; i < n; i++ {
		x = tc.arr[i]
		if x > i+1 {
			ok = false
		}
		if ok {
			sl[ct] = i
			ct++
			for y < x {
				ct--
				pos := sl[ct]
				ans[pos] = y
				y++
			}
		}
	}
	if !ok {
		return "-1"
	}
	fillVal := x + 1
	for ct > 0 {
		ct--
		pos := sl[ct]
		ans[pos] = fillVal
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func parseTestcases() ([]testCase, error) {
	rawLines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	start := 0
	count := 0
	if fields := strings.Fields(lines[0]); len(fields) == 1 {
		if v, err := strconv.Atoi(fields[0]); err == nil {
			count = v
			start = 1
		}
	}
	if count == 0 {
		count = len(lines)
		start = 0
	}
	if start+count != len(lines) {
		return nil, fmt.Errorf("testcase count mismatch: declared %d actual %d", count, len(lines)-start)
	}
	cases := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		fields := strings.Fields(lines[start+i])
		if len(fields) < 1 {
			return nil, fmt.Errorf("empty case %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields) < 1+n {
			return nil, fmt.Errorf("case %d expected at least %d numbers got %d", i+1, 1+n, len(fields))
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[1+j])
			if err != nil {
				return nil, err
			}
			arr[j] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		var input strings.Builder
		input.WriteString(strconv.Itoa(tc.n))
		for _, v := range tc.arr {
			input.WriteByte(' ')
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		want := solve(tc)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
