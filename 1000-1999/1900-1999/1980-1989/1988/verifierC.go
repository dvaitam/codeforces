package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func genTest() uint64 {
	// Problem constraints: n <= 10^18
	maxN := int64(1e18)
	// Generate n with random number of bits set
	n := uint64(rand.Int63n(maxN)) + 1
	return n
}

func check(n uint64, seq []uint64) error {
	// Expected length: popcount(n) + 1, unless popcount(n) == 1 then 1
	ones := bits.OnesCount64(n)
	expectedLen := ones + 1
	if ones == 1 { // power of 2
		expectedLen = 1
	}

	if len(seq) != expectedLen {
		return fmt.Errorf("length mismatch: expected %d, got %d", expectedLen, len(seq))
	}

	if len(seq) == 0 {
		return fmt.Errorf("empty sequence")
	}

	// Last element must be n
	if seq[len(seq)-1] != n {
		return fmt.Errorf("last element %d != n (%d)", seq[len(seq)-1], n)
	}

	// Strictly increasing
	for i := 0; i < len(seq)-1; i++ {
		if seq[i] >= seq[i+1] {
			return fmt.Errorf("not strictly increasing: seq[%d]=%d >= seq[%d]=%d", i, seq[i], i+1, seq[i+1])
		}
	}

	// OR condition: a_i | a_{i-1} == n
	for i := 1; i < len(seq); i++ {
		if (seq[i] | seq[i-1]) != n {
			return fmt.Errorf("OR condition failed at index %d: %d | %d = %d != %d", i, seq[i], seq[i-1], seq[i]|seq[i-1], n)
		}
	}
	
	// Elements <= n and > 0
	for i, v := range seq {
		if v > n {
			return fmt.Errorf("element %d at index %d > n", v, i)
		}
		if v == 0 {
			return fmt.Errorf("element at index %d is 0", i)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go <binary_path>")
		os.Exit(1)
	}
	binPath := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	const NumTests = 100
	var inputBuf bytes.Buffer
	var testCases []uint64

	fmt.Fprintf(&inputBuf, "%d\n", NumTests)
	for i := 0; i < NumTests; i++ {
		n := genTest()
		testCases = append(testCases, n)
		fmt.Fprintf(&inputBuf, "%d\n", n)
	}

	var cmd *exec.Cmd
	if strings.HasSuffix(binPath, ".go") {
		cmd = exec.Command("go", "run", binPath)
	} else {
		cmd = exec.Command(binPath)
	}

	cmd.Stdin = &inputBuf
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Runtime error executing solution: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(&outBuf)
	scanner.Split(bufio.ScanWords)

	for i, n := range testCases {
		if !scanner.Scan() {
			fmt.Printf("Test %d: missing length output\n", i+1)
			os.Exit(1)
		}
		kStr := scanner.Text()
		k, err := strconv.Atoi(kStr)
		if err != nil {
			fmt.Printf("Test %d: invalid length %q\n", i+1, kStr)
			os.Exit(1)
		}

		seq := make([]uint64, 0, k)
		for j := 0; j < k; j++ {
			if !scanner.Scan() {
				fmt.Printf("Test %d: missing element %d\n", i+1, j+1)
				os.Exit(1)
			}
			valStr := scanner.Text()
			val, err := strconv.ParseUint(valStr, 10, 64)
			if err != nil {
				fmt.Printf("Test %d: invalid element %q\n", i+1, valStr)
				os.Exit(1)
			}
			seq = append(seq, val)
		}

		if err := check(n, seq); err != nil {
			fmt.Printf("Test %d (n=%d) failed: %v\n", i+1, n, err)
			fmt.Printf("Output: %v\n", seq)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed!\n", NumTests)
}