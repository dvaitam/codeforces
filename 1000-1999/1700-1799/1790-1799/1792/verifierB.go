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

func solveCaseB(a1, a2, a3, a4 int) int {
	if a1 == 0 {
		return 1
	}
	pair := a2
	if a3 < pair {
		pair = a3
	}
	ans := a1 + 2*pair
	diff := a2 - a3
	if diff < 0 {
		diff = -diff
	}
	extra := diff + a4
	if a1+1 < extra {
		ans += a1 + 1
	} else {
		ans += extra
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(43)
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]int, t)
	for i := 0; i < t; i++ {
		a1 := rand.Intn(10)
		a2 := rand.Intn(10)
		a3 := rand.Intn(10)
		a4 := rand.Intn(10)
		fmt.Fprintf(&input, "%d %d %d %d\n", a1, a2, a3, a4)
		expected[i] = solveCaseB(a1, a2, a3, a4)
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
