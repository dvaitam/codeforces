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

type pair struct{ u, v int }

func expectedOps(arr []int) []pair {
	n := len(arr)
	var ops []pair
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if arr[i] > arr[j] {
				ops = append(ops, pair{i, j})
			}
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		ai, aj := arr[ops[i].u], arr[ops[j].u]
		if ai != aj {
			return ai < aj
		}
		return ops[i].v > ops[j].v
	})
	return ops
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesE.txt")
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Printf("test %d: wrong values\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			arr[i] = v
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", v))
		}
		input.WriteByte('\n')

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
		if len(outLines) == 0 {
			fmt.Printf("Test %d: empty output\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(strings.TrimSpace(outLines[0]))
		if err != nil {
			fmt.Printf("Test %d: invalid count\n", idx)
			os.Exit(1)
		}
		if m != len(outLines)-1 {
			fmt.Printf("Test %d: expected %d lines of pairs got %d\n", idx, m, len(outLines)-1)
			os.Exit(1)
		}
		expect := expectedOps(arr)
		if len(expect) != m {
			fmt.Printf("Test %d: expected %d operations but got %d\n", idx, len(expect), m)
			os.Exit(1)
		}
		for i, line2 := range outLines[1:] {
			fields := strings.Fields(line2)
			if len(fields) != 2 {
				fmt.Printf("Test %d line %d: expected two ints\n", idx, i+1)
				os.Exit(1)
			}
			u, _ := strconv.Atoi(fields[0])
			v, _ := strconv.Atoi(fields[1])
			exp := expect[i]
			if u != exp.u+1 || v != exp.v+1 {
				fmt.Printf("Test %d: mismatch at op %d\n", idx, i+1)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
