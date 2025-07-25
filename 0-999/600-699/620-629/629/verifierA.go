package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func solveA(n int, grid []string) int {
	ans := 0
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < n; j++ {
			if grid[i][j] == 'C' {
				cnt++
			}
		}
		ans += cnt * (cnt - 1) / 2
	}
	for j := 0; j < n; j++ {
		cnt := 0
		for i := 0; i < n; i++ {
			if grid[i][j] == 'C' {
				cnt++
			}
		}
		ans += cnt * (cnt - 1) / 2
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(100) + 1
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			row := make([]byte, n)
			for j := 0; j < n; j++ {
				if rand.Intn(2) == 0 {
					row[j] = '.'
				} else {
					row[j] = 'C'
				}
			}
			grid[i] = string(row)
		}
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for i := 0; i < n; i++ {
			fmt.Fprintln(&input, grid[i])
		}
		expected := solveA(n, grid)
		cmd := exec.Command(binary)
		cmd.Stdin = &input
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: binary error: %v\n", t, err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(&out)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output\n", t)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
