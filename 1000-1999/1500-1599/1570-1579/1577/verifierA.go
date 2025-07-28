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

func computeSum(nums []int) int {
	s := 0
	for _, v := range nums {
		s += v
	}
	return s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open testcasesA.txt:", err)
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
		if len(parts) == 0 {
			fmt.Fprintf(os.Stderr, "invalid test %d: empty line\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil || len(parts)-1 < n {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			nums[i] = v
		}
		expect := computeSum(nums)

		var in strings.Builder
		in.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			in.WriteString(parts[i+1])
		}
		in.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		output := strings.TrimSpace(outBuf.String())
		if output != fmt.Sprintf("%d", expect) {
			fmt.Printf("test %d failed: expected %d got %s\n", idx, expect, output)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
