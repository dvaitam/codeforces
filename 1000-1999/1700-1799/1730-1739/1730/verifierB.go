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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, x, t []int) string {
	minDiff := int64(1<<63 - 1)
	maxSum := int64(-1 << 63)
	for i := 0; i < n; i++ {
		diff := int64(x[i] - t[i])
		sum := int64(x[i] + t[i])
		if diff < minDiff {
			minDiff = diff
		}
		if sum > maxSum {
			maxSum = sum
		}
	}
	res := float64(minDiff+maxSum) * 0.5
	return fmt.Sprintf("%.6f", res)
}

func main() {
	arg := ""
	if len(os.Args) == 2 {
		arg = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		arg = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := arg
	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesB.txt: %v\n", err)
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
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != 1+2*n {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		xs := make([]int, n)
		ts := make([]int, n)
		for i := 0; i < n; i++ {
			xs[i], _ = strconv.Atoi(parts[1+i])
		}
		for i := 0; i < n; i++ {
			ts[i], _ = strconv.Atoi(parts[1+n+i])
		}
		expect := expected(n, xs, ts)
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, strings.Join(parts[1:1+n], " "), strings.Join(parts[1+n:], " "))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
