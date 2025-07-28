package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func expectedCount(s string) int64 {
	var res int64 = 1
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			res += int64(i + 1)
		} else {
			res++
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := ioutil.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcasesB.txt: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "binary execution failed: %v\n", err)
		os.Exit(1)
	}

	inScanner := bufio.NewScanner(bytes.NewReader(data))
	if !inScanner.Scan() {
		fmt.Fprintln(os.Stderr, "testcase file empty")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(inScanner.Text()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad test count")
		os.Exit(1)
	}

	outScanner := bufio.NewScanner(&out)
	outScanner.Split(bufio.ScanWords)

	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !inScanner.Scan() {
			fmt.Fprintf(os.Stderr, "bad testcase format at case %d\n", caseIdx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(inScanner.Text()))
		if !inScanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing string for case %d\n", caseIdx)
			os.Exit(1)
		}
		s := strings.TrimSpace(inScanner.Text())
		if len(s) != n {
			fmt.Fprintf(os.Stderr, "case %d invalid test length\n", caseIdx)
			os.Exit(1)
		}
		if !outScanner.Scan() {
			fmt.Fprintf(os.Stderr, "not enough output for case %d\n", caseIdx)
			os.Exit(1)
		}
		ansStr := outScanner.Text()
		ans, err := strconv.ParseInt(ansStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output not integer\n", caseIdx)
			os.Exit(1)
		}
		exp := expectedCount(s)
		if ans != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong answer got %d expected %d\n", caseIdx, ans, exp)
			os.Exit(1)
		}
	}
	if outScanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
