package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expectedAnswer(n int, s, t string) string {
	f1 := strings.IndexByte(s, '1')
	if f1 == -1 {
		f1 = n
	}
	f2 := strings.IndexByte(t, '1')
	if f2 == -1 {
		f2 = n
	}
	if f1 <= f2 {
		return "YES"
	}
	return "NO"
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
		var n int
		var s, t string
		fmt.Sscan(line, &n, &s, &t)
		expect := expectedAnswer(n, s, t)
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, s, t)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
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
		got = strings.ToUpper(got)
		expect = strings.ToUpper(expect)
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
