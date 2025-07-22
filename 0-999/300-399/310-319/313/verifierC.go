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

func solve(arr []int64) int64 {
	n := len(arr)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var sumAll int64
	for _, v := range arr {
		sumAll += v
	}
	m := (n - 1) / 3
	var sumMax int64
	for i := n - m; i < n; i++ {
		sumMax += arr[i]
	}
	return sumAll + sumMax
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	scan.Scan()
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[i] = v
		}
		expected := fmt.Sprintf("%d", solve(arr))
		var input strings.Builder
		input.WriteString(strconv.Itoa(n))
		input.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(v, 10))
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
