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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
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
		if len(parts) < 3 {
			fmt.Fprintf(os.Stderr, "case %d malformed\n", idx)
			os.Exit(1)
		}
		w, _ := strconv.Atoi(parts[0])
		h, _ := strconv.Atoi(parts[1])
		n, _ := strconv.Atoi(parts[2])
		if len(parts) != 3+2*n {
			fmt.Fprintf(os.Stderr, "case %d expected %d rectangle pairs got %d numbers\n", idx, n, len(parts)-3)
			os.Exit(1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", w, h)
		fmt.Fprintf(&sb, "%d\n", n)
		k := 3
		for i := 0; i < n; i++ {
			wi, _ := strconv.Atoi(parts[k])
			hi, _ := strconv.Atoi(parts[k+1])
			fmt.Fprintf(&sb, "%d %d\n", wi, hi)
			k += 2
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d lines of output got %d\n", idx, n, len(lines))
			os.Exit(1)
		}
		for li, l := range lines {
			fields := strings.Fields(l)
			if len(fields) != 4 {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected 4 numbers got %d\n", idx, li+1, len(fields))
				os.Exit(1)
			}
			for _, f := range fields {
				if _, err := strconv.ParseFloat(f, 64); err != nil {
					fmt.Fprintf(os.Stderr, "case %d line %d: invalid number %s\n", idx, li+1, f)
					os.Exit(1)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
