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

func solveB(n, k int, s []int64) string {
	if k == 1 {
		return "YES"
	}
	diff := make([]int64, k-1)
	for i := 1; i < k; i++ {
		diff[i-1] = s[i] - s[i-1]
	}
	for i := 1; i < len(diff); i++ {
		if diff[i] < diff[i-1] {
			return "NO"
		}
	}
	firstCount := int64(n - k + 1)
	if diff[0]*firstCount < s[0] {
		return "NO"
	}
	return "YES"
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
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+k {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		svals := make([]int64, k)
		for i := 0; i < k; i++ {
			v, _ := strconv.ParseInt(parts[2+i], 10, 64)
			svals[i] = v
		}
		expect := solveB(n, k, svals)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < k; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", svals[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(strings.ToUpper(out.String()))
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
