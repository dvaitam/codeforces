package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
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

func parseCase(scanner *bufio.Scanner) (string, bool) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var n, m int
		fmt.Sscan(line, &n, &m)
		var sb strings.Builder
		sb.WriteString(line + "\n")
		if !scanner.Scan() {
			return "", false
		}
		sb.WriteString(strings.TrimSpace(scanner.Text()) + "\n")
		for i := 0; i < m; i++ {
			if !scanner.Scan() {
				return "", false
			}
			sb.WriteString(strings.TrimSpace(scanner.Text()) + "\n")
		}
		// consume blank line
		scanner.Scan()
		return sb.String(), true
	}
	return "", false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		input, ok := parseCase(scanner)
		if !ok {
			break
		}
		idx++
		want, err := runProg("266E.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected:\n%s\nGot:\n%s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
