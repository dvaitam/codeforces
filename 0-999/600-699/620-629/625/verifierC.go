package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testC struct {
	n, k int
}

func solveC(n, k int) (int, [][]int) {
	ans := make([][]int, n)
	for i := range ans {
		ans[i] = make([]int, n)
	}
	times := 0
	sum := 0
	for c := 0; c < k-1; c++ {
		for r := 0; r < n; r++ {
			times++
			ans[r][c] = times
		}
	}
	for r := 0; r < n; r++ {
		for c := k - 1; c < n; c++ {
			times++
			ans[r][c] = times
			if c == k-1 {
				sum += times
			}
		}
	}
	return sum, ans
}

func genTests() []testC {
	rand.Seed(3)
	tests := make([]testC, 100)
	for i := range tests {
		n := rand.Intn(8) + 1
		k := rand.Intn(n) + 1
		tests[i] = testC{n: n, k: k}
	}
	tests = append(tests, testC{n: 1, k: 1})
	tests = append(tests, testC{n: 4, k: 2})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.k)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		sumExp, matExp := solveC(t.n, t.k)
		reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
		var sumOut int
		if _, err := fmt.Fscan(reader, &sumOut); err != nil {
			fmt.Printf("Test %d failed to read sum: %v\n", i+1, err)
			os.Exit(1)
		}
		if sumOut != sumExp {
			fmt.Printf("Test %d sum mismatch\nInput:%sExpected %d got %d\n", i+1, input, sumExp, sumOut)
			os.Exit(1)
		}
		for r := 0; r < t.n; r++ {
			for c := 0; c < t.n; c++ {
				var v int
				if _, err := fmt.Fscan(reader, &v); err != nil {
					fmt.Printf("Test %d failed reading matrix\n", i+1)
					os.Exit(1)
				}
				if v != matExp[r][c] {
					fmt.Printf("Test %d matrix mismatch at (%d,%d)\n", i+1, r, c)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
