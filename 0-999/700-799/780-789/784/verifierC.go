package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	total := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Println("invalid test case:", line)
			continue
		}
		n := parts[0]
		nVal := 0
		fmt.Sscan(n, &nVal)
		if len(parts) < nVal+2 {
			fmt.Println("invalid test case:", line)
			continue
		}
		numbers := parts[1 : 1+nVal]
		expected := parts[1+nVal]
		input := n + "\n" + strings.Join(numbers, " ") + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		total++
		result := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", total, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", total, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, total)
}
