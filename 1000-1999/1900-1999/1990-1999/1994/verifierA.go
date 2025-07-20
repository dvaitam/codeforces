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

func expectedOutput(n, m int, a []int) string {
	if n == 1 && m == 1 {
		return "-1"
	}
	total := n * m
	out := make([]string, 0, n)
	idx := 0
	for i := 0; i < n; i++ {
		row := make([]string, m)
		for j := 0; j < m; j++ {
			x := (a[idx] % total) + 1
			row[j] = strconv.Itoa(x)
			idx++
		}
		out = append(out, strings.Join(row, " "))
	}
	return strings.Join(out, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
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
		m, _ := strconv.Atoi(parts[1])
		need := n*m + 2
		if len(parts) != need {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, need, len(parts))
			os.Exit(1)
		}
		a := make([]int, n*m)
		for i := 0; i < n*m; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			a[i] = v
		}
		expect := expectedOutput(n, m, a)
		// prepare input
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		idx2 := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				input.WriteString(fmt.Sprintf("%d", a[idx2]))
				idx2++
				if j+1 < m {
					input.WriteByte(' ')
				}
			}
			input.WriteByte('\n')
		}
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
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
