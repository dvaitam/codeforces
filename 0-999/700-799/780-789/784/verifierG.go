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
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemG.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	passed := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		idx++
		split := strings.SplitN(line, " ", 2)
		if len(split) != 2 {
			fmt.Println("invalid test case:", line)
			continue
		}
		expr := split[0]
		expected := strings.ReplaceAll(split[1], "\\n", "\n")
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(expr + "\n")
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		if err != nil {
			fmt.Printf("Case %d: runtime error: %v\n", idx, err)
			fmt.Printf("Output: %s\n", result)
			continue
		}
		if result == expected {
			passed++
		} else {
			fmt.Printf("Case %d failed: expected %s got %s\n", idx, expected, result)
		}
	}
	fmt.Printf("%d/%d cases passed\n", passed, idx)
}
