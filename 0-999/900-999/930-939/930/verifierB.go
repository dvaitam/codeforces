package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expectedB(s string) string {
	n := len(s)
	positions := make([][]int, 26)
	for i := 0; i < n; i++ {
		ch := s[i] - 'a'
		if ch >= 0 && ch < 26 {
			positions[ch] = append(positions[ch], i)
		}
	}
	total := 0
	var book [26]int
	for ch := 0; ch < 26; ch++ {
		if len(positions[ch]) == 0 {
			continue
		}
		maxRes := 0
		for l := 1; l < n; l++ {
			for i := 0; i < 26; i++ {
				book[i] = 0
			}
			for _, pos := range positions[ch] {
				np := pos + l
				if np >= n {
					np -= n
				}
				book[s[np]-'a']++
			}
			res := 0
			for i := 0; i < 26; i++ {
				if book[i] == 1 {
					res++
				}
			}
			if res > maxRes {
				maxRes = res
			}
		}
		total += maxRes
	}
	return fmt.Sprintf("%.12f", float64(total)/float64(n))
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
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		expect := expectedB(s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(s + "\n")
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
