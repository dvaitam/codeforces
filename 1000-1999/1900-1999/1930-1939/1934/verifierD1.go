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
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD1.go <binary_path>")
		os.Exit(1)
	}
	binPath := os.Args[1]

	if strings.HasSuffix(binPath, ".go") {
		tmpBin := "./verifier_tmp_sol"
		cmd := exec.Command("go", "build", "-o", tmpBin, binPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Failed to compile %s: %v\n", binPath, err)
			os.Exit(1)
		}
		defer os.Remove(tmpBin)
		binPath = tmpBin
	}

	rand.Seed(time.Now().UnixNano())
	const numTests = 100

	var inputBuf bytes.Buffer
	fmt.Fprintf(&inputBuf, "%d\n", numTests)
	
	type TestCase struct {
		n, m uint64
	}
	var cases []TestCase

	for i := 0; i < numTests; i++ {
		// Generate n, m
		// n up to 10^18.
		// Ensure n >= 2 so m can be >= 1.
		n := uint64(rand.Int63n(1e18)) + 2
		m := uint64(rand.Int63n(int64(n-1))) + 1
		cases = append(cases, TestCase{n, m})
		fmt.Fprintf(&inputBuf, "%d %d\n", n, m)
	}

	cmd := exec.Command(binPath)
	cmd.Stdin = &inputBuf
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Runtime error: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(&outBuf)
	scanner.Split(bufio.ScanWords)

	for i, c := range cases {
		n, m := c.n, c.m
		if !scanner.Scan() {
			fmt.Printf("Test %d: missing output\n", i+1)
			os.Exit(1)
		}
		kStr := scanner.Text()
		k, err := strconv.Atoi(kStr)
		if err != nil {
			fmt.Printf("Test %d: invalid k %q\n", i+1, kStr)
			os.Exit(1)
		}

		if k == -1 {
			continue
		}

		if k < 1 || k > 63 {
			fmt.Printf("Test %d: k out of range %d\n", i+1, k)
			os.Exit(1)
		}

		path := make([]uint64, 0, k+1)
		for j := 0; j < k+1; j++ {
			if !scanner.Scan() {
				fmt.Printf("Test %d: missing path element %d\n", i+1, j)
				os.Exit(1)
			}
			valStr := scanner.Text()
			val, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				fmt.Printf("Test %d: invalid path element %q\n", i+1, valStr)
				os.Exit(1)
			}
			path = append(path, val)
		}

		if len(path) != k+1 {
			fmt.Printf("Test %d: expected %d elements, got %d\n", i+1, k+1, len(path))
			os.Exit(1)
		}
		if path[0] != n {
			fmt.Printf("Test %d: path start %d != n %d\n", i+1, path[0], n)
			os.Exit(1)
		}
		if path[k] != m {
			fmt.Printf("Test %d: path end %d != m %d\n", i+1, path[k], m)
			os.Exit(1)
		}

		for j := 0; j < k; j++ {
			curr := path[j]
			next := path[j+1]

			if next >= curr {
				fmt.Printf("Test %d: step %d -> %d not strictly decreasing\n", i+1, curr, next)
				os.Exit(1)
			}
			if (curr ^ next) >= curr {
				fmt.Printf("Test %d: step %d -> %d XOR condition failed (curr^next >= curr)\n", i+1, curr, next)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed!\n", numTests)
}