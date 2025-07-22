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

func expectedNavigation(n, p, k int) string {
	start := p - k
	if start < 1 {
		start = 1
	}
	end := p + k
	if end > n {
		end = n
	}
	var tokens []string
	if start > 1 {
		tokens = append(tokens, "<<")
	}
	for i := start; i <= end; i++ {
		if i == p {
			tokens = append(tokens, fmt.Sprintf("(%d)", i))
		} else {
			tokens = append(tokens, strconv.Itoa(i))
		}
	}
	if end < n {
		tokens = append(tokens, ">>")
	}
	return strings.Join(tokens, " ")
}

func runBinary(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
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
		fields := strings.Fields(line)
		if len(fields) != 3 {
			fmt.Fprintf(os.Stderr, "case %d invalid format\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		p, _ := strconv.Atoi(fields[1])
		k, _ := strconv.Atoi(fields[2])
		expect := expectedNavigation(n, p, k)
		input := fmt.Sprintf("%d %d %d\n", n, p, k)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n   got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
