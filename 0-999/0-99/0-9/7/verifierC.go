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

func extGCD(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g2, x1, y1 := extGCD(b, a%b)
	return g2, y1, x1 - (a/b)*y1
}

func solve(A, B, C int64) string {
	g, x0, y0 := extGCD(abs(A), abs(B))
	if A < 0 {
		x0 = -x0
	}
	if B < 0 {
		y0 = -y0
	}
	if (-C)%g != 0 {
		return "-1"
	}
	factor := (-C) / g
	x := x0 * factor
	y := y0 * factor
	return fmt.Sprintf("%d %d", x, y)
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
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
		if len(parts) != 3 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		A, _ := strconv.ParseInt(parts[0], 10, 64)
		B, _ := strconv.ParseInt(parts[1], 10, 64)
		C, _ := strconv.ParseInt(parts[2], 10, 64)
		exp := solve(A, B, C)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", A, B, C)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != exp {
			fmt.Printf("Test %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
