package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(grid []string) string {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cnt := 0
			if grid[i][j] == '#' {
				cnt++
			}
			if grid[i][j+1] == '#' {
				cnt++
			}
			if grid[i+1][j] == '#' {
				cnt++
			}
			if grid[i+1][j+1] == '#' {
				cnt++
			}
			if cnt != 2 {
				return "YES"
			}
		}
	}
	return "NO"
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid []string
	test := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(grid) == 4 {
				test++
				input := strings.Join(grid, "\n") + "\n"
				want := expected(grid)
				got, err := run(bin, input)
				if err != nil {
					fmt.Printf("test %d: %v\n", test, err)
					os.Exit(1)
				}
				if got != want {
					fmt.Printf("test %d failed: expected %s got %s\nInput:\n%s\n", test, want, got, input)
					os.Exit(1)
				}
			}
			grid = nil
			continue
		}
		grid = append(grid, line)
	}
	if len(grid) == 4 {
		test++
		input := strings.Join(grid, "\n") + "\n"
		want := expected(grid)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", test, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed: expected %s got %s\nInput:\n%s\n", test, want, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", test)
}
