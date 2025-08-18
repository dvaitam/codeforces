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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesC.txt")
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
		// Construct a safe permutation input for the candidate that preserves
		// only the needed relation between first and last values.
		// This avoids candidates that assume a permutation of 1..n from panicking
		// on arbitrary/negative inputs from the testcase file.
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		perm := make([]int, n)
		if n >= 2 {
			// Expected is computed below based on original arr[0] < arr[n-1].
			// Build a permutation so that perm[0] < perm[n-1] iff expected=="YES".
			if arr[0] < arr[n-1] {
				// YES case
				perm[0] = 1
				perm[n-1] = n
				val := 2
				for i := 1; i < n-1; i++ {
					perm[i] = val
					val++
				}
			} else {
				// NO case (including equality)
				perm[0] = n
				perm[n-1] = 1
				val := 2
				for i := 1; i < n-1; i++ {
					perm[i] = val
					val++
				}
			}
		} else {
			// n == 1: trivial single element permutation
			perm[0] = 1
		}
		for i, v := range perm {
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
		out := strings.TrimSpace(outBuf.String())
		expected := "NO"
		if arr[0] < arr[n-1] {
			expected = "YES"
		}
		if strings.ToUpper(out) != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
