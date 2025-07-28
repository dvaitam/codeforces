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

func solve(r, b, d int64) string {
	if r < b {
		r, b = b, r
	}
	if r <= b*(d+1) {
		return "YES"
	}
	return "NO"
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
		if len(parts) != 3 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		r, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		d, _ := strconv.ParseInt(parts[2], 10, 64)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d %d %d\n", r, b, d)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.ToUpper(strings.TrimSpace(out.String()))
		exp := solve(r, b, d)
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
