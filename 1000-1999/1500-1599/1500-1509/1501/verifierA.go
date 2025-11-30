package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesA.txt so the verifier is self-contained.
const testcasesRaw = `100
4 7 8 13 22 30 37 42 50 2 4 1 4
2 5 8 10 20 2 4
5 3 8 10 12 18 26 35 37 43 50 2 4 5 1 4
4 8 17 22 23 32 33 35 42 5 5 5 0
5 8 14 18 24 26 30 40 44 48 51 4 3 0 0 2
5 8 10 15 24 29 31 40 46 55 59 4 4 4 2 3
1 10 17 2
5 4 9 12 16 19 20 30 35 43 45 0 5 1 1 0
1 9 16 5
5 5 14 18 22 32 39 49 54 62 70 5 5 5 2 0
3 10 12 20 30 36 40 1 0 5
3 2 6 12 15 21 28 0 0 1
2 1 11 20 30 5 0
1 2 6 4
5 2 9 11 17 19 20 30 31 35 38 5 0 3 1 5
1 1 10 3
5 2 7 9 13 15 20 26 33 36 37 4 3 0 4 0
4 4 9 15 23 33 36 40 41 5 1 1 2
5 5 7 17 25 28 29 37 44 54 63 2 5 2 3 5
3 3 12 13 21 23 29 5 0 4
3 3 7 15 21 31 36 5 2 4
5 3 8 15 22 24 25 35 39 45 48 1 1 5 3 3
5 7 8 15 25 32 33 36 44 46 51 5 1 3 4 3
5 10 11 12 20 26 31 39 40 47 51 4 5 0 5 1
1 7 14 2
1 4 5 5
1 9 19 0
2 2 12 16 21 2 5
2 2 10 17 19 0 2
4 2 7 10 19 25 27 30 35 0 0 0 1
3 9 15 21 31 32 42 5 3 5
4 7 13 22 25 29 36 46 51 0 1 1 2
3 6 12 14 20 30 31 0 2 1
2 10 15 21 28 4 1
3 2 10 14 15 20 23 4 5 0
3 7 13 18 25 27 29 4 3 3
3 6 8 16 18 26 33 0 2 2
2 3 13 20 22 0 0
2 4 5 12 13 0 3
5 9 14 22 30 40 44 51 53 59 63 2 4 1 3 1
3 2 4 5 14 22 26 0 3 3
3 4 5 9 19 22 24 1 3 3
3 9 12 14 24 32 35 4 3 5
4 9 17 23 31 39 43 52 62 1 0 2 5
3 6 7 16 19 24 34 1 3 4
3 8 10 12 21 22 24 1 1 0
3 1 9 15 18 21 29 2 4 3
5 9 10 20 22 31 41 43 50 54 59 4 4 3 3 3
5 10 14 15 16 19 24 33 43 48 54 0 3 2 2 3
4 7 8 11 14 18 23 29 30 0 3 3 1
4 10 12 15 21 28 29 39 47 3 3 0 0
4 3 4 5 15 25 28 34 36 5 4 5 2
2 7 15 17 18 4 5
4 10 16 18 28 33 36 43 48 5 5 0 4
2 1 8 16 22 1 3
3 2 3 4 12 17 18 4 5 4
5 4 8 10 19 28 35 44 49 51 54 3 4 3 0 0
4 2 4 11 14 15 23 30 37 0 3 2 5
3 2 8 10 12 18 19 2 2 1
1 4 10 0
5 3 7 8 12 14 15 20 26 27 37 1 1 1 3 0
4 6 11 14 15 19 25 31 39 2 2 4 5
3 3 13 15 17 26 36 2 1 3
2 3 7 13 22 1 1
2 5 11 18 19 1 4
1 7 9 5
1 3 10 2
5 7 10 20 27 32 38 40 44 52 58 5 4 0 3 3
1 7 13 3
2 6 11 19 21 1 0
3 2 11 21 24 32 39 1 3 3
2 4 12 18 27 1 2
4 2 10 14 19 20 28 38 46 0 1 2 0
3 9 19 22 29 37 39 5 3 1
5 7 12 13 15 20 21 22 27 34 43 4 5 3 3 0
3 6 11 15 25 27 28 0 2 2
5 6 8 17 21 24 26 33 38 43 52 1 4 4 5 1
5 2 9 18 25 30 35 43 49 59 62 1 0 5 0 3
4 10 18 21 30 35 41 49 56 1 3 3 5
5 6 14 15 23 28 31 39 40 50 54 0 2 3 3 0
5 2 4 11 12 18 19 21 31 32 37 5 5 2 5 1
2 10 15 19 21 3 3
3 7 10 16 23 30 33 3 5 1
5 6 9 13 16 24 30 37 44 52 59 5 1 1 3 1
5 1 8 9 13 15 18 24 25 28 32 4 2 4 0 5
5 5 11 18 26 27 36 45 52 62 70 3 2 5 3 1
3 5 6 7 8 11 17 0 2 5
1 3 5 3
2 10 17 26 30 3 1
3 10 12 22 24 30 36 4 3 2
3 1 10 11 15 21 23 1 4 2
2 4 9 14 19 4 3
3 8 14 18 19 24 33 0 0 3
4 8 9 16 24 32 40 42 44 0 1 0 1
4 4 12 22 24 31 40 47 48 1 1 3 1
2 5 11 17 24 0 4
3 10 19 23 28 36 45 4 3 4
3 5 9 10 12 22 24 1 5 3
2 4 9 10 19 4 3
1 2 9 5`

type testCase struct {
	n  int
	a  []int
	b  []int
	tm []int
}

func solveCase(tc testCase) int {
	prevDepart := 0
	prevB := 0
	for i := 0; i < tc.n; i++ {
		travel := tc.a[i] - prevB + tc.tm[i]
		arrival := prevDepart + travel
		if i == tc.n-1 {
			return arrival
		}
		stay := (tc.b[i] - tc.a[i] + 1) / 2
		depart := arrival + stay
		if depart < tc.b[i] {
			depart = tc.b[i]
		}
		prevDepart = depart
		prevB = tc.b[i]
	}
	return 0
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
		if idx >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseIdx+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", caseIdx+1, err)
		}
		idx++
		if idx+3*n > len(fields) {
			return nil, fmt.Errorf("case %d: not enough numbers", caseIdx+1)
		}
		a := make([]int, n)
		b := make([]int, n)
		tm := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a%d: %v", caseIdx+1, i+1, err)
			}
			a[i] = v
		}
		idx += n
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b%d: %v", caseIdx+1, i+1, err)
			}
			b[i] = v
		}
		idx += n
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("case %d: parse tm%d: %v", caseIdx+1, i+1, err)
			}
			tm[i] = v
		}
		idx += n
		res = append(res, testCase{n: n, a: a, b: b, tm: tm})
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for j := 0; j < tc.n; j++ {
			sb.WriteString(strconv.Itoa(tc.a[j]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(tc.b[j]))
			sb.WriteByte('\n')
		}
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.tm[j]))
		}
		sb.WriteByte('\n')

		expected := strconv.Itoa(solveCase(tc))
		got, err := runCandidate(bin, sb.String())
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
