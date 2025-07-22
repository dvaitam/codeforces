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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// compute expected output using reference solution 172E.go
func expected(doc string, queries []string) ([]int, error) {
	var input strings.Builder
	input.WriteString(doc)
	input.WriteByte('\n')
	input.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		input.WriteString(q)
		input.WriteByte('\n')
	}
	out, err := runProg("172E.go", input.String())
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(strings.NewReader(out))
	res := []int{}
	for scan.Scan() {
		v, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		res = append(res, v)
	}
	if len(res) != len(queries) {
		return nil, fmt.Errorf("reference output length mismatch")
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
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
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		m, err := strconv.Atoi(parts[1])
		if err != nil || len(parts) != m+2 {
			fmt.Fprintf(os.Stderr, "test %d bad m\n", idx)
			os.Exit(1)
		}
		doc := parts[0]
		queries := parts[2:]
		want, err := expected(doc, queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d reference error: %v\n", idx, err)
			os.Exit(1)
		}
		var input strings.Builder
		input.WriteString(doc)
		input.WriteByte('\n')
		input.WriteString(fmt.Sprintf("%d\n", m))
		for _, q := range queries {
			input.WriteString(q)
			input.WriteByte('\n')
		}
		gotStr, err := runProg(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(gotStr))
		outVals := []int{}
		for scan.Scan() {
			v, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
			outVals = append(outVals, v)
		}
		if len(outVals) != m {
			fmt.Fprintf(os.Stderr, "case %d wrong output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			if outVals[i] != want[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at query %d: expected %d got %d\n", idx, i+1, want[i], outVals[i])
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
