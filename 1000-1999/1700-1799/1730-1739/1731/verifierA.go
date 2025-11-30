package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `5 8 30 4 47 8
3 39 47 26
3 37 39 47
5 42 2 34 15 19
3 37 32 38
6 46 26 31 20 43 45
6 17 18 24 47 40 3
3 8 23 27
2 5 2
2 16 48
3 26 9 4
4 21 48 38 6
5 12 32 8 14 24
4 24 40 44 33
5 9 18 33 2 48
3 10 13 42
4 23 22 36 37
5 12 29 27 47 45
3 41 33 20
4 27 44 48 13
2 44 25
2 1 47
4 45 22 20 11
3 46 35 18
6 3 7 22 15 46 15
5 31 16 30 7 40
2 6 7
2 39 41
2 6 16
5 18 40 16 2 24
4 41 30 1 26
4 30 9 15 49
6 35 45 36 45 1 6
4 39 20 6 7
5 28 40 7 49 50
2 12 50
5 22 38 7 37 23
3 28 10 34
2 42 42
2 8 49
4 8 22 26 5
3 13 8 12
4 6 16 41 35
3 5 22 3
4 12 5 18 19
5 36 10 11 6 12
2 22 27
2 16 23
4 46 29 1 44
4 8 12 47 13
6 20 25 31 20 30 7
4 44 27 15 21
2 9 26
5 49 3 43 25 49
5 23 33 50 2 27
3 35 5 12
5 25 50 45 30 44
6 34 21 27 50 16 46
4 32 21 22 39
6 7 34 1 35 17 16
6 32 49 47 44 15 46
3 39 3 24
3 2 32 6
2 9 49
2 32 18
4 38 22 4 25
4 42 8 5 45
5 31 48 18 17 27
2 46 6
2 25 7
5 48 20 14 24 10
4 13 33 2 6
5 50 20 15 14 20
4 30 41 39 37
4 29 35 31 31
3 6 35 32
5 3 12 26 23 14
2 26 13
3 7 19 2
6 45 45 27 11 17 25
6 45 33 34 49 10 19
4 38 45 39 12
4 42 2 11 31
5 29 22 48 11 17
6 35 5 29 18 12 41
5 42 13 4 42 37
3 25 6 43
3 25 16 2
5 5 8 4 50 27
5 40 24 49 31 18
2 4 46
3 12 1 23
6 40 45 26 4 14 47
6 16 9 4 13 36 32
6 33 9 4 30 48 26
2 15 44
4 49 3 23 45
6 14 25 48 28 1 14
6 14 32 11 46 32 42
4 45 7 29 8`

type testCase struct {
	n    int
	vals []int64
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(vals []int64) int64 {
	product := int64(1)
	for _, v := range vals {
		product *= v
	}
	ans := product + int64(len(vals)-1)
	return ans * 2022
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d expected %d numbers got %d", idx+1, n, len(fields)-1)
		}
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d parse value %d: %v", idx+1, i, err)
			}
			vals[i] = v
		}
		cases = append(cases, testCase{n: n, vals: vals})
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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

	for idx, tc := range cases {
		want := solve(tc.vals)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.vals {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')

		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
