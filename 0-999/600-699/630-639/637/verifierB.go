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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func expected(names []string) []string {
	seen := make(map[string]bool)
	res := make([]string, 0, len(names))
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if !seen[name] {
			seen[name] = true
			res = append(res, name)
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "empty testcases file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid test count")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing test case %d\n", caseIdx)
			os.Exit(1)
		}
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil || n != len(fields)-1 {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseIdx)
			os.Exit(1)
		}
		names := fields[1:]
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(names, " "))
		want := strings.Join(expected(names), "\n")
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", caseIdx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
