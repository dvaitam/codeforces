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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(grid []string) string {
	score := 0
	for i := 0; i < 10; i++ {
		row := grid[i]
		for j := 0; j < 10 && j < len(row); j++ {
			if row[j] == 'X' {
				val := min(min(i, 9-i), min(j, 9-j)) + 1
				score += val
			}
		}
	}
	return fmt.Sprintf("%d\n", score)
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		grid := make([]string, 10)
		for i := 0; i < 10; i++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			grid[i] = scan.Text()
		}
		var input strings.Builder
		input.WriteString("1\n")
		for i := 0; i < 10; i++ {
			input.WriteString(grid[i])
			input.WriteByte('\n')
		}
		exp := expected(grid)
		if err := runCase(exe, input.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
