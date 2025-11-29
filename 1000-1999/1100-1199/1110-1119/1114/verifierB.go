package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesRaw = `4 3 1 -2 -7 5 4
6 3 2 -4 -7 5 -10 2 3
7 1 3 4 -2 -3 8 -7 0 -10
3 1 3 7 -10 2
8 1 2 -10 6 -3 4 5 7 -3 1
4 3 1 4 -1 -10 3
7 3 1 -5 10 -1 -7 0 6 3
7 3 1 -1 -1 8 5 6 2 8
3 2 1 2 3 -5
9 3 3 1 -8 4 6 -7 -5 6 2 1
6 3 1 5 -9 -1 9 8 8
6 3 1 -5 6 -3 -10 -4 7
7 1 2 6 1 8 1 4 -2 7
7 3 1 2 6 -6 6 7 -4 3
4 2 2 8 7 -4 6
6 2 2 3 1 -10 7 7 9
7 2 2 9 -10 -3 10 -5 7 8
4 1 3 -2 -9 -8 -8
3 2 1 -2 -3 -2
3 3 1 1 -1 -8
4 1 2 6 -5 -2 10
8 2 2 0 5 5 -7 -10 -1 2 0
6 1 2 -7 -2 6 -4 9 3
3 1 1 2 -6 -9
8 1 2 6 3 7 -3 10 6 4 -3
7 3 1 2 8 0 10 3 -9 -1
4 1 1 -1 -8 -8 -1
5 3 1 3 8 -2 -6 -10
7 1 3 -4 8 4 -5 9 6 -9
6 1 2 -7 -4 8 3 8 -4
6 1 3 2 -1 6 5 -10 0
7 2 2 -10 -5 -4 0 8 -6 0
6 1 2 -7 2 7 1 7 5
7 1 1 -9 -8 -6 -5 -5 7 -4
6 2 3 6 -2 1 0 0 -7
5 1 3 5 -6 8 7 -7
5 1 2 -8 2 -6 -6 0
9 3 3 2 -8 8 7 -3 8 -8 -2 1
9 3 3 -7 4 -2 -7 -9 -1 -10 9 -10
3 2 1 -9 -4 -3
7 2 1 -7 4 -5 -3 -5 -7 3
6 3 2 7 -2 5 0 -7 -4
8 2 1 -10 -10 -1 9 0 4 2 0
6 1 1 0 9 4 -7 -2 -4
9 3 3 5 1 -2 -5 7 -4 -1 -4 -3
5 1 2 -8 4 -8 10 8
8 2 1 2 -1 -9 0 -5 0 8 -1
4 2 1 7 9 8 9
3 1 1 -10 -3 2
6 2 3 -8 -8 -10 10 -10 -1
5 2 2 -6 -7 6 0 -8
7 3 1 -5 -6 -6 0 -1 -7 6
7 2 1 -4 -6 7 -9 0 9 7
8 3 1 -5 -1 3 7 -5 -9 -3 -2
6 3 2 3 7 -2 7 4 7
6 1 2 0 -5 -2 5 -10 10
6 3 1 -9 1 8 -6 8 -6
4 2 2 2 8 2 -5
7 1 1 5 -10 -5 6 0 6 10
9 3 3 -3 -3 0 5 5 -3 3 0 7
9 3 3 -2 10 -3 -9 -8 6 10 1 -5
7 1 2 -1 -1 7 1 -5 4 9
3 1 3 6 8 2
4 1 2 3 -4 8 -9
6 3 2 10 1 2 6 -5 7
8 1 3 -8 -2 10 -7 -2 -8 -6 9
9 3 3 -8 4 -3 2 3 2 -5 0 4
6 3 2 -4 -7 3 9 7 3
6 3 2 -2 -3 2 7 -10 -4
7 2 3 -10 -10 10 9 -3 -2 -4
4 2 1 7 -4 -2 -1
7 2 3 4 -5 7 1 5 3 -7
6 3 2 -4 -1 -7 -10 -7 8
8 1 3 -1 10 -6 -8 6 1 8 -1
9 3 3 1 6 0 -10 -7 4 4 1 -1
7 2 2 8 5 -7 10 2 2 -4
7 1 2 10 9 6 -4 4 9 6
9 3 3 -1 -5 4 9 6 -4 1 6 -10
8 2 3 3 2 0 9 8 -8 5 -3
8 3 2 10 -10 3 10 -6 10 2 -2
4 1 3 -10 1 -2 3
8 3 2 -6 4 -2 5 -5 4 6 -9
5 3 1 8 3 -8 1 -8
8 2 1 -5 6 -5 -8 2 10 -2 9
5 1 3 -4 -3 0 -2 -8
9 3 3 1 4 6 7 -9 -5 -1 10 7
6 2 3 -3 2 7 2 -5 5
6 3 2 -3 -2 9 -3 -10 9
6 2 2 -3 -2 -4 -8 10 -5
7 2 3 -6 9 -2 4 6 -5 -6
6 3 2 1 -1 2 -3 -7 -4
8 3 2 -8 -7 -3 2 0 5 -7 -5
3 1 3 -10 -4 -9
9 3 3 9 4 0 -2 -7 9 -5 -7 -3
6 1 2 4 2 -5 -3 -3 -1
9 3 3 2 -4 4 -2 0 5 8 -7 -4
3 1 1 -10 5 0
6 3 2 -4 2 -5 10 -6 -10
3 2 1 7 -9 8
6 2 1 -8 4 10 -1 -10 -9
`

type node struct {
	val int64
	pos int
}

// solve mirrors 1114B.go logic.
func solve(n, m, k int, arr []int64) (int64, []int) {
	b := make([]node, n)
	for i := 0; i < n; i++ {
		b[i] = node{arr[i], i}
	}
	c := make([]node, n)
	copy(c, b)
	sort.Slice(c, func(i, j int) bool { return c[i].val > c[j].val })
	tar := m * k
	vis := make([]bool, n)
	var sum int64
	for i := 0; i < tar; i++ {
		sum += c[i].val
		vis[c[i].pos] = true
	}
	cuts := make([]int, 0, k-1)
	cnt, t := 0, 0
	for i := 0; i < n && t < k-1; i++ {
		if vis[i] {
			cnt++
		}
		if cnt == m {
			cuts = append(cuts, i+1)
			cnt = 0
			t++
		}
	}
	return sum, cuts
}

type testCase struct {
	n   int
	m   int
	k   int
	arr []int64
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	var tests []testCase
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("invalid test line: %q", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
		if len(fields) != 3+n {
			return nil, fmt.Errorf("line expects %d numbers got %d", 3+n, len(fields))
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[3+i], 10, 64)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		tests = append(tests, testCase{n: n, m: m, k: k, arr: arr})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(got)
	if len(fields) < 1 {
		return fmt.Errorf("no output")
	}
	sum, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid sum %q", fields[0])
	}
	expectedSum, cuts := solve(tc.n, tc.m, tc.k, tc.arr)
	if sum != expectedSum {
		return fmt.Errorf("expected sum %d got %d", expectedSum, sum)
	}
	if tc.k > 1 {
		if len(fields)-1 != len(cuts) {
			return fmt.Errorf("expected %d cuts got %d", len(cuts), len(fields)-1)
		}
		for i, f := range fields[1:] {
			v, err := strconv.Atoi(f)
			if err != nil || v != cuts[i] {
				return fmt.Errorf("cut %d expected %d got %q", i+1, cuts[i], f)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
