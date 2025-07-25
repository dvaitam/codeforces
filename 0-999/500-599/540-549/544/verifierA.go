package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func possible(k int, q string) bool {
	seen := make([]bool, 26)
	cnt := 0
	for i := 0; i < len(q) && cnt < k; i++ {
		c := q[i] - 'a'
		if !seen[c] {
			seen[c] = true
			cnt++
		}
	}
	return cnt >= k
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 2 {
			fmt.Printf("test %d malformed line\n", idx)
			os.Exit(1)
		}
		k := 0
		for _, ch := range parts[0] {
			k = k*10 + int(ch-'0')
		}
		q := parts[1]
		input := fmt.Sprintf("%d\n%s\n", k, q)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Printf("test %d: empty output\n", idx)
			os.Exit(1)
		}
		ans := strings.TrimSpace(lines[0])
		if ans == "NO" {
			if possible(k, q) {
				fmt.Printf("test %d: answer should be YES\n", idx)
				os.Exit(1)
			}
			continue
		}
		if ans != "YES" {
			fmt.Printf("test %d: first line must be YES or NO\n", idx)
			os.Exit(1)
		}
		if !possible(k, q) {
			fmt.Printf("test %d: answer should be NO\n", idx)
			os.Exit(1)
		}
		if len(lines)-1 != k {
			fmt.Printf("test %d: expected %d lines, got %d\n", idx, k, len(lines)-1)
			os.Exit(1)
		}
		firsts := make(map[byte]bool)
		var concat strings.Builder
		for i := 1; i <= k; i++ {
			s := lines[i]
			if len(s) == 0 {
				fmt.Printf("test %d: empty string on line %d\n", idx, i)
				os.Exit(1)
			}
			if firsts[s[0]] {
				fmt.Printf("test %d: first characters not distinct\n", idx)
				os.Exit(1)
			}
			firsts[s[0]] = true
			concat.WriteString(s)
		}
		if concat.String() != q {
			fmt.Printf("test %d: concatenation mismatch\n", idx)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
