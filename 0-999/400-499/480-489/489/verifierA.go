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

func parseInts(fields []string) ([]int, error) {
	vals := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		vals[i] = v
	}
	return vals, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		fields := strings.Fields(line)
		values, err := parseInts(fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test line %d: %v\n", idx, err)
			os.Exit(1)
		}
		if len(values) < 1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		n := values[0]
		if len(values) != n+1 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		arr := append([]int(nil), values[1:]...)
		in := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		tokens := strings.Fields(out)
		if len(tokens) == 0 {
			fmt.Fprintf(os.Stderr, "case %d: no output\n", idx)
			os.Exit(1)
		}
		k, err := strconv.Atoi(tokens[0])
		if err != nil || k < 0 || k > n {
			fmt.Fprintf(os.Stderr, "case %d: invalid k\n", idx)
			os.Exit(1)
		}
		if len(tokens) != 1+2*k {
			fmt.Fprintf(os.Stderr, "case %d: wrong number of tokens\n", idx)
			os.Exit(1)
		}
		p := 1
		for i := 0; i < k; i++ {
			a, err1 := strconv.Atoi(tokens[p])
			b, err2 := strconv.Atoi(tokens[p+1])
			if err1 != nil || err2 != nil || a < 0 || a >= n || b < 0 || b >= n {
				fmt.Fprintf(os.Stderr, "case %d: invalid swap index\n", idx)
				os.Exit(1)
			}
			arr[a], arr[b] = arr[b], arr[a]
			p += 2
		}
		for i := 1; i < n; i++ {
			if arr[i-1] > arr[i] {
				fmt.Fprintf(os.Stderr, "case %d: array not sorted\n", idx)
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
