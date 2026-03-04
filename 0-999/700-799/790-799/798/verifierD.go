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

func run(bin, input string) (string, error) {
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

func validateOutput(n int, a, b []int64, out string) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %q", fields[0])
	}
	limit := n/2 + 1
	if k < 1 || k > limit {
		return fmt.Errorf("k out of range: got %d, allowed [1,%d]", k, limit)
	}
	if len(fields) != 1+k {
		return fmt.Errorf("expected %d indices, got %d", k, len(fields)-1)
	}

	seen := make(map[int]bool, k)
	var sumSelA, sumSelB int64
	var sumAllA, sumAllB int64
	for i := 1; i <= n; i++ {
		sumAllA += a[i]
		sumAllB += b[i]
	}

	for i := 0; i < k; i++ {
		idx, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return fmt.Errorf("invalid index token: %q", fields[1+i])
		}
		if idx < 1 || idx > n {
			return fmt.Errorf("index out of range: %d", idx)
		}
		if seen[idx] {
			return fmt.Errorf("duplicate index: %d", idx)
		}
		seen[idx] = true
		sumSelA += a[idx]
		sumSelB += b[idx]
	}

	if 2*sumSelA <= sumAllA {
		return fmt.Errorf("A inequality failed: 2*%d <= %d", sumSelA, sumAllA)
	}
	if 2*sumSelB <= sumAllB {
		return fmt.Errorf("B inequality failed: 2*%d <= %d", sumSelB, sumAllB)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesD.txt")
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
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+2*n {
			fmt.Printf("test %d wrong number of values\n", idx)
			os.Exit(1)
		}
		aVals := make([]int64, n+1)
		bVals := make([]int64, n+1)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				fmt.Printf("test %d invalid A value\n", idx)
				os.Exit(1)
			}
			aVals[i+1] = v
		}
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[1+n+i], 10, 64)
			if err != nil {
				fmt.Printf("test %d invalid B value\n", idx)
				os.Exit(1)
			}
			bVals[i+1] = v
		}
		input := fmt.Sprintf("%d\n%s\n%s\n", n, strings.Join(fields[1:1+n], " "), strings.Join(fields[1+n:], " "))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if err := validateOutput(n, aVals, bVals, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s got: %s\n", idx, err, input, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
