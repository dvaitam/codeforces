package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

func solveCaseC(p []int) int {
	n := len(p)
	pre := 0
	for pre < n && p[pre] == pre+1 {
		pre++
	}
	suf := 0
	for suf < n-pre && p[n-1-suf] == n-suf {
		suf++
	}
	rem := n - pre - suf
	if rem < 0 {
		rem = 0
	}
	return (rem + 1) / 2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(44)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		fmt.Fprintln(&input, n)
		p := rand.Perm(n)
		for j := 0; j < n; j++ {
			p[j]++
		}
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, p[j])
		}
		input.WriteByte('\n')
		expected[i] = solveCaseC(p)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "binary execution failed:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(outBytes))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("not enough output on test %d\n", i+1)
			os.Exit(1)
		}
		got, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("invalid integer on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Printf("mismatch on test %d: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
