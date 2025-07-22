package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(digits []int) string {
	freq := make([]int, 10)
	for _, d := range digits {
		if d >= 0 && d <= 9 {
			freq[d]++
		}
	}
	if freq[0] == 0 {
		return "-1"
	}
	sum := 0
	for d, c := range freq {
		sum += d * c
	}
	mod := sum % 3
	mod1 := []int{}
	mod2 := []int{}
	for d := 1; d <= 9; d++ {
		for i := 0; i < freq[d]; i++ {
			if d%3 == 1 {
				mod1 = append(mod1, d)
			} else if d%3 == 2 {
				mod2 = append(mod2, d)
			}
		}
	}
	remove := func(list []int, cnt int) bool {
		if len(list) < cnt {
			return false
		}
		for i := 0; i < cnt; i++ {
			freq[list[i]]--
		}
		return true
	}
	switch mod {
	case 1:
		if !remove(mod1, 1) {
			if !remove(mod2, 2) {
				return "-1"
			}
		}
	case 2:
		if !remove(mod2, 1) {
			if !remove(mod1, 2) {
				return "-1"
			}
		}
	}
	if freq[0] == 0 {
		return "-1"
	}
	nonZero := 0
	for d := 1; d <= 9; d++ {
		nonZero += freq[d]
	}
	if nonZero == 0 {
		return "0"
	}
	var sb strings.Builder
	for d := 9; d >= 0; d-- {
		for i := 0; i < freq[d]; i++ {
			sb.WriteByte(byte('0' + d))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		if len(parts) < 1 {
			fmt.Fprintf(os.Stderr, "case %d: empty line\n", idx)
			os.Exit(1)
		}
		n := 0
		fmt.Sscan(parts[0], &n)
		if len(parts)-1 != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d digits got %d\n", idx, n, len(parts)-1)
			os.Exit(1)
		}
		digits := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(parts[i+1], &digits[i])
		}
		expect := solve(digits)
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(parts[1:], " "))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", idx, expect, outStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
