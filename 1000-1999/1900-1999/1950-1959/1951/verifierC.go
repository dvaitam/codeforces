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

func solve(n int, m, k int64, a []int) int64 {
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return a[idx[i]] < a[idx[j]] })
	x := make([]int64, n)
	left := k
	for _, id := range idx {
		if left == 0 {
			break
		}
		take := m
		if left < take {
			take = left
		}
		x[id] = take
		left -= take
	}
	var prefix, ans int64
	for i := 0; i < n; i++ {
		ans += x[i] * (int64(a[i]) + prefix)
		prefix += x[i]
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idxCase := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idxCase++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		mVal, _ := strconv.ParseInt(fields[1], 10, 64)
		kVal, _ := strconv.ParseInt(fields[2], 10, 64)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[3+i])
			arr[i] = v
		}
		exp := solve(n, mVal, kVal, arr)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, mVal, kVal))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
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
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idxCase, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idxCase, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idxCase)
}
