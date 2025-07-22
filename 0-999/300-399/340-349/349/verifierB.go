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

func expectedB(v int, cost []int) string {
	minCost := cost[0]
	for _, c := range cost[1:] {
		if c < minCost {
			minCost = c
		}
	}
	if v < minCost {
		return "-1"
	}
	length := v / minCost
	rem := v - length*minCost
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		for d := 8; d >= 0; d-- {
			extra := cost[d] - minCost
			if extra <= rem {
				res[i] = byte('1' + d)
				rem -= extra
				break
			}
		}
	}
	return string(res)
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
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 10 {
			fmt.Fprintf(os.Stderr, "test %d invalid line\n", idx)
			os.Exit(1)
		}
		v, _ := strconv.Atoi(parts[0])
		cost := make([]int, 9)
		for i := 0; i < 9; i++ {
			cost[i], _ = strconv.Atoi(parts[i+1])
		}
		expect := expectedB(v, cost)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", v))
		for i := 0; i < 9; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(parts[i+1])
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
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
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
