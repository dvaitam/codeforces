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

func solve(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	b, err := strconv.Atoi(fields[1])
	if err != nil {
		return "", err
	}
	if len(fields) != 2+n {
		return "", fmt.Errorf("expected %d numbers, got %d", n, len(fields)-2)
	}
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[2+i])
		if err != nil {
			return "", err
		}
		a[i] = v
		sum += v
	}
	V := float64(sum+b) / float64(n)
	for i := 0; i < n; i++ {
		if float64(a[i]) > V {
			return "-1", nil
		}
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		c := V - float64(a[i])
		res[i] = fmt.Sprintf("%.6f", c)
	}
	return strings.Join(res, "\n"), nil
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
		expected, err := solve(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
