package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC2.txt so the verifier is self-contained.
const testcasesRaw = `7 4 3 4 5 6 7 2 1 2 7 1 5 5 7 1 4
9 1 3 5 2 1 6 7 8 9 4 1 7
10 2 5 10 8 2 4 3 7 9 6 1 2 5 2 5
1 5 1 1 1 1 1 1 1 1 1 1 1
2 4 1 2 1 2 2 2 1 2 2 2
2 4 1 2 1 2 1 2 2 2 1 2
1 4 1 1 1 1 1 1 1 1 1
3 1 1 3 2 2 3
8 4 8 6 4 2 1 5 7 3 1 3 1 6 3 6 7 8
4 3 2 3 1 4 1 4 2 3 2 3
9 4 2 1 8 4 5 9 6 3 7 1 8 4 8 6 9 1 4
1 0 1
3 1 3 1 2 3 3
2 4 1 2 2 2 1 2 1 2 2 2
7 3 2 7 5 1 4 3 6 1 7 2 7 4 4
9 1 9 3 5 7 2 4 8 1 6 4 9
1 5 1 1 1 1 1 1 1 1 1 1 1
10 2 4 6 3 8 10 9 7 1 2 5 3 8 2 5
10 3 8 2 1 3 10 6 9 7 4 5 1 2 3 6 10 10
8 0 6 2 5 1 7 3 4 8
10 5 9 5 7 2 4 6 8 1 10 3 5 7 2 6 5 6 4 10 9 9
1 2 1 1 1 1 1
5 2 4 5 1 3 2 1 3 1 5
9 0 6 4 7 9 2 5 1 8 3
3 5 2 1 3 1 2 3 3 2 3 1 2 2 3
6 4 3 1 5 4 6 2 3 5 1 6 1 6 3 3
6 0 3 4 2 6 5 1
3 5 3 1 2 3 3 2 2 1 2 3 3 1 2
2 4 1 2 1 2 1 2 2 2 1 1
1 5 1 1 1 1 1 1 1 1 1 1 1
7 3 1 4 6 7 2 5 3 4 6 1 4 2 3
9 0 3 7 5 6 2 1 9 8 4
3 1 3 1 2 2 3
7 0 1 7 6 3 4 2 5
7 3 5 4 7 3 2 1 6 2 2 3 7 2 7
3 4 3 1 2 1 1 2 3 1 3 1 3
3 4 2 1 3 1 3 2 3 3 3 1 3
8 1 3 6 8 5 7 2 4 1 2 8
3 0 2 1 3
2 0 2 1
7 3 1 2 5 7 4 3 6 3 6 1 7 2 5
4 5 4 3 1 2 1 4 2 4 2 3 1 2 3 3
6 2 6 4 1 2 3 5 2 6 4 5
6 3 1 6 5 2 4 3 3 6 3 6 4 5
7 4 2 1 7 3 5 6 4 2 4 2 3 2 5 2 5
10 4 5 1 2 6 4 7 3 8 10 9 4 5 2 8 1 4 2 10
5 0 4 1 3 5 2
6 4 2 6 4 5 1 3 2 3 4 4 3 6 1 3
10 3 4 7 10 8 6 5 9 3 1 2 4 7 3 6 3 8
3 4 1 2 3 3 3 1 1 1 2 1 3
3 4 3 2 1 2 2 3 3 1 3 1 3
6 1 1 5 2 4 6 3 3 6
2 3 2 1 1 2 2 2 1 2
10 3 6 1 5 10 4 8 7 9 2 3 5 7 2 5 6 9
9 0 4 3 7 6 2 8 5 9 1
8 0 7 5 3 4 1 2 8 6
1 1 1 1 1
1 4 1 1 1 1 1 1 1 1 1
9 1 6 8 2 9 5 1 7 3 4 4 5
7 5 3 2 4 5 7 6 1 2 7 2 7 5 7 1 2 1 3
6 1 6 2 4 1 5 3 2 3
8 2 3 2 7 1 6 5 4 8 7 8 7 8
9 3 6 5 1 3 4 7 8 2 9 3 3 1 7 2 2
9 2 1 2 3 9 4 7 8 6 5 8 9 4 6
9 0 4 8 9 2 1 5 3 6 7
8 0 8 7 3 1 2 5 4 6
1 4 1 1 1 1 1 1 1 1 1
7 0 1 5 3 4 2 6 7
3 0 1 2 3
9 5 2 8 1 4 5 9 6 7 3 3 8 6 8 6 7 3 9 7 7
4 2 3 1 2 4 2 3 2 3
8 1 7 4 1 3 2 8 5 6 3 5
9 5 5 1 7 3 6 9 4 8 2 2 9 4 8 5 7 5 9 4 6
7 2 4 3 5 2 7 1 6 5 7 1 4
5 4 1 5 2 4 3 1 4 1 1 2 3 3 5
9 3 6 1 9 2 3 4 8 7 5 2 2 1 9 6 9
10 0 9 2 5 6 3 7 8 1 4 10
3 4 2 1 3 2 2 1 1 1 3 1 1
4 2 3 1 2 4 2 3 1 2
7 1 6 7 4 3 1 5 2 2 6
5 0 2 1 4 3 5
1 1 1 1 1
7 1 6 7 3 1 5 4 2 2 6
3 3 1 3 2 1 3 3 3 2 3
4 0 1 4 2 3
6 4 1 6 4 2 3 5 1 2 1 4 5 6 1 5
3 5 3 1 2 2 2 3 3 1 1 1 3 2 2
6 3 6 4 2 3 1 5 2 5 6 6 5 5
3 5 3 2 1 2 3 3 3 1 2 1 2 1 2
2 4 1 2 1 2 1 1 1 1 1 1
10 3 10 3 1 8 6 9 5 4 2 7 2 9 3 10 4 9
10 4 2 1 5 4 3 6 8 9 7 10 6 6 1 4 4 8 4 9
4 0 3 4 1 2
8 0 1 8 2 4 3 6 7 5
2 4 2 1 2 2 1 1 1 2 1 1
8 2 7 4 1 2 8 6 5 3 2 6 4 6
4 2 1 2 3 4 1 4 2 4
2 2 2 1 1 2 1 2
10 1 3 4 8 9 1 2 6 5 7 10 7 9
3 0 1 3 2`

