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

func mutate(s string, l, r int) string {
	seg := s[l-1 : r]
	even := make([]byte, 0, (len(seg)+1)/2)
	odd := make([]byte, 0, len(seg)/2)
	for i := 1; i < len(seg); i += 2 {
		even = append(even, seg[i])
	}
	for i := 0; i < len(seg); i += 2 {
		odd = append(odd, seg[i])
	}
	copySeg := string(even) + string(odd)
	return s[:r] + copySeg + s[r:]
}

func expected(lines []string) string {
	s := strings.TrimSpace(lines[0])
	k, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
	n, _ := strconv.Atoi(strings.TrimSpace(lines[2]))
	cur := s
	idx := 3
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[idx])
		idx++
		l, _ := strconv.Atoi(parts[0])
		r, _ := strconv.Atoi(parts[1])
		cur = mutate(cur, l, r)
	}
	if k > len(cur) {
		k = len(cur)
	}
	return cur[:k]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if len(lines) == 0 {
				continue
			}
			idx++
			expect := expected(lines)
			input := strings.Join(lines, "\n") + "\n"
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
			if got != expect {
				fmt.Printf("test %d failed\nexpected: %s\n     got: %s\n", idx, expect, got)
				os.Exit(1)
			}
			lines = lines[:0]
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) > 0 {
		idx++
		expect := expected(lines)
		input := strings.Join(lines, "\n") + "\n"
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
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n     got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
