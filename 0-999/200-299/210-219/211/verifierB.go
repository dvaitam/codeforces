package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}
	s := fields[0]
	m, _ := strconv.Atoi(fields[1])
	sets := fields[2:]
	qMasks := make([]uint32, m)
	for i := 0; i < m; i++ {
		mask := uint32(0)
		for _, c := range sets[i] {
			mask |= 1 << (c - 'a')
		}
		qMasks[i] = mask
	}
	counts := make(map[uint32]int)
	for _, m := range qMasks {
		counts[m] = 0
	}
	n := len(s)
	for a := 0; a < n; a++ {
		mask := uint32(0)
		for b := a; b < n; b++ {
			mask |= 1 << (s[b] - 'a')
			if _, ok := counts[mask]; ok {
				leftOk := a > 0 && (mask|(1<<(s[a-1]-'a'))) == mask
				rightOk := b+1 < n && (mask|(1<<(s[b+1]-'a'))) == mask
				if !leftOk && !rightOk {
					counts[mask]++
				}
			}
		}
	}
	var sb strings.Builder
	for i, mask := range qMasks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", counts[mask]))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
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
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
