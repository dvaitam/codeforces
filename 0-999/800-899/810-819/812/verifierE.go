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
	cmd := exec.Command("go", "build", "-o", oracle, "812E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

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

	const testcasesRaw = `5 4 1 3 3 1 1 1 1 4
4 1 1 3 1 1 2 3
4 3 2 3 2 1 2 1
5 2 2 1 4 2 1 1 1 3
2 5 1 1
8 1 4 2 2 2 5 1 4 1 1 3 4 5 6 3
4 3 3 2 4 1 1 1
5 5 5 2 1 2 1 1 3 2
5 2 5 1 5 2 1 2 1 4
5 2 1 1 4 3 1 2 2 1
8 1 3 4 1 5 2 3 1 1 1 1 2 4 1 5
2 4 4 1
7 3 3 3 4 1 1 3 1 2 2 1 1 2
6 1 5 2 1 1 3 1 2 3 2 4
7 3 1 4 5 4 1 3 1 2 2 1 5 3
5 4 1 4 3 4 1 2 1 2
3 4 4 4 1 2
3 5 1 3 1 2
8 2 5 2 3 3 3 5 3 1 1 2 2 1 4 5
3 5 2 2 1 1
6 1 5 5 2 3 2 1 2 2 3 3
8 2 2 3 1 3 5 2 2 1 2 3 4 3 2 5
2 5 2 1
7 1 2 5 2 5 2 2 1 2 3 3 4 6
4 5 1 1 2 1 2 2
3 4 1 2 1 2
2 1 1 1
2 1 1 1
3 2 5 1 1 1
4 2 4 5 5 1 2 3
8 2 5 2 5 3 2 4 2 1 2 3 3 4 6 2
6 4 4 4 5 3 5 1 1 3 4 3
6 2 1 1 1 1 5 1 2 2 1 2
8 2 5 4 1 4 3 5 5 1 2 3 4 2 5 2
8 5 4 4 1 5 3 2 4 1 1 2 1 4 5 6
3 5 2 5 1 1
8 1 5 5 5 4 4 1 5 1 2 2 4 3 2 2
3 5 1 3 1 2
2 4 5 1
4 2 2 3 1 1 2 2
4 4 4 1 4 1 1 2
4 4 2 1 2 1 1 2
2 5 3 1
4 5 1 2 1 1 2 2
4 5 2 5 2 1 2 2
7 5 1 5 3 1 5 2 1 1 3 1 1 5
6 4 5 4 4 4 4 1 1 2 4 3
3 4 4 2 1 1
7 2 4 4 5 3 5 2 1 2 1 1 5 1
8 2 3 4 4 4 5 3 5 1 1 3 1 4 2 7
3 1 4 2 1 1
4 3 4 2 3 1 1 3
6 1 2 4 3 3 3 1 2 3 4 3
6 4 1 2 1 2 3 1 2 2 1 4
2 2 3 1
8 1 4 1 2 1 1 3 3 1 1 3 4 5 4 5
2 4 3 1
3 3 4 2 1 1
2 2 3 1
6 3 2 1 3 2 1 1 1 2 1 3
8 1 3 2 4 1 5 5 4 1 2 2 1 3 5 2
4 2 5 4 5 1 1 1
6 4 3 5 5 4 5 1 2 2 3 5
3 2 4 5 1 2
8 5 2 5 2 2 5 5 1 1 1 3 1 3 3 3
5 1 5 4 4 1 1 2 3 1
4 3 5 3 5 1 2 1
4 4 3 4 3 1 1 1
6 5 5 5 1 2 4 1 1 3 2 1
5 4 3 3 3 1 1 1 2 4
2 3 2 1
7 2 2 2 4 1 3 4 1 1 3 3 3 2
6 3 1 2 4 3 1 1 1 2 2 4
3 3 3 5 1 2
4 3 2 5 3 1 1 2
2 2 1 1
6 2 5 1 4 1 2 1 2 2 4 2
6 2 5 2 2 3 4 1 2 2 4 3
5 5 3 1 4 3 1 2 2 4
5 3 4 5 4 1 1 1 3 3
7 3 2 5 2 4 1 3 1 1 3 1 5 6
5 1 1 5 4 5 1 1 3 4
8 1 1 2 1 3 2 2 1 1 1 2 3 4 2 3
4 2 5 2 4 1 2 3
7 2 4 1 3 2 4 3 1 2 1 3 1 1
7 4 1 5 1 5 4 5 1 1 3 4 1 4
8 1 2 2 2 1 4 3 2 1 2 3 1 5 1 7
2 1 1 1
5 4 1 2 2 2 1 2 2 3
7 3 4 5 1 3 3 2 1 1 3 3 5 1
8 2 2 3 1 4 3 2 3 1 2 1 2 4 6 7
8 3 3 4 5 2 1 3 3 1 1 2 4 2 1 5
8 5 5 1 5 2 4 1 1 1 2 2 3 2 1 3
5 2 5 5 5 5 1 2 2 2
7 2 3 3 5 3 1 4 1 2 3 2 3 5
3 1 5 5 1 2
2 4 5 1
5 5 3 4 5 2 1 1 3 4
5 4 3 4 2 4 1 2 3 4
4 1 1 5 5 1 1 3`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("test %d invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		expectedLen := 1 + n + (n - 1)
		if len(parts) != expectedLen {
			fmt.Printf("test %d invalid number of values\n", idx)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(parts[0])
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(parts[1+i])
		}
		sb.WriteByte('\n')
		for i := 0; i < n-1; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(parts[1+n+i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
