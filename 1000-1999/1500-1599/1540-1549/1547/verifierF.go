package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `
100
5
5 6 9 1 8
2
1 3
1
6
4
4 7 9 2
5
4 1 4 7 5
2
7 3
1
3
5
10 8 3 3 1
1
4
2
3 3
3
6 4 9
2
3 4
4
5 1 6 7
2
3 5
1
6
3
10 10 1
5
6 2 5 6 5
4
6 3 8 8
2
1 5
1
6
4
1 9 7 6
4
10 1 8 1
2
10 4
1
4
4
6 9 6 9
3
8 2 10
3
5 1 7
1
4
3
9 10 6
2
6 5
5
2 5 6 5 3
1
3
3
8 3 1
1
10
5
7 1 4 10 6
3
8 7 3
1
1
4
6 4 3 10
2
7 2
2
7 6
2
1 7
3
3 8 10
2
9 8
4
6 8 5 5
4
7 3 2 7
5
3 8 6 3 2
4
5 9 9 9
3
2 6 10
1
5
3
9 5 8
3
5 6 3
5
1 8 9 5 6
3
8 5 9
3
6 5 6
4
6 3 8 6
3
9 3 9
2
4 6
4
5 2 7 3
5
10 9 7 5 10
5
5 1 4 3 10
4
10 3 4 3
1
8
2
3 1
2
2 6
2
8 4
5
1 7 8 6 7
5
2 10 4 4 6
1
6
4
5 7 2 9
3
1 9 10
3
2 5 9
5
6 10 5 6 3
4
7 10 9 6
4
3 3 10 7
5
8 4 3 10 2
3
1 7 2
3
10 10 9
2
6 10
4
7 7 4 8
3
8 7 7
2
10 10
3
5 8 5
4
1 6 5 8
3
3 8 1
1
10
4
4 5 1 3
4
1 8 9 9
3
4 8 1
2
8 5
2
5 5
4
10 8 9 10
1
1
2
5 5
5
6 10 5 9 1
4
6 6 10 3
1
1
3
9 8 2
5
4 1 7 7 10
5
8 7 8 7 4
3
8 2 5

`

const mask = (1 << 30) - 1

type testCase struct {
	arr []int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildSparse(arr []int) ([][]int, []int) {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	K := log[n] + 1
	st := make([][]int, K)
	st[0] = make([]int, n)
	copy(st[0], arr)
	for k := 1; k < K; k++ {
		size := n - (1 << k) + 1
		st[k] = make([]int, size)
		step := 1 << (k - 1)
		for i := 0; i < size; i++ {
			st[k][i] = gcd(st[k-1][i], st[k-1][i+step])
		}
	}
	return st, log
}

func query(st [][]int, log []int, l, r int) int {
	length := r - l
	k := log[length]
	return gcd(st[k][l], st[k][r-(1<<k)])
}

func allEqual(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] != arr[0] {
			return false
		}
	}
	return true
}

func solveCase(tc testCase) int {
	n := len(tc.arr)
	if allEqual(tc.arr) {
		return 0
	}
	g := 0
	for _, v := range tc.arr {
		g = gcd(g, v)
	}
	b := make([]int, 2*n)
	copy(b, tc.arr)
	copy(b[n:], tc.arr)
	st, log := buildSparse(b)
	low, high := 1, n
	for low < high {
		mid := (low + high) / 2
		ok := true
		for i := 0; i < n; i++ {
			if query(st, log, i, i+mid) != g {
				ok = false
				break
			}
		}
		if ok {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low - 1
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	pos := 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", i+1)
		}
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		pos++
		if pos+n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough elements", i+1)
		}
		tc := testCase{arr: make([]int, n)}
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse arr[%d]: %v", i+1, j, err)
			}
			tc.arr[j] = v
		}
		pos += n
		cases = append(cases, tc)
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("trailing data after parsing")
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	var inputBuilder strings.Builder
	inputBuilder.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		inputBuilder.WriteString(strconv.Itoa(len(tc.arr)))
		inputBuilder.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			inputBuilder.WriteString(strconv.Itoa(v))
		}
		inputBuilder.WriteByte('\n')
	}

	got, err := runCandidate(bin, inputBuilder.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	outs := strings.Fields(got)
	if len(outs) != len(cases) {
		fmt.Printf("expected %d outputs, got %d\n", len(cases), len(outs))
		os.Exit(1)
	}

	for i, tc := range cases {
		expected := solveCase(tc)
		gotVal, err := strconv.Atoi(outs[i])
		if err != nil {
			fmt.Printf("case %d: non-integer output %q\n", i+1, outs[i])
			os.Exit(1)
		}
		if gotVal != expected {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %d\n", i+1, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
