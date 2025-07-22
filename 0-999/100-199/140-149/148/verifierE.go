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
	cmd := exec.Command("go", "build", "-o", oracle, "148E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func readCases(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty test file")
	}
	t, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	var cases []string
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing header for case %d", i+1)
		}
		header := scanner.Text()
		parts := strings.Fields(header)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad header: %s", header)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		var sb strings.Builder
		sb.WriteString(header)
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("missing line in case %d", i+1)
			}
			sb.WriteString(scanner.Text())
			sb.WriteByte('\n')
		}
		cases = append(cases, sb.String())
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	cases, err := readCases("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}

	for i, input := range cases {
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
