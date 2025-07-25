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

func run(bin, input string) (string, error) {
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

func expected(n int, s string) int {
	zeros := 0
	for _, c := range s {
		if c == '0' {
			zeros++
		}
	}
	ones := n - zeros
	if zeros > ones {
		return zeros - ones
	}
	return ones - zeros
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("failed to open testcasesA.txt:", err)
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
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("bad test format on line %d\n", idx+1)
			os.Exit(1)
		}
		idx++
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Printf("bad n on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		s := fields[1]
		input := fmt.Sprintf("%d\n%s\n", n, s)
		exp := expected(n, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", idx, exp, got, input)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
