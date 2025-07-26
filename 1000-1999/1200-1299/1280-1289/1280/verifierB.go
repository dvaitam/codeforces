package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func rowAllA(row string) bool {
	for i := 0; i < len(row); i++ {
		if row[i] != 'A' {
			return false
		}
	}
	return true
}

func colAllA(grid []string, col int) bool {
	for i := 0; i < len(grid); i++ {
		if grid[i][col] != 'A' {
			return false
		}
	}
	return true
}

func solve(grid []string) int {
	r := len(grid)
	c := len(grid[0])
	hasA := false
	allA := true
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if grid[i][j] == 'A' {
				hasA = true
			} else {
				allA = false
			}
		}
	}
	if !hasA {
		return -1
	}
	if allA {
		return 0
	}
	if rowAllA(grid[0]) || rowAllA(grid[r-1]) || colAllA(grid, 0) || colAllA(grid, c-1) {
		return 1
	}
	if grid[0][0] == 'A' || grid[0][c-1] == 'A' || grid[r-1][0] == 'A' || grid[r-1][c-1] == 'A' {
		return 2
	}
	for i := 0; i < r; i++ {
		if rowAllA(grid[i]) {
			return 2
		}
	}
	for j := 0; j < c; j++ {
		if colAllA(grid, j) {
			return 2
		}
	}
	for j := 0; j < c; j++ {
		if grid[0][j] == 'A' || grid[r-1][j] == 'A' {
			return 3
		}
	}
	for i := 0; i < r; i++ {
		if grid[i][0] == 'A' || grid[i][c-1] == 'A' {
			return 3
		}
	}
	return 4
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	t := 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		r := rand.Intn(4) + 1
		c := rand.Intn(4) + 1
		grid := make([]string, r)
		for j := 0; j < r; j++ {
			row := make([]byte, c)
			for k := 0; k < c; k++ {
				if rand.Intn(2) == 0 {
					row[k] = 'A'
				} else {
					row[k] = 'P'
				}
			}
			grid[j] = string(row)
		}
		fmt.Fprintln(&input, r, c)
		for _, row := range grid {
			fmt.Fprintln(&input, row)
		}
		res := solve(grid)
		if res == -1 {
			expected[i] = "MORTAL"
		} else {
			expected[i] = fmt.Sprintf("%d", res)
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary error:", err)
		fmt.Print(string(out))
		return
	}
	outputs := strings.Fields(string(out))
	if len(outputs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outputs))
		fmt.Print(string(out))
		return
	}
	for i := 0; i < t; i++ {
		if outputs[i] != expected[i] {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected[i], outputs[i])
			return
		}
	}
	fmt.Println("All tests passed!")
}
