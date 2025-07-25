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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func digitSumStr(s string) int {
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += int(s[i] - '0')
	}
	return sum
}

func bestNumber(x int64) int64 {
	s := strconv.FormatInt(x, 10)
	best := x
	bestSum := digitSumStr(s)
	for i := 0; i < len(s); i++ {
		if s[i] == '0' {
			continue
		}
		t := s[:i] + string(s[i]-1) + strings.Repeat("9", len(s)-i-1)
		num, _ := strconv.ParseInt(t, 10, 64)
		if num <= 0 {
			continue
		}
		sum := digitSumStr(t)
		if sum > bestSum || (sum == bestSum && num > best) {
			best = num
			bestSum = sum
		}
	}
	return best
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
		x, _ := strconv.ParseInt(line, 10, 64)
		expect := strconv.FormatInt(bestNumber(x), 10)
		input := fmt.Sprintf("%d\n", x)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expect, out)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
