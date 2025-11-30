package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
100
5
5 1 2 0 1
4
3 5
3 5
3 4
2 4
3
5 2 1
4
1 2
2 3
1 2
2 3
4
2 5 5 0
1
2 4
4
2 2 0 1
3
3 4
1 4
1 3
3
0 5 4
3
1 2
2 3
2 3
5
5 1 5 1 0
1
2 5
2
3 4
1
1 2
6
4 0 1 5 2 3
2
5 6
3 5
3
5 3 2
2
1 2
1 3
5
5 0 2 4 3
3
4 5
2 4
2 4
3
1 0 1
4
2 3
2 3
1 3
1 2
3
4 4 2
4
2 3
2 3
1 2
2 3
6
4 5 0 2 3 1
1
1 5
6
4 4 4 4 4 5
4
4 6
5 6
2 5
4 5
6
5 5 0 5 3 1
1
1 6
2
1 0
1
1 2
3
0 5 0
2
1 3
2 3
5
1 3 5 0 4
2
1 5
1 5
3
1 0 4
2
2 3
2 3
4
4 3 3 2
3
3 4
2 4
2 4
5
2 5 0 1 5
4
3 5
4 5
2 4
2 3
6
1 0 2 4 4 1
5
2 6
1 5
2 6
2 3
1 6
2
5 5
1
1 2
3
4 1 5
3
2 3
1 3
2 3
4
4 4 1 4
2
3 4
3 4
2
1 3
1
1 2
3
5 2 2
3
1 3
2 3
1 3
2
0 1
3
1 2
1 2
1 2
2
3 4
4
1 2
1 2
1 2
1 2
6
2 3 0 4 5 1
2
5 6
4 5
5
2 3 2 0 3
1
4 5
3
2 2 2
5
1 3
2 3
2 3
2 3
2 3
6
1 0 3 4 3 4
2
4 6
2 5
6
1 0 0 1 0 1
4
2 3
4 6
3 5
5 6
5
4 5 0 2 3
4
2 5
4 5
4 5
2 5
4
2 1 1 5
3
1 2
2 4
1 3
4
5 4 1 3
1
3 4
3
3 3 0
1
1 2
3
2 5 3
5
1 2
1 2
2 3
1 2
2 3
2
5 0
5
1 2
1 2
1 2
1 2
1 2
4
4 1 4 2
1
3 4
4
0 2 2 2
1
2 3
4
0 5 0 5
5
2 3
2 3
2 4
2 3
1 4
6
3 4 0 5 0 5
1
1 3
5
3 4 1 5 1
5
1 2
1 5
4 5
2 4
1 4
5
4 1 4 2 4
2
1 5
4 5
3
2 3 5
2
1 2
1 2
2
2 2
3
1 2
1 2
1 2
5
0 4 3 1 3
2
4 5
2 5
5
4 5 2 0 5
5
2 3
4 5
2 5
1 2
2 3
4
0 2 0 1
5
3 4
3 4
2 3
2 3
1 4
3
5 3 5
5
2 3
1 2
2 3
2 3
2 3
4
2 0 2 5
1
2 4
3
2 1 1
2
1 3
1 3
3
2 3 0
5
1 3
1 3
1 3
1 2
2 3
3
1 0 2
3
1 3
2 3
1 2
5
1 2 3 2 2
1
2 4
4
4 1 0 3
4
2 4
3 4
2 3
2 3
5
5 0 3 4 1
1
1 5
3
2 5 0
4
2 3
1 3
2 3
1 2
3
5 5 0
1
1 3
4
3 4 1 3
2
3 4
2 3
2
2 1
1
1 2
2
5 5
5
1 2
1 2
1 2
1 2
1 2
6
5 2 1 4 2 2
1
4 5
2
5 0
4
1 2
1 2
1 2
1 2
4
5 2 2 5
4
1 2
2 4
2 3
3 4
4
0 3 1 0
2
1 4
3 4
3
5 2 1
5
1 2
2 3
2 3
1 3
2 3
3
1 5 3
1
1 2
3
4 5 0
1
2 3
3
1 3 2
3
2 3
2 3
1 2
3
1 0 1
1
1 2
6
5 4 3 3 5 3
4
1 6
5 6
5 6
2 5
6
2 5 3 4 4 0
4
5 6
4 6
5 6
1 4
6
5 2 5 4 2 3
4
1 6
3 5
1 5
3 6
2
0 5
3
1 2
1 2
1 2
5
0 4 5 0 3
2
3 5
3 4
5
4 3 3 2 2
1
3 4
4
5 0 1 3
4
2 4
2 3
3 4
3 4
5
2 3 0 0 5
2
1 4
3 5
4
0 3 5 2
2
2 4
1 4
6
2 2 2 4 5 2
1
3 4
5
4 1 2 4 2
3
4 5
4 5
2 4
3
0 5 0
4
2 3
1 3
1 2
2 3
3
5 1 3
1
2 3
4
0 1 0 0
3
2 4
1 3
2 4
4
2 1 0 5
5
1 3
3 4
3 4
2 3
3 4
2
4 1
1
1 2
3
4 3 4
4
2 3
2 3
1 2
2 3
4
5 2 1 3
1
1 4
2
1 4
3
1 2
1 2
1 2
3
1 4 1
4
2 3
2 3
2 3
1 2
4
2 2 0 4
3
2 3
2 3
2 4
3
1 4 1
5
1 2
2 3
1 3
2 3
2 3
3
3 3 2
1
2 3
4
1 2 3 0
3
2 4
3 4
1 3
3
1 1 4
4
2 3
1 3
1 2
2 3
6
0 0 3 2 4 5
2
4 5
5 6
5
4 3 3 4 3
4
4 5
3 4
2 4
1 5
`

type testCase struct {
	n       int
	arr     []int
	queries [][2]int
}

func parseTestcases(raw string) ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(raw))
	sc.Split(bufio.ScanWords)

	scanInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, err
		}
		return v, nil
	}

	t, err := scanInt()
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		n, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: n: %w", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := scanInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: value %d: %w", i+1, j+1, err)
			}
			arr[j] = v
		}
		q, err := scanInt()
		if err != nil {
			return nil, fmt.Errorf("case %d: q: %w", i+1, err)
		}
		queries := make([][2]int, q)
		for j := 0; j < q; j++ {
			l, err := scanInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: query %d l: %w", i+1, j+1, err)
			}
			r, err := scanInt()
			if err != nil {
				return nil, fmt.Errorf("case %d: query %d r: %w", i+1, j+1, err)
			}
			queries[j] = [2]int{l, r}
		}
		cases = append(cases, testCase{n: n, arr: arr, queries: queries})
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}
	return cases, nil
}

func solve(tc testCase) string {
	var sb strings.Builder
	for idx, qr := range tc.queries {
		l := qr[0] - 1
		r := qr[1] - 1
		good := false
		for i := l; i < r && !good; i++ {
			maxLeft := tc.arr[l]
			for j := l; j <= i; j++ {
				if tc.arr[j] > maxLeft {
					maxLeft = tc.arr[j]
				}
			}
			minRight := tc.arr[i+1]
			for j := i + 1; j <= r; j++ {
				if tc.arr[j] < minRight {
					minRight = tc.arr[j]
				}
			}
			if maxLeft < minRight {
				good = true
			}
		}
		if good {
			sb.WriteString("Yes")
		} else {
			sb.WriteString("No")
		}
		if idx+1 < len(tc.queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.Itoa(q[0]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q[1]))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expected := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
