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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(m, s int) (string, string) {
	if (s == 0 && m > 1) || s > 9*m {
		return "-1", "-1"
	}
	if s == 0 && m == 1 {
		return "0", "0"
	}
	rem := s
	max := make([]byte, m)
	for i := 0; i < m; i++ {
		d := 9
		if rem < 9 {
			d = rem
		}
		max[i] = byte('0' + d)
		rem -= d
	}
	rem = s
	min := make([]byte, m)
	for i := 0; i < m; i++ {
		low := 0
		if i == 0 {
			low = 1
		}
		for d := low; d <= 9; d++ {
			if rem-d < 0 {
				break
			}
			if rem-d <= 9*(m-i-1) {
				min[i] = byte('0' + d)
				rem -= d
				break
			}
		}
	}
	return string(min), string(max)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		mVal, _ := strconv.Atoi(fields[0])
		sVal, _ := strconv.Atoi(fields[1])
		in := line + "\n"
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		wantMin, wantMax := solve(mVal, sVal)
		tokens := strings.Fields(out)
		if len(tokens) != 2 {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		if tokens[0] != wantMin || tokens[1] != wantMax {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s %s got %s %s\n", idx, wantMin, wantMax, tokens[0], tokens[1])
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
