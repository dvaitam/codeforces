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

type testCaseC struct {
	x   int
	arr []int
}

func parseTestcases(path string) ([]testCaseC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseC, T)
	for i := 0; i < T; i++ {
		var n, x int
		fmt.Fscan(in, &n, &x)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &arr[j])
		}
		cases[i] = testCaseC{x: x, arr: arr}
	}
	return cases, nil
}

func solveCase(tc testCaseC) int {
	a := append([]int(nil), tc.arr...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	teams := 0
	size := 0
	for _, v := range a {
		size++
		if size*v >= tc.x {
			teams++
			size = 0
		}
	}
	return teams
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.arr), tc.x))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected := solveCase(tc)
		got, err := strconv.Atoi(strings.TrimSpace(outStr))
		if err != nil || got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, expected, strings.TrimSpace(outStr))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
