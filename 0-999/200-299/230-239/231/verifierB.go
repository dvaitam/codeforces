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

type TestCaseB struct {
	n   int
	d   int
	l   int
	ans string
}

func solveCaseB(n, d, l int) string {
	arr := make([]int, n)
	m := d
	k := l
	for i := 0; i < n; i++ {
		if m > 0 {
			arr[i] = k
		} else {
			arr[i] = 1
		}
		m = arr[i] - m
	}
	arr[n-1] -= m
	if arr[n-1] > 0 && arr[n-1] <= k {
		var sb strings.Builder
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		return sb.String()
	}
	return "-1"
}

func readCasesB(path string) ([]TestCaseB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]TestCaseB, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		l, _ := strconv.Atoi(scan.Text())
		ans := solveCaseB(n, d, l)
		cases[i] = TestCaseB{n, d, l, ans}
	}
	return cases, nil
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := readCasesB("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.d, tc.l)
		expected := tc.ans
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
