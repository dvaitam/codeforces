package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
2 2 15
3 5
2 6
4 1 16
4 5 1 7
10
4 5 15
1 4 5 2
9 12 11 6 8
4 1 16
4 9 7 3
10
1 1 9
1
2
3 4 17
10 7 8
1 6 3 5
3 2 20
2 1 7
5 8
5 3 4
9 4 3 2 7
12 6 11
5 2 19
9 11 4 12 2
6 8
5 4 9
5 8 6 3 9
1 11 12 7
4 5 15
3 9 5 6
12 10 7 2 11
4 5 11
8 11 1 12
6 3 7 5 2
5 2 1
6 8 7 1 2
11 3
4 1 12
1 2 10 3
9
3 2 5
5 2 1
7 8
3 4 6
6 7 9
2 4 10 3
2 2 6
4 3
6 2
4 4 8
4 8 10 5
9 1 6 7
2 1 6
3 1
4
2 5 10
10 2
11 6 7 5 8
1 5 18
12
7 10 8 9 3
4 2 11
5 1 9 8
3 4
1 3 1
3
1 5 2
5 4 18
4 8 12 6 2
9 1 5 11
5 4 11
5 1 9 11 4
6 2 3 12
2 2 9
6 3
4 5
4 3 16
6 4 1 3
9 2 10
4 4 15
1 7 8 4
6 2 9 10
2 1 5
4 2
6
5 1 14
9 7 1 3 4
8
2 2 9
3 6
5 1
5 3 20
9 4 5 8 12
7 10 3
3 2 1
2 5 6
1 3
4 2 7
5 1 7 8
2 10
4 3 4
10 6 4 9
8 7 3
2 2 3
5 3
4 1
3 4 10
10 3 9
8 1 6 4
3 5 1
3 7 12
4 2 10 9 5
5 2 8
12 3 4 7 6
9 11
2 4 4
10 1
7 5 8 4
1 2 18
6
2 4
4 5 11
12 2 5 3
10 9 4 6 8
4 1 7
1 6 4 3
8
2 3 12
3 8
1 4 7
2 3 1
8 6
1 7 2
1 1 8
3
1
1 3 14
3
2 5 7
2 3 10
3 6
4 5 7
1 4 9
9
10 8 7 1
5 1 14
7 3 1 9 11
8
5 2 3
6 4 3 11 1
12 5
5 2 3
7 10 2 8 3
9 12
1 3 11
7
1 6 8
4 1 12
5 3 8 2
9
3 2 13
6 3 4
7 1
3 5 10
9 8 1
12 6 7 2 5
4 1 13
6 8 1 9
4
1 3 19
5
7 2 6
5 3 13
5 4 2 6 11
9 10 12
3 2 5
6 7 8
1 5
5 2 12
6 5 11 12 8
10 4
5 4 6
1 3 10 12 8
4 5 2 6
3 5 7
2 9 7
6 4 12 3 1
2 5 4
9 11
10 7 12 3 8
3 1 14
5 3 8
7
5 3 11
4 7 3 1 10
12 11 9
4 3 15
1 9 7 6
8 3 2
3 5 6
11 4 6
10 1 8 9 2
5 2 9
3 7 2 8 4
11 12
4 5 4
4 3 8 2
11 10 7 6 9
4 3 20
6 2 5 1
8 3 7
2 4 13
7 10
1 5 4 3
5 2 19
12 5 6 1 7
8 9
2 1 3
5 3
4
1 1 7
1
4
5 4 2
6 1 12 7 3
11 10 5 8
2 5 5
7 5
11 1 3 2 9
4 1 17
1 9 7 2
5
5 1 3
9 3 5 4 10
6
3 3 17
8 2 4
7 3 1
5 2 17
1 6 2 4 8
10 5
4 5 6
12 6 3 8
1 7 11 10 4
1 1 4
4
3
3 5 15
6 7 9
8 2 11 3 1
2 2 1
3 5
2 1
4 1 10
7 1 3 8
2
4 1 15
6 5 9 8
4
4 2 1
8 5 4 10
1 3
1 1 1
3
1
2 4 20
1 6
10 9 7 5
4 1 2
9 7 5 1
10
5 1 3
6 12 2 8 1
4
5 3 2
1 7 6 3 9
4 12 11
2 5 11
1 8
9 2 5 3 7
3 3 6
4 3 2
5 7 8
3 5 5
7 10 1
4 2 8 5 3
2 1 1
3 5
1
1 1 4
2
1
4 1 14
7 6 4 2
5
5 2 18
7 2 3 12 6
1 8
`

type testCase struct {
	n int
	m int
	k int64
	a []int
	b []int
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func exgcd(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := exgcd(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func solveCase(tc testCase) int64 {
	n := tc.n
	m := tc.m
	k := tc.k
	a := tc.a
	b := tc.b

	maxVal := n
	if m > maxVal {
		maxVal = m
	}
	maxVal = 2*maxVal + 5
	posA := make([]int, maxVal)
	posB := make([]int, maxVal)
	for i := range posA {
		posA[i] = -1
		posB[i] = -1
	}
	for i, v := range a {
		if v < len(posA) {
			posA[v] = i
		}
	}
	for i, v := range b {
		if v < len(posB) {
			posB[v] = i
		}
	}

	n1 := int64(n)
	m1 := int64(m)
	g := gcd(n1, m1)
	lcm := n1 / g * m1
	inv := int64(0)
	if g == 1 {
		_, invTmp, _ := exgcd(n1, m1)
		inv = (invTmp%m1 + m1) % m1
	} else {
		_, invTmp, _ := exgcd(n1/g, m1/g)
		inv = (invTmp%(m1/g) + m1/g) % (m1 / g)
	}

	matches := make([]int64, 0)
	for c := 0; c < maxVal; c++ {
		ia := posA[c]
		ib := posB[c]
		if ia == -1 || ib == -1 {
			continue
		}
		i64 := int64(ia)
		j64 := int64(ib)
		if (i64-j64)%g != 0 {
			continue
		}
		var t int64
		if g == 1 {
			diff := (j64 - i64) % m1
			if diff < 0 {
				diff += m1
			}
			t = (i64 + n1*((diff*inv)%m1)) % lcm
		} else {
			mg := m1 / g
			diff := ((j64 - i64) / g) % mg
			if diff < 0 {
				diff += mg
			}
			t = (i64 + n1*((diff*inv)%mg)) % lcm
		}
		matches = append(matches, t+1)
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i] < matches[j] })

	diffsPerCycle := lcm - int64(len(matches))
	cycles := int64(0)
	if diffsPerCycle > 0 {
		cycles = (k - 1) / diffsPerCycle
	}
	k -= cycles * diffsPerCycle
	base := cycles * lcm

	left, right := int64(1), lcm
	for left < right {
		mid := (left + right) / 2
		eq := int64(sort.Search(len(matches), func(i int) bool { return matches[i] > mid }))
		diff := mid - eq
		if diff >= k {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return base + left
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	res := make([]testCase, 0, t)
	idx := 1
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx+2 >= len(fields) {
			return nil, fmt.Errorf("case %d: not enough header fields", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		m, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse m: %v", caseIdx+1, err)
		}
		k, err := strconv.ParseInt(fields[idx+2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse k: %v", caseIdx+1, err)
		}
		idx += 3
		if idx+n+m > len(fields) {
			return nil, fmt.Errorf("case %d: not enough numbers", caseIdx+1)
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a%d: %v", caseIdx+1, i+1, err)
			}
			a[i] = v
		}
		idx += n
		b := make([]int, m)
		for i := 0; i < m; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b%d: %v", caseIdx+1, i+1, err)
			}
			b[i] = v
		}
		idx += m
		res = append(res, testCase{n: n, m: m, k: k, a: a, b: b})
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after parsing")
	}
	return res, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, tc.k)
		for idx, v := range tc.a {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		expected := strconv.FormatInt(solveCase(tc), 10)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
