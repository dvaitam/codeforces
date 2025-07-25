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

type testCase struct {
	n   int
	arr []int
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func loadTests() ([]testCase, error) {
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("not enough data")
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("not enough numbers for case %d", i+1)
			}
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		tests[i] = testCase{n, arr}
	}
	return tests, nil
}

func expected(t testCase) string {
	even := 0
	var sb strings.Builder
	for i, v := range t.arr {
		if v%2 == 0 {
			even++
		}
		if even%2 == 1 {
			sb.WriteString("1")
		} else {
			sb.WriteString("2")
		}
		if i < len(t.arr)-1 {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	tests, err := loadTests()
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
		input := b.String()
		exp := expected(tc)
		out, err := runBinary(exe, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\ninput:%sexpected:%s\ngot:%s", idx+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
