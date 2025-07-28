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

func solve(k int) []int {
	n := 1 << k
	res := make([]int, k)
	for t := 1; t <= k; t++ {
		res[t-1] = n - (1 << t) + 1
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesH.txt")
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
		kVal, _ := strconv.Atoi(fields[0])
		n := 1 << kVal
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			arr[i] = v
		}
		expSlice := solve(kVal)
		exp := strings.TrimSpace(strings.Join(func() []string {
			s := make([]string, len(expSlice))
			for i, v := range expSlice {
				s[i] = strconv.Itoa(v)
			}
			return s
		}(), " "))

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(strconv.Itoa(kVal))
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')

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
