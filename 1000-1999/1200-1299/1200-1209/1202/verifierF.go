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

func expected(a, b int) int {
	n := a + b
	count := 0
	for k := 1; k <= n; k++ {
		m := n / k
		if m == 0 {
			continue
		}
		r := n % k
		for x := 0; x <= r; x++ {
			rem := a - x*(m+1)
			if rem < 0 {
				continue
			}
			if rem%m != 0 {
				continue
			}
			y := rem / m
			if y >= 0 && y <= k-r {
				count++
				break
			}
		}
	}
	return count
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesF.txt:", err)
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
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		input := fmt.Sprintf("%d %d\n", a, b)
		expectedOut := fmt.Sprintf("%d", expected(a, b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expectedOut, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
