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

func expectedD(n, w int, a, b []int64) int {
	if w == 1 {
		return n
	}
	m := w - 1
	pattern := make([]int64, m)
	for i := 0; i < m; i++ {
		pattern[i] = b[i+1] - b[i]
	}
	text := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		text[i] = a[i+1] - a[i]
	}
	pi := make([]int, m)
	for i, j := 1, 0; i < m; i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	count := 0
	for i, j := 0, 0; i < n-1; i++ {
		for j > 0 && text[i] != pattern[j] {
			j = pi[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == m {
			count++
			j = pi[j-1]
		}
	}
	return count
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		wVal, _ := strconv.Atoi(parts[1])
		need := 2 + nVal + wVal
		if len(parts) != need {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		a := make([]int64, nVal)
		b := make([]int64, wVal)
		for i := 0; i < nVal; i++ {
			v, _ := strconv.ParseInt(parts[2+i], 10, 64)
			a[i] = v
		}
		for i := 0; i < wVal; i++ {
			v, _ := strconv.ParseInt(parts[2+nVal+i], 10, 64)
			b[i] = v
		}
		expect := expectedD(nVal, wVal, a, b)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", nVal, wVal))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: invalid output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
