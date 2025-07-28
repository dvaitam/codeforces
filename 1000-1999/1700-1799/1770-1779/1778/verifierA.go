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

type testCaseA struct {
	n   int
	arr []int
}

func parseCases(path string) ([]testCaseA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseA, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &arr[j]); err != nil {
				return nil, err
			}
		}
		cases[i] = testCaseA{n: n, arr: arr}
	}
	return cases, nil
}

func solve(tc testCaseA) int {
	sum := 0
	for _, v := range tc.arr {
		sum += v
	}
	hasNegPair := false
	hasOppPair := false
	for i := 0; i+1 < tc.n; i++ {
		if tc.arr[i] == -1 && tc.arr[i+1] == -1 {
			hasNegPair = true
		}
		if tc.arr[i] != tc.arr[i+1] {
			hasOppPair = true
		}
	}
	if hasNegPair {
		return sum + 4
	}
	if hasOppPair {
		return sum
	}
	return sum - 4
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		// build input with t=1
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "1\n%d\n", tc.n)
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solve(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
