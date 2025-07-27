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

func mex(arr []int) int {
	present := make(map[int]bool)
	for _, v := range arr {
		present[v] = true
	}
	m := 0
	for {
		if !present[m] {
			return m
		}
		m++
	}
}

func isSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			return false
		}
	}
	return true
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[i+1])
			arr[i] = v
		}
		var input strings.Builder
		input.WriteString("1\n")
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
		outFields := strings.Fields(strings.TrimSpace(outBuf.String()))
		if len(outFields) == 0 {
			fmt.Printf("Test %d: empty output\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(outFields[0])
		if err != nil {
			fmt.Printf("Test %d: invalid operations count\n", idx)
			os.Exit(1)
		}
		if len(outFields)-1 != m {
			fmt.Printf("Test %d: expected %d positions got %d\n", idx, m, len(outFields)-1)
			os.Exit(1)
		}
		ops := make([]int, m)
		for i := 0; i < m; i++ {
			pos, err := strconv.Atoi(outFields[i+1])
			if err != nil || pos < 1 || pos > n {
				fmt.Printf("Test %d: invalid position %q\n", idx, outFields[i+1])
				os.Exit(1)
			}
			ops[i] = pos - 1
		}
		a := append([]int(nil), arr...)
		for _, p := range ops {
			m := mex(a)
			a[p] = m
		}
		for i := 0; i < n; i++ {
			if a[i] != i {
				fmt.Printf("Test %d failed: array not sorted after operations\n", idx)
				os.Exit(1)
			}
		}
		if !isSorted(a) {
			fmt.Printf("Test %d failed: array not sorted\n", idx)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
