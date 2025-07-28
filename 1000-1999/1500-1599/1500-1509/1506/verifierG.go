package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(s string) string {
	last := make([]int, 26)
	for i := range last {
		last[i] = -1
	}
	for i := 0; i < len(s); i++ {
		last[int(s[i]-'a')] = i
	}
	used := make([]bool, 26)
	stack := make([]byte, 0, 26)
	for i := 0; i < len(s); i++ {
		c := s[i]
		idx := int(c - 'a')
		if used[idx] {
			continue
		}
		for len(stack) > 0 && stack[len(stack)-1] < c && last[int(stack[len(stack)-1]-'a')] > i {
			used[int(stack[len(stack)-1]-'a')] = false
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, c)
		used[idx] = true
	}
	return string(stack)
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesG.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesG.txt: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		expect := expected(s)
		input := fmt.Sprintf("1\n%s\n", s)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
