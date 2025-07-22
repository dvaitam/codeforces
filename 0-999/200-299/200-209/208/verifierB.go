package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func match(a, b string) bool {
	return a[0] == b[0] || a[1] == b[1]
}

func canSolve(piles []string) bool {
	type state struct {
		piles []string
	}
	encode := func(p []string) string {
		return strings.Join(p, ",")
	}
	seen := make(map[string]bool)
	queue := [][]string{append([]string(nil), piles...)}
	seen[encode(queue[0])] = true
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if len(cur) == 1 {
			return true
		}
		x := len(cur) - 1
		// try move to x-3
		if x >= 3 && match(cur[x], cur[x-3]) {
			nxt := append([]string(nil), cur...)
			nxt[x-3] = cur[x]
			nxt = append(nxt[:x], nxt[x+1:]...)
			key := encode(nxt)
			if !seen[key] {
				seen[key] = true
				queue = append(queue, nxt)
			}
		}
		// try move to x-1
		if x >= 1 && match(cur[x], cur[x-1]) {
			nxt := append([]string(nil), cur...)
			nxt[x-1] = cur[x]
			nxt = append(nxt[:x], nxt[x+1:]...)
			key := encode(nxt)
			if !seen[key] {
				seen[key] = true
				queue = append(queue, nxt)
			}
		}
	}
	return false
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
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n := 0
		fmt.Sscan(fields[0], &n)
		if len(fields) != n+1 {
			fmt.Printf("test %d: expected %d cards got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		cards := fields[1:]
		expect := "NO"
		if canSolve(cards) {
			expect = "YES"
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for i, c := range cards {
			sb.WriteString(c)
			if i+1 < len(cards) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
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
		got = strings.ToUpper(got)
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
