package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func expected(n int, path string) int {
	sx, sy := 0, 0
	x, y := 0, 0
	count := 1
	stepCount := 0
	for _, c := range path {
		px, py := x, y
		switch c {
		case 'L':
			x--
		case 'R':
			x++
		case 'U':
			y++
		case 'D':
			y--
		}
		stepCount++
		dist := abs(x-sx) + abs(y-sy)
		if dist < stepCount {
			count++
			sx, sy = px, py
			stepCount = 1
		}
	}
	return count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
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
		var n int
		var path string
		fmt.Sscan(line, &n, &path)
		exp := expected(n, path)
		input := fmt.Sprintf("%d\n%s\n", n, path)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(strings.TrimSpace(string(out)), &got)
		if got != exp {
			fmt.Printf("Test %d failed. Expected %d got %s\n", idx, exp, strings.TrimSpace(string(out)))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
