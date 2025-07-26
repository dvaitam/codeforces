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

func countCompromised(arr []int) int64 {
	n := len(arr)
	var ans int64
	for l := 0; l < n; l++ {
		orVal := 0
		for r := l; r < n; r++ {
			orVal |= arr[r]
			present := false
			for k := l; k <= r; k++ {
				if arr[k] == orVal {
					present = true
					break
				}
			}
			if present {
				ans++
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesA.txt")
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
		if len(fields) == 0 {
			fmt.Printf("invalid test case %d: empty line\n", idx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("invalid n at test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(fields)-1 < n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[1+i])
			if err != nil {
				fmt.Printf("test %d: invalid number %q\n", idx, fields[1+i])
				os.Exit(1)
			}
			arr[i] = v
		}
		expected := countCompromised(arr)

		// Prepare input for binary
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
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
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
