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

func parseTestsA() ([]testCaseA, error) {
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var tests []testCaseA
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("invalid n: %v", err)
		}
		if len(fields)-1 < n {
			return nil, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("invalid number %q", fields[1+i])
			}
			arr[i] = v
		}
		tests = append(tests, testCaseA{arr})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func expectedA(arr []int) string {
	maxVal := 0
	for _, v := range arr {
		if v > maxVal {
			maxVal = v
		}
	}
	freq := make([]int, maxVal+1)
	for _, v := range arr {
		freq[v]++
	}
	for i := 1; i <= maxVal; i++ {
		if freq[i] > freq[i-1] {
			return "NO"
		}
	}
	return "YES"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestsA()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range tests {
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		expect := expectedA(tc.arr)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(out)) != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
