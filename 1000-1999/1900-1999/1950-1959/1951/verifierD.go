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

func solve(n, k int64) string {
	if k > n {
		return "NO"
	}
	if k == n {
		return "YES\n1\n1"
	}
	if k <= (n+1)/2 {
		return fmt.Sprintf("YES\n2\n%d %d", n-k+1, 1)
	}
	return "NO"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		nVal, _ := strconv.ParseInt(fields[0], 10, 64)
		kVal, _ := strconv.ParseInt(fields[1], 10, 64)
		exp := solve(nVal, kVal)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", nVal, kVal))

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != exp {
			fmt.Printf("Test %d failed: expected %q got %q\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
