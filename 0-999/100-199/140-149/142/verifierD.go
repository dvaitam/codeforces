package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(n, m, k int, rows []string) string {
	draw := false
	for i := 0; i < n; i++ {
		cnt := 0
		for _, c := range rows[i] {
			if c == 'G' || c == 'R' {
				cnt++
			}
		}
		if cnt <= 1 {
			draw = true
		}
	}
	if draw || k >= 2 {
		return "Draw"
	}
	xor := 0
	for i := 0; i < n; i++ {
		gpos, rpos := -1, -1
		for j, c := range rows[i] {
			if c == 'G' {
				gpos = j
			} else if c == 'R' {
				rpos = j
			}
		}
		if gpos >= 0 && rpos >= 0 {
			d := abs(gpos-rpos) - 1
			xor ^= d
		}
	}
	if xor != 0 {
		return "First"
	}
	return "Second"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not open testcasesD.txt:", err)
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
		var n, m, k int
		var board string
		if _, err := fmt.Sscan(line, &n, &m, &k, &board); err != nil {
			fmt.Printf("bad test case on line %d\n", idx)
			os.Exit(1)
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			start := i * m
			rows[i] = board[start : start+m]
		}
		expect := solveCase(n, m, k, rows)
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < n; i++ {
			b.WriteString(rows[i])
			b.WriteByte('\n')
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(b.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		var errBuf bytes.Buffer
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
