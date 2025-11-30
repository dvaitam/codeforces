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

// Embedded testcases from testcasesB.txt.
const embeddedTestcasesB = `100
3 -2 3
6 -4 2
5 0 2
4 7 1
2 0 3
4 -9 1
6 -10 3
2 -5 3
3 6 5
4 4 1
5 6 5
5 -7 1
6 4 3
3 -1 2
5 -1 1
2 -10 3
2 0 3
6 10 3
2 -3 5
3 -2 5
2 2 4
4 8 2
4 0 5
5 4 1
5 1 2
6 -9 2
5 2 1
3 -8 3
2 -5 2
5 -10 3
6 3 1
6 -3 4
4 -2 4
5 5 1
3 6 3
3 -3 2
3 -1 2
3 7 3
6 -10 2
3 3 2
3 7 2
6 -5 2
6 -1 1
5 -4 1
5 -5 4
6 2 2
5 7 5
5 -7 2
2 7 1
2 -10 2
5 4 3
5 1 4
4 -8 2
6 -3 2
2 -9 5
2 4 2
5 -6 5
6 -1 1
6 -6 3
4 1 3
6 10 3
3 -10 4
3 9 4
4 5 3
5 -6 4
4 0 3
6 2 4
6 2 5
4 8 2
3 3 2
4 6 2
6 -6 5
6 7 5
5 0 3
5 -1 4
6 -9 2
4 4 1
2 -4 5
2 -2 1
3 -5 5
3 1 5
4 -10 1
2 9 2
4 -3 4
5 1 3
2 8 4
2 -10 2
2 9 5
4 1 2
5 -4 3
3 -5 4
3 9 5
6 3 5
3 -10 1
5 9 1
3 6 1
5 -1 3
6 -4 1
4 -6 3
4 5 4`

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

func readCasesB() ([]TestCaseB, error) {
	scan := bufio.NewScanner(strings.NewReader(embeddedTestcasesB))
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
	cases, err := readCasesB()
	if err != nil {
		fmt.Println("could not read embedded testcases:", err)
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