type testCase struct {
	n       int
	q       int
	arr     []int
	queries [][2]int
}

func solveCase(tc testCase) string {
	a := make([]int, tc.n+2)
	a[0] = 0
	copy(a[1:], tc.arr)
	var ans int64
	for i := 1; i <= tc.n; i++ {
		d := a[i] - a[i-1]
		if d > 0 {
			ans += int64(d)
		}
	}
	var out strings.Builder
	out.WriteString(strconv.FormatInt(ans, 10))
	out.WriteByte('\n')

	adjust := func(idx int, delta *int64) {
		if idx < 1 || idx > tc.n {
			return
		}
		d := a[idx] - a[idx-1]
		if d > 0 {
			*delta += int64(d)
		}
	}

	for _, qr := range tc.queries {
		l, r := qr[0], qr[1]
		idxs := [4]int{l, l + 1, r, r + 1}
		seen := make(map[int]struct{}, 4)
		var delta int64
		for _, idx := range idxs {
			if idx < 1 || idx > tc.n {
				continue
			}
			if _, ok := seen[idx]; ok {
				continue
			}
			seen[idx] = struct{}{}
			adjust(idx, &delta)
		}
		ans -= delta
		// swap
		a[l], a[r] = a[r], a[l]
		seen = make(map[int]struct{}, 4)
		delta = 0
		for _, idx := range idxs {
			if idx < 1 || idx > tc.n {
				continue
			}
			if _, ok := seen[idx]; ok {
				continue
			}
			seen[idx] = struct{}{}
			adjust(idx, &delta)
		}
		ans += delta
		out.WriteString(strconv.FormatInt(ans, 10))
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		q, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: parse n/q: %v %v", idx+1, err1, err2)
		}
		expected := 2 + n + 2*q
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, expected, len(fields))
		}
		arr := make([]int, n)
		pos := 2
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
			pos++
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, errL := strconv.Atoi(fields[pos])
			r, errR := strconv.Atoi(fields[pos+1])
			if errL != nil || errR != nil {
				return nil, fmt.Errorf("line %d: parse query %d: %v %v", idx+1, i+1, errL, errR)
			}
			queries[i] = [2]int{l, r}
			pos += 2
		}
		cases = append(cases, testCase{n: n, q: q, arr: arr, queries: queries})
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.q)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for _, qr := range tc.queries {
			fmt.Fprintf(&input, "%d %d\n", qr[0], qr[1])
		}
	}

	got, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var expected strings.Builder
	for i, tc := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		expected.WriteString(solveCase(tc))
	}

	if strings.TrimSpace(got) != strings.TrimSpace(expected.String()) {
		fmt.Printf("output mismatch\nexpected:\n%s\n\ngot:\n%s\n", expected.String(), got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
