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

func expectedBuses(n, m int, groups []int) int {
	buses := 0
	current := 0
	for _, a := range groups {
		if current+a > m {
			buses++
			current = a
		} else {
			current += a
		}
	}
	if current > 0 {
		buses++
	}
	return buses
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
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad test case on line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Fprintf(os.Stderr, "bad number of groups on line %d\n", idx)
			os.Exit(1)
		}
		groups := make([]int, n)
		for i := 0; i < n; i++ {
			groups[i], _ = strconv.Atoi(parts[2+i])
		}
		expect := expectedBuses(n, m, groups)
		input := fmt.Sprintf("%d %d\n", n, m)
		for i, a := range groups {
			if i > 0 {
				input += " "
			}
			input += strconv.Itoa(a)
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Fscan(&out, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, out.String())
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
