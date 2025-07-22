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

type TestCaseA struct {
	n   int
	tri [][3]int
	ans int
}

func solveCaseA(n int, tri [][3]int) int {
	cnt := 0
	for _, t := range tri {
		if t[0]+t[1]+t[2] >= 2 {
			cnt++
		}
	}
	return cnt
}

func readCasesA(path string) ([]TestCaseA, error) {
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
	cases := make([]TestCaseA, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		n, _ := strconv.Atoi(scan.Text())
		tri := make([][3]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			tri[j] = [3]int{a, b, c}
		}
		cases[i] = TestCaseA{n: n, tri: tri, ans: solveCaseA(n, tri)}
	}
	return cases, nil
}

func runCase(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := readCasesA("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, tc.n)
		for _, t := range tc.tri {
			fmt.Fprintf(&sb, "%d %d %d\n", t[0], t[1], t[2])
		}
		expected := fmt.Sprintf("%d", tc.ans)
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
