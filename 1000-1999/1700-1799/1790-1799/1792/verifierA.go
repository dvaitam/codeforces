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

func solveCaseA(n int, arr []int) int {
	ones := 0
	for _, x := range arr {
		if x == 1 {
			ones++
		}
	}
	return n - ones/2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(100) + 1
		fmt.Fprintln(&input, n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(100) + 1
		}
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, arr[j])
		}
		input.WriteByte('\n')
		expected[i] = solveCaseA(n, arr)
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
