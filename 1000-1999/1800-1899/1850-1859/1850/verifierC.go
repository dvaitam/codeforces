package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solveGrid(lines []string) string {
	var sb strings.Builder
	for _, line := range lines {
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch != '.' {
				sb.WriteByte(ch)
			}
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputs []string
	var exps []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 8 {
			fmt.Fprintf(os.Stderr, "invalid testcase line: %s\n", line)
			os.Exit(1)
		}
		grid := parts
		exp := solveGrid(grid)
		var sb strings.Builder
		sb.WriteString("1\n")
		for i := 0; i < 8; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		inputs = append(inputs, sb.String())
		exps = append(exps, exp)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	for idx, input := range inputs {
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
