package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(line string) string {
	s := strings.TrimSpace(line)
	n := len(s)
	target := make([]int, n)
	for i, c := range s {
		if c == 'A' {
			target[i] = 0
		} else {
			target[i] = 1
		}
	}
	count := 0
	for mask := 0; mask < (1 << n); mask++ {
		state := make([]int, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				state[i] = 1
			}
		}
		next := make([]int, n)
		for i := 0; i < n; i++ {
			left := (i - 1 + n) % n
			right := (i + 1) % n
			flip := (state[i] == 0 && state[right] == 1) || (state[left] == 0 && state[i] == 1)
			if flip {
				next[i] = 1 - state[i]
			} else {
				next[i] = state[i]
			}
		}
		ok := true
		for i := 0; i < n; i++ {
			if next[i] != target[i] {
				ok = false
				break
			}
		}
		if ok {
			count++
		}
	}
	return fmt.Sprintf("%d", count)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
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
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
