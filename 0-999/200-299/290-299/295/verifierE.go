package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "295E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return oracle, nil
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcases:", err)
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
		if len(parts) < 2 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		p := 0
		n := atoi(parts[p])
		p++
		if len(parts) < 1+n {
			fmt.Printf("test %d incomplete\n", idx)
			os.Exit(1)
		}
		coords := parts[p : p+n]
		p += n
		if len(parts) <= p {
			fmt.Printf("test %d incomplete\n", idx)
			os.Exit(1)
		}
		m := atoi(parts[p])
		p++
		queries := parts[p:]
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		input.WriteString(strings.Join(coords, " ") + "\n")
		fmt.Fprintf(&input, "%d\n", m)
		// Build queries lines
		qp := 0
		for i := 0; i < m; i++ {
			t := atoi(queries[qp])
			qp++
			if t == 1 {
				a := queries[qp]
				b := queries[qp+1]
				qp += 2
				input.WriteString(fmt.Sprintf("1 %s %s\n", a, b))
			} else {
				a := queries[qp]
				b := queries[qp+1]
				qp += 2
				input.WriteString(fmt.Sprintf("2 %s %s\n", a, b))
			}
		}

		expectCmd := exec.Command(oracle)
		expectCmd.Stdin = strings.NewReader(input.String())
		var exp bytes.Buffer
		expectCmd.Stdout = &exp
		if err := expectCmd.Run(); err != nil {
			fmt.Printf("oracle runtime error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(exp.String())

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
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
