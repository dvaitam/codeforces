package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isPalindrome(b []byte) bool {
	i, j := 0, len(b)-1
	for i < j {
		if b[i] != b[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func solveH(ops []string) []int {
	queue := make([]byte, 0)
	res := make([]int, 0, len(ops))
	for i := 0; i < len(ops); i++ {
		parts := strings.Fields(ops[i])
		if parts[0] == "push" {
			queue = append(queue, parts[1][0])
		} else if parts[0] == "pop" {
			if len(queue) > 0 {
				queue = queue[1:]
			}
		}
		seen := make(map[string]struct{})
		for i := 0; i < len(queue); i++ {
			for j := i; j < len(queue); j++ {
				if isPalindrome(queue[i : j+1]) {
					seen[string(queue[i:j+1])] = struct{}{}
				}
			}
		}
		res = append(res, len(seen))
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesH.txt")
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
		q := 0
		fmt.Sscan(parts[0], &q)
		ops := parts[1:]
		if len(ops) != q && len(ops) != q*2 {
			// operations may have spaces, but stored as tokens
		}
		// rebuild ops to phrases
		rebuilt := make([]string, 0, q)
		i := 1
		for len(rebuilt) < q {
			if parts[i] == "push" {
				rebuilt = append(rebuilt, parts[i]+" "+parts[i+1])
				i += 2
			} else { // pop
				rebuilt = append(rebuilt, parts[i])
				i++
			}
		}
		expect := solveH(rebuilt)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", q))
		for _, op := range rebuilt {
			input.WriteString(op)
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
		outParts := strings.Fields(strings.TrimSpace(out.String()))
		if len(outParts) != q {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, q, len(outParts))
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			got := outParts[i]
			if got != fmt.Sprintf("%d", expect[i]) {
				fmt.Printf("test %d failed at step %d expected %d got %s\n", idx, i+1, expect[i], got)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
