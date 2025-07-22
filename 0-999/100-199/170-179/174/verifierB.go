package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(s string) (string, error) {
	n := len(s)
	dots := make([]int, 0)
	for i := 0; i < n; i++ {
		if s[i] == '.' {
			dots = append(dots, i)
		}
	}
	m := len(dots)
	if m == 0 {
		return "NO", nil
	}
	if dots[0] < 1 || dots[0] > 8 {
		return "NO", nil
	}
	segLen := make([]int, m)
	for k := 0; k < m-1; k++ {
		delta := dots[k+1] - dots[k]
		low := delta - 9
		if low < 1 {
			low = 1
		}
		high := delta - 2
		if high > 3 {
			high = 3
		}
		if low > high {
			return "NO", nil
		}
		found := false
		for l := low; l <= high; l++ {
			ok := true
			for j := 1; j <= l; j++ {
				if dots[k]+j >= n || s[dots[k]+j] == '.' {
					ok = false
					break
				}
			}
			if ok {
				segLen[k] = l
				found = true
				break
			}
		}
		if !found {
			return "NO", nil
		}
	}
	lastL := n - 1 - dots[m-1]
	if lastL < 1 || lastL > 3 {
		return "NO", nil
	}
	for j := 1; j <= lastL; j++ {
		if dots[m-1]+j >= n || s[dots[m-1]+j] == '.' {
			return "NO", nil
		}
	}
	segLen[m-1] = lastL
	res := []string{"YES"}
	start := 0
	for k := 0; k < m; k++ {
		end := dots[k] + segLen[k]
		res = append(res, s[start:end+1])
		start = end + 1
	}
	return strings.Join(res, "\n"), nil
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
		expected, err := solveCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
