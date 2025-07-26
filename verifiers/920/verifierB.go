package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveB(l, r []int) []int {
	n := len(l)
	ans := make([]int, n)
	cur := 1
	for i := 0; i < n; i++ {
		if cur < l[i] {
			cur = l[i]
		}
		if cur > r[i] {
			ans[i] = 0
		} else {
			ans[i] = cur
			cur++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(2)
	const tests = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, tests)
	expectedAll := make([][]int, tests)
	for t := 0; t < tests; t++ {
		n := rand.Intn(20) + 1
		fmt.Fprintln(&input, n)
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			l[i] = rand.Intn(30) + 1
			r[i] = l[i] + rand.Intn(10)
			fmt.Fprintf(&input, "%d %d\n", l[i], r[i])
		}
		expectedAll[t] = solveB(l, r)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to run binary:", err)
		os.Exit(1)
	}
	parts := strings.Fields(string(out))
	// gather expected output count
	expectedCount := 0
	for _, arr := range expectedAll {
		expectedCount += len(arr)
	}
	if len(parts) != expectedCount {
		fmt.Printf("Expected %d numbers, got %d\n", expectedCount, len(parts))
		os.Exit(1)
	}
	idx := 0
	for t, arr := range expectedAll {
		for i, exp := range arr {
			if idx >= len(parts) {
				fmt.Println("Output ended early")
				os.Exit(1)
			}
			got, err := strconv.Atoi(parts[idx])
			if err != nil || got != exp {
				fmt.Printf("Test %d position %d failed: expected %d got %s\n", t+1, i+1, exp, parts[idx])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}
