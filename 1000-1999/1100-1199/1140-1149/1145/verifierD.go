package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open problemD.txt:", err)
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
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+2 {
			fmt.Println("invalid test case")
			os.Exit(1)
		}
		arr := strings.Join(parts[1:n+1], " ")
		expected := parts[n+1]
		input := fmt.Sprintf("%d\n%s\n", n, arr)
		var cmd *exec.Cmd
		if strings.HasSuffix(bin, ".go") {
			cmd = exec.Command("go", "run", bin)
		} else {
			cmd = exec.Command(bin)
		}
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		result := strings.TrimSpace(string(out))
		total++
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
