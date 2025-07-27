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
	arr []int
}

func parseTestcases(path string) ([]testCaseA, error) {
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
	cases := make([]testCaseA, T)
	for i := 0; i < T; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, err
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &arr[j])
		}
		cases[i] = testCaseA{arr: arr}
	}
	return cases, nil
}

func solveCase(arr []int) (bool, [3]int) {
	stack := []int{}
	found := false
	var x, y, z int
	for i, v := range arr {
		for len(stack) > 0 && arr[stack[len(stack)-1]] > v {
			if len(stack) > 1 && !found {
				x = stack[len(stack)-2] + 1
				y = stack[len(stack)-1] + 1
				z = i + 1
				found = true
			}
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	if found {
		return true, [3]int{x, y, z}
	}
	return false, [3]int{}
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

func checkOutput(arr []int, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("no output")
	}
	lines := strings.Split(output, "\n")
	ok, _ := solveCase(arr)
	first := strings.TrimSpace(lines[0])
	if strings.EqualFold(first, "NO") {
		if ok {
			return fmt.Errorf("expected YES but got NO")
		}
		if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if !strings.EqualFold(first, "YES") {
		return fmt.Errorf("output should start with YES or NO")
	}
	if !ok {
		return fmt.Errorf("expected NO but got YES")
	}
	if len(lines) < 2 {
		return fmt.Errorf("missing indices line")
	}
	fields := strings.Fields(lines[1])
	if len(fields) != 3 {
		return fmt.Errorf("expected 3 indices")
	}
	idx := make([]int, 3)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		idx[i] = v
	}
	n := len(arr)
	i, j, k := idx[0], idx[1], idx[2]
	if !(1 <= i && i < j && j < k && k <= n) {
		return fmt.Errorf("indices out of order or range")
	}
	if !(arr[i-1] < arr[j-1] && arr[j-1] > arr[k-1]) {
		return fmt.Errorf("indices do not satisfy condition")
	}
	if len(lines) > 2 && strings.TrimSpace(lines[2]) != "" {
		return fmt.Errorf("extra output")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
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
		if e := checkOutput(tc.arr, outStr); e != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, e)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
