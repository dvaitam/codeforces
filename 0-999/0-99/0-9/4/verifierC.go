package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(names []string) []string {
	cnt := make(map[string]int)
	res := make([]string, len(names))
	for i, n := range names {
		if cnt[n] == 0 {
			res[i] = "OK"
			cnt[n] = 1
		} else {
			res[i] = fmt.Sprintf("%s%d", n, cnt[n])
			cnt[n]++
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
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
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(parts[0], &n)
		if len(parts)-1 < n {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		names := parts[1 : 1+n]
		exp := expected(names)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d\n", n)
		for _, name := range names {
			fmt.Fprintf(&buf, "%s\n", name)
		}
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(outLines) != len(exp) {
			fmt.Printf("Test %d failed: expected %d lines got %d\n", idx, len(exp), len(outLines))
			os.Exit(1)
		}
		for i := range exp {
			if strings.TrimSpace(outLines[i]) != exp[i] {
				fmt.Printf("Test %d failed line %d: expected %s got %s\n", idx, i+1, exp[i], outLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
