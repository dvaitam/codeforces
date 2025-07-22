package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, m int) int {
	small, big := n, m
	if small > big {
		small, big = big, small
	}
	switch small {
	case 1:
		return big
	case 2:
		blocks := big / 4
		rem := big % 4
		return blocks*4 + min(rem*2, 4)
	default:
		return (n*m + 1) / 2
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
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
		if _, err := fmt.Sscan(line, &n, &m); err != nil {
			fmt.Printf("bad test case on line %d\n", idx)
			os.Exit(1)
		}
		expect := solve(n, m)
		input := fmt.Sprintf("%d %d\n", n, m)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
