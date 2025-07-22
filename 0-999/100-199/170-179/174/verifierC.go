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

func solveCase(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return "", fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return "", err
	}
	if len(fields) != n+1 {
		return "", fmt.Errorf("expected %d numbers, got %d", n, len(fields)-1)
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return "", err
		}
		a[i] = v
	}
	var starts, ends []int
	s := 0
	for i := 1; i <= n+1; i++ {
		x := 0
		if i <= n {
			x = a[i-1]
		}
		for s < x {
			starts = append(starts, i)
			s++
		}
		for s > x {
			ends = append(ends, i-1)
			s--
		}
	}
	res := []string{fmt.Sprintf("%d", len(starts))}
	for i := 0; i < len(starts); i++ {
		res = append(res, fmt.Sprintf("%d %d", starts[i], ends[i]))
	}
	return strings.Join(res, "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		expected, err := solveCase(line)
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
