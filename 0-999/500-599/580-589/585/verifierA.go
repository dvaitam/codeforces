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

func expectedA(n int, vals []int) string {
	v := make([]int, n)
	d := make([]int, n)
	p := make([]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		v[i] = vals[idx]
		d[i] = vals[idx+1]
		p[i] = vals[idx+2]
		idx += 3
	}
	removed := make([]bool, n)
	res := []int{}
	for i := 0; i < n; i++ {
		if removed[i] || p[i] < 0 {
			continue
		}
		res = append(res, i+1)
		cur := v[i]
		for j := i + 1; j < n && cur > 0; j++ {
			if !removed[j] {
				p[j] -= cur
			}
			cur--
		}
		for j := i + 1; j < n; j++ {
			if removed[j] {
				continue
			}
			if p[j] < 0 {
				removed[j] = true
				for k := j + 1; k < n; k++ {
					if !removed[k] {
						p[k] -= d[j]
					}
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for i, val := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(val))
	}
	if len(res) > 0 {
		sb.WriteByte('\n')
	}
	return sb.String()
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
		if len(parts) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		need := 1 + 3*n
		if len(parts) != need {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, need, len(parts))
			os.Exit(1)
		}
		vals := make([]int, 3*n)
		for i := 0; i < 3*n; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			vals[i] = v
		}
		expect := expectedA(n, vals)
		// prepare input
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", vals[3*i], vals[3*i+1], vals[3*i+2]))
		}
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
		got := strings.TrimSpace(out.String())
		expectTrim := strings.TrimSpace(expect)
		if got != expectTrim {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expectTrim, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
