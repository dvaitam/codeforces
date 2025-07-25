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

func apply(arr []int, l, r int) {
	for i := l; i < r; i += 2 {
		arr[i], arr[i+1] = arr[i+1], arr[i]
	}
}

func isSorted(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
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
		if len(parts) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+n {
			fmt.Printf("test %d: wrong number of elements\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[1+i])
			if err != nil {
				fmt.Printf("test %d: invalid number %q\n", idx, parts[1+i])
				os.Exit(1)
			}
			arr[i] = v
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errb.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		tokens := []string{}
		if outStr != "" {
			tokens = strings.Fields(outStr)
		}
		if len(tokens)%2 != 0 {
			fmt.Printf("test %d: odd number of integers in output\n", idx)
			os.Exit(1)
		}
		if len(tokens)/2 > 20000 {
			fmt.Printf("test %d: too many operations (%d)\n", idx, len(tokens)/2)
			os.Exit(1)
		}
		opsArr := append([]int(nil), arr...)
		for i := 0; i < len(tokens); i += 2 {
			l, _ := strconv.Atoi(tokens[i])
			r, _ := strconv.Atoi(tokens[i+1])
			if !(1 <= l && l < r && r <= n) {
				fmt.Printf("test %d: invalid segment %d %d\n", idx, l, r)
				os.Exit(1)
			}
			if (r-l+1)%2 != 0 {
				fmt.Printf("test %d: segment length not even: %d %d\n", idx, l, r)
				os.Exit(1)
			}
			apply(opsArr, l-1, r-1)
		}
		if !isSorted(opsArr) {
			fmt.Printf("test %d failed: array not sorted\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
