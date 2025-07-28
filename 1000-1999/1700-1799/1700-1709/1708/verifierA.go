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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expectedOutputA(nums []int64) string {
	if len(nums) == 0 {
		return ""
	}
	g := nums[0]
	for i := 1; i < len(nums); i++ {
		g = gcd(g, nums[i])
	}
	if g == nums[0] {
		return "YES"
	}
	return "NO"
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesA.txt:", err)
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
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || n <= 0 || len(fields)-1 < n {
			fmt.Fprintf(os.Stderr, "bad test case at line %d\n", idx+1)
			os.Exit(1)
		}
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			val, err := strconv.ParseInt(fields[1+i], 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "bad number at line %d\n", idx+1)
				os.Exit(1)
			}
			nums[i] = val
		}
		expected := expectedOutputA(nums)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		got, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expected {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
