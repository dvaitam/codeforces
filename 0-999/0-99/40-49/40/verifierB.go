package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func countBlack(n, m int, k int64) int64 {
	nn := int64(n) - 2*k
	mm := int64(m) - 2*k
	if nn <= 0 || mm <= 0 {
		return 0
	}
	area := nn * mm
	return (area + 1) / 2
}

func expected(n, m int, x int64) int64 {
	k := x - 1
	b0 := countBlack(n, m, k)
	b1 := countBlack(n, m, k+1)
	return b0 - b1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, m int
		var x int64
		fmt.Sscan(line, &n, &m, &x)
		exp := expected(n, m, x)
		input := fmt.Sprintf("%d %d\n%d\n", n, m, x)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
