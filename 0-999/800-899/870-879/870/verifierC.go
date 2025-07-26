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

func solveN(n int64) int64 {
	switch {
	case n < 4:
		return -1
	case n%4 == 0:
		return n / 4
	case n%4 == 1:
		if n >= 9 {
			return (n-9)/4 + 1
		}
		return -1
	case n%4 == 2:
		if n >= 6 {
			return (n-6)/4 + 1
		}
		return -1
	default:
		if n >= 15 {
			return (n-15)/4 + 2
		}
		return -1
	}
}

func runCase(exe string, input, expect string) error {
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, q)
		for i := 0; i < q; i++ {
			scan.Scan()
			val, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[i] = val
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "%d\n", q)
		for i, v := range arr {
			fmt.Fprintf(inputBuilder, "%d\n", v)
			_ = i
		}
		expectBuilder := &strings.Builder{}
		for _, v := range arr {
			fmt.Fprintf(expectBuilder, "%d\n", solveN(v))
		}
		if err := runCase(exe, inputBuilder.String(), expectBuilder.String()); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
