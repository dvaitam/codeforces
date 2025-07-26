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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func positions(n, k, seg int) []int {
	rem := seg - (k+1)*2
	start := 1 + min(rem, k)
	step := 2*k + 1
	var pos []int
	for p := start; p <= n; p += step {
		pos = append(pos, p)
	}
	return pos
}

func expected(n, k int) string {
	if n < (k+1)*2 {
		pos := min(n, k+1)
		return fmt.Sprintf("1\n%d\n", pos)
	}
	step := 2*k + 1
	baseMin := (k + 1) * 2
	for i := baseMin; i <= baseMin+step; i++ {
		if (n-i)%step == 0 {
			pos := positions(n, k, i)
			var b strings.Builder
			fmt.Fprintf(&b, "%d\n", len(pos))
			for j, v := range pos {
				if j > 0 {
					b.WriteByte(' ')
				}
				fmt.Fprintf(&b, "%d", v)
			}
			b.WriteByte('\n')
			return b.String()
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", i)
	}
	b.WriteByte('\n')
	return b.String()
}

func runCase(bin, input, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expect) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expect, got)
	}
	return nil
}

func parseCase(line string) (int, int, error) {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid case")
	}
	n, _ := strconv.Atoi(parts[0])
	k, _ := strconv.Atoi(parts[1])
	return n, k, nil
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
		n, k, err := parseCase(line)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		expect := expected(n, k)
		if err := runCase(bin, input, expect); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
