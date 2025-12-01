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
	n   int
	k   int
	arr []int
}

// Testcases embedded from testcasesB.txt (one case per line, n k then n numbers).
const rawTestcases = `2 3
7 6
5 4
2 2 9 8 8
6 3
2 8 2 8 7 1
5 3
3 3 10 7 2
2 1
4 4
1 4
1
2 4
9 9
5 4
8 10 4 7 2
6 2
5 10 3 7 4 6
2 1
1 9
8 2
2 8 7 5 4 1 4 10
3 1
4 8 7
6 5
3 2 10 8 3 10
7 4
9 8 6 8 8 4 9
4 1
6 6 6 1
3 3
10 3 7
5 4
2 2 9 1 2
4 2
1 5 1 8
6 2
3 8 6 9 7 9
1 5
2
2 4
4 5
7 4
7 10 10 4 1 1 3
5 5
10 5 6 2 8
5 3
7 7 7 1 3
3 2
5 6 1
1 4
7
3 4
10 2 3
6 4
1 10 8 7 8 1
2 4
3 1
1 5
10
3 3
2 9 6
4 4
8 2 1 10
8 5
6 2 10 5 3 7 5 2
4 1
7 8 6 4
8 3
2 1 1 8 5 1 9 10
4 2
2 9 9 7
5 1
3 7 10 7 2
2 4
2 2
7 2
1 8 7 7 1 8 6
5 1
6 2 2 6 1
6 3
3 1 4 6 2 10
3 2
1 4 2
1 3
6
1 5
4
3 2
8 2 8
6 3
3 1 4 6 6 8
5 3
9 6 3 10 2
2 5
10 5
3 4
3 3 4
6 5
4 4 3 5 6 7
1 2
10
1 4
2
2 2
7 5
7 2
10 7 5 6 2 4 8
6 5
1 7 7 1 7 6
8 2
6 5 8 2 3 2 5 2
3 4
7 3 7
7 2
4 8 6 9 3 6 8
2 4
4 5
1 4
10
8 1
4 5 2 5 9 10 3 7
8 1
8 4 9 7 5 1 2 5
1 1
5
7 5
10 7 8 2 5 6 5
4 5
2 1 2 5
5 5
6 2 9 4 3
2 4
5 5
3 5
9 4 9
2 4
9 7
5 3
8 6 10 3 3
2 1
7 7
8 2
9 5 6 8 7 4 8 8
6 4
1 8 5 3 8 1
4 1
6 8 7 1
2 1
7 1
6 1
2 10 1 5 5 4
3 5
5 4 2
7 4
6 7 3 6 7 7 3
8 2
9 6 3 4 3 8 6 7
7 4
7 4 4 8 4 10 1
7 1
4 2 3 6 1 3 4
5 5
2 9 5 6 7
8 1
9 9 7 10 8 8 5 8
4 3
5 1 1 1
3 3
1 5 1
3 1
7 4 10
7 5
4 8 4 6 10 2 10
2 3
6 9
8 3
5 1 9 1 4 6 2 4
6 2
4 5 5 5 9 7
5 4
6 4 1 5 9
2 1
8 8
8 1
7 8 8 8 2 2 2 4
2 2
7 4
8 5
2 7 9 7 1 3 4 8
4 2
5 6 6 7
2 5
5 10
4 3
8 9 10 8
5 3
4 1 2 10 2
3 4
4 4 5`

func parseTestcases() ([]testcase, error) {
	tokens := strings.Fields(rawTestcases)
	pos := 0
	var res []testcase
	for pos < len(tokens) {
		if pos+1 >= len(tokens) {
			return nil, fmt.Errorf("incomplete testcase header")
		}
		n, err1 := strconv.Atoi(tokens[pos])
		k, err2 := strconv.Atoi(tokens[pos+1])
		pos += 2
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid n or k")
		}
		if pos+n > len(tokens) {
			return nil, fmt.Errorf("not enough values for n=%d", n)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(tokens[pos+i])
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		pos += n
		res = append(res, testcase{n: n, k: k, arr: arr})
	}
	return res, nil
}

func solve(n, k int, arr []int) (int, []string) {
	freq := make(map[int]int)
	for i := 0; i < n; i++ {
		bi := arr[i] - i*k
		if bi >= 1 {
			freq[bi]++
		}
	}
	bestH, bestCnt := 1, 0
	for h, cnt := range freq {
		if cnt > bestCnt {
			bestCnt = cnt
			bestH = h
		}
	}
	var ops []string
	for i := 0; i < n; i++ {
		desired := bestH + i*k
		if arr[i] != desired {
			if arr[i] < desired {
				ops = append(ops, fmt.Sprintf("+ %d %d", i+1, desired-arr[i]))
			} else {
				ops = append(ops, fmt.Sprintf("- %d %d", i+1, arr[i]-desired))
			}
		}
	}
	return len(ops), ops
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		count, ops := solve(tc.n, tc.k, tc.arr)
		expected := fmt.Sprintf("%d", count)
		if count > 0 {
			expected += "\n" + strings.Join(ops, "\n")
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, out.String())
			os.Exit(1)
		}

		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\n got:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
