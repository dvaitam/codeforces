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

func solve(n, m int, aDigits, bDigits []int) string {
	a := [10]bool{}
	b := [10]bool{}
	minA, minB := 10, 10
	for _, v := range aDigits {
		a[v] = true
		if v < minA {
			minA = v
		}
	}
	for _, v := range bDigits {
		b[v] = true
		if v < minB {
			minB = v
		}
	}
	common := 10
	for d := 1; d <= 9; d++ {
		if a[d] && b[d] && d < common {
			common = d
		}
	}
	if common < 10 {
		return fmt.Sprintf("%d\n", common)
	}
	if minA < minB {
		return fmt.Sprintf("%d%d\n", minA, minB)
	}
	return fmt.Sprintf("%d%d\n", minB, minA)
}

func runCase(exe string, input string, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		aDigits := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			aDigits[i], _ = strconv.Atoi(scan.Text())
		}
		bDigits := make([]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			bDigits[i], _ = strconv.Atoi(scan.Text())
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "%d %d\n", n, m)
		for i, v := range aDigits {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			fmt.Fprintf(inputBuilder, "%d", v)
		}
		inputBuilder.WriteByte('\n')
		for i, v := range bDigits {
			if i > 0 {
				inputBuilder.WriteByte(' ')
			}
			fmt.Fprintf(inputBuilder, "%d", v)
		}
		inputBuilder.WriteByte('\n')
		expect := solve(n, m, aDigits, bDigits)
		if err := runCase(exe, inputBuilder.String(), expect); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
