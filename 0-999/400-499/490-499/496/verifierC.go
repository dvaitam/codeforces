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

type testCase struct {
	n    int
	m    int
	rows []string
}

func parseTestcases(path string) ([]testCase, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cases []testCase
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields)-2 != n {
			return nil, fmt.Errorf("row count mismatch")
		}
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			if len(fields[2+i]) != m {
				return nil, fmt.Errorf("row length mismatch")
			}
			rows[i] = fields[2+i]
		}
		cases = append(cases, testCase{n: n, m: m, rows: rows})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveCase(tc testCase) int {
	n := tc.n
	m := tc.m
	s := tc.rows
	sorted := make([]bool, n-1)
	ans := 0
	for j := 0; j < m; j++ {
		del := false
		for i := 0; i < n-1; i++ {
			if !sorted[i] && s[i][j] > s[i+1][j] {
				del = true
				break
			}
		}
		if del {
			ans++
			continue
		}
		for i := 0; i < n-1; i++ {
			if !sorted[i] && s[i][j] < s[i+1][j] {
				sorted[i] = true
			}
		}
	}
	return ans
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i := 0; i < tc.n; i++ {
			sb.WriteString(tc.rows[i])
			sb.WriteByte('\n')
		}
		expected := strconv.Itoa(solveCase(tc))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
