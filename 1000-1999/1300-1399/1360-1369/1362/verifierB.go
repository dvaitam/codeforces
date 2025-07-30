package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	arr      []int
	expected int
}

func computeExpectedB(arr []int) int {
	set := make(map[int]bool, len(arr))
	for _, v := range arr {
		set[v] = true
	}
	for k := 1; k < 1024; k++ {
		valid := true
		for _, v := range arr {
			if !set[v^k] {
				valid = false
				break
			}
		}
		if valid {
			return k
		}
	}
	return -1
}

func generateTestsB() []testCaseB {
	const numTests = 100
	rand.Seed(2)
	tests := make([]testCaseB, 0, numTests+5)
	for i := 0; i < numTests; i++ {
		n := rand.Intn(10) + 1
		used := make(map[int]bool)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			for {
				x := rand.Intn(1024)
				if !used[x] {
					used[x] = true
					arr[j] = x
					break
				}
			}
		}
		tests = append(tests, testCaseB{arr: arr, expected: computeExpectedB(arr)})
	}
	edge := []testCaseB{
		{arr: []int{0, 1, 2, 3}, expected: computeExpectedB([]int{0, 1, 2, 3})},
		{arr: []int{1}, expected: computeExpectedB([]int{1})},
		{arr: []int{5, 9, 12}, expected: computeExpectedB([]int{5, 9, 12})},
	}
	tests = append(tests, edge...)
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTestsB()
	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, len(tc.arr))
		for i, v := range tc.arr {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
	}

	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running binary: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(&out)
	for i, tc := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		got, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
