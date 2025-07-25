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

func solveB(n, k int, x int64, arr []int64) int64 {
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] | arr[i]
	}
	suffix := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] | arr[i]
	}
	pow := int64(1)
	for i := 0; i < k; i++ {
		pow *= x
	}
	best := int64(0)
	for i := 0; i < n; i++ {
		val := prefix[i] | (arr[i] * pow) | suffix[i+1]
		if val > best {
			best = val
		}
	}
	return best
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	for {
		if !scanner.Scan() {
			break
		}
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		idx++
		parts := strings.Fields(line1)
		if len(parts) != 3 {
			fmt.Fprintf(os.Stderr, "case %d bad header\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		x64, _ := strconv.ParseInt(parts[2], 10, 64)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing array line\n", idx)
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scanner.Text())
		nums := strings.Fields(line2)
		if len(nums) != n {
			fmt.Fprintf(os.Stderr, "case %d expected %d numbers got %d\n", idx, n, len(nums))
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i, s := range nums {
			v, _ := strconv.ParseInt(s, 10, 64)
			arr[i] = v
		}
		expected := solveB(n, k, x64, arr)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", n, k, x64)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %s\n", idx, got)
			os.Exit(1)
		}
		if val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expected, val)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
