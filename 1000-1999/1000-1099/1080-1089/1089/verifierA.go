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
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	caseNum := 0
	passed := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		input := strings.Join(parts[:len(parts)-1], " ")
		expected := parts[len(parts)-1]
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input + "\n")
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		caseNum++
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", caseNum, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", caseNum, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, caseNum)
}
