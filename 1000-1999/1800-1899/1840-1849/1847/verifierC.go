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

func solveC(n int, arr []int) int {
	prefix := 0
	seen := make([]bool, 256)
	seen[0] = true
	for _, x := range arr {
		prefix ^= x
		seen[prefix] = true
	}
	ans := 0
	for i := 0; i < 256; i++ {
		if !seen[i] {
			continue
		}
		for j := 0; j < 256; j++ {
			if seen[j] {
				v := i ^ j
				if v > ans {
					ans = v
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(256)
		}
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			if j+1 == n {
				fmt.Fprintf(&input, "%d\n", arr[j])
			} else {
				fmt.Fprintf(&input, "%d ", arr[j])
			}
		}
		expected[i] = strconv.Itoa(solveC(n, arr))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error executing binary: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(string(bytes.TrimSpace(out)))
	if len(lines) != t {
		fmt.Fprintf(os.Stderr, "expected %d lines of output, got %d\n", t, len(lines))
		fmt.Fprint(os.Stderr, string(out))
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if lines[i] != expected[i] {
			fmt.Fprintf(os.Stderr, "mismatch on case %d: expected %s got %s\n", i+1, expected[i], lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
