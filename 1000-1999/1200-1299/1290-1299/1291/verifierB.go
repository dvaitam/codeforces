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

func expectedB(n int, arr []int64) string {
	prefix := make([]bool, n)
	ok := true
	for i := 0; i < n; i++ {
		if ok && arr[i] >= int64(i) {
			prefix[i] = true
		} else {
			ok = false
		}
	}
	suffix := make([]bool, n)
	ok = true
	for i := n - 1; i >= 0; i-- {
		if ok && arr[i] >= int64(n-1-i) {
			suffix[i] = true
		} else {
			ok = false
		}
	}
	for i := 0; i < n; i++ {
		if prefix[i] && suffix[i] {
			return "Yes"
		}
	}
	return "No"
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
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n+1 {
			fmt.Printf("test %d expected %d numbers got %d\n", idx, n, len(parts)-1)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+i], 10, 64)
			arr[i] = v
		}
		expect := expectedB(n, arr)
		// build input string
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			input.WriteString(fmt.Sprint(arr[i]))
			if i+1 < n {
				input.WriteByte(' ')
			}
		}
		input.WriteByte('\n')
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(strings.ToLower(got))
		expectLower := strings.ToLower(expect)
		if got != expectLower {
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
