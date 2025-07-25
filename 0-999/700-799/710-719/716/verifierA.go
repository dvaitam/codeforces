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

func expectedA(times []int64, c int64) int {
	count := 0
	var last int64
	for i, t := range times {
		if i == 0 || t-last <= c {
			count++
		} else {
			count = 1
		}
		last = t
	}
	return count
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcasesA.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid format\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		c, _ := strconv.ParseInt(fields[1], 10, 64)
		if len(fields) != 2+n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d timestamps got %d\n", idx, n, len(fields)-2)
			os.Exit(1)
		}
		times := make([]int64, n)
		for i := 0; i < n; i++ {
			val, _ := strconv.ParseInt(fields[2+i], 10, 64)
			times[i] = val
		}
		input := fmt.Sprintf("%d %d\n", n, c)
		for i, t := range times {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", t)
		}
		input += "\n"
		expected := expectedA(times, c)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx, err, out)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output %q\n", idx, out)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
