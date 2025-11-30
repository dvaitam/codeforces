package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `20 20 12 12 17 4 2 9 20 12 8 15 11 4 9 13 4 12 2 12 19
12 10 5 9 9 7 10 1 4 12 7 12 11
2 1 2
16 15 4 7 13 10 7 8 5 13 14 2 16 2 1 8 4
2 1 2
20 2 6 6 17 1 9 1 5 3 8 5 13 5 3 2 5 16 8 2 6
16 16 12 6 10 7 10 16 5 11 13 11 11 15 10 13 1
14 5 9 8 9 4 11 4 10 4 10 11 12 6 14
20 10 9 20 19 19 5 4 4 1 13 14 5 1 4 15 8 3 17 10 6
3 2 2 2
10 5 9 10 5 9 9 3 10 6 8
5 1 3 3 3 4
2 1 1
19 8 8 11 16 6 1 4 3 9 2 17 8 7 12 6 13 1 18 11
5 5 3 4 2 3
10 10 1 8 4 1 2 9 1 3 3
4 4 4 1 4
6 3 4 1 3 5 1
18 10 4 13 4 10 17 17 7 10 18 14 4 6 14 12 1 18 12
13 5 12 4 6 3 1 8 8 13 2 5 2 1`

type testCase struct {
	n   int
	arr []int
}

func runCandidate(bin, input string) (string, error) {
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

func solve(arr []int) int64 {
	const limit = 1 << 18
	squares := make([]int, 0, 512)
	for i := 0; i*i < limit; i++ {
		squares = append(squares, i*i)
	}
	freq := make([]int64, limit)
	prefix := 0
	freq[0] = 1
	var bad int64
	for _, v := range arr {
		prefix ^= v
		for _, s := range squares {
			x := prefix ^ s
			if x < limit {
				bad += freq[x]
			}
		}
		freq[prefix]++
	}
	total := int64(len(arr)) * int64(len(arr)+1) / 2
	return total - bad
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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("line %d parse arr[%d]: %v", idx+1, i, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCase{n: n, arr: arr})
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		want := solve(tc.arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(tc.n))
		input.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		gotStr, err := runCandidate(bin, input.String())
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
