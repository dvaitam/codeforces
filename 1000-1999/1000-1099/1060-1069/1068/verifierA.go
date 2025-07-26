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

func solveCase(n, m, k, l int64) string {
	if n < m || l > n || k >= n {
		return "-1"
	}
	need := k + l
	packs := need / m
	if need%m != 0 {
		packs++
	}
	if packs*m <= n {
		return fmt.Sprintf("%d", packs)
	}
	return "-1"
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
		parts := strings.Fields(line)
		if len(parts) != 4 {
			fmt.Printf("invalid test line %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.ParseInt(parts[0], 10, 64)
		mVal, _ := strconv.ParseInt(parts[1], 10, 64)
		kVal, _ := strconv.ParseInt(parts[2], 10, 64)
		lVal, _ := strconv.ParseInt(parts[3], 10, 64)
		expected := solveCase(nVal, mVal, kVal, lVal)
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
