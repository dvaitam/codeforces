package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1830C.go")
	exe := filepath.Join(dir, "oracleC.bin")
	cmd := exec.Command("go", "build", "-o", exe, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func run(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesC.txt:", err)
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
		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			fmt.Fprintf(os.Stderr, "bad test case on line %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(tokens[0])
		k, _ := strconv.Atoi(tokens[1])
		need := 2 + 2*k
		if len(tokens) != need {
			fmt.Fprintf(os.Stderr, "line %d: expected %d numbers got %d\n", idx, need, len(tokens))
			os.Exit(1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, k)
		for i := 0; i < k; i++ {
			sb.WriteString(tokens[2+2*i])
			sb.WriteByte(' ')
			sb.WriteString(tokens[2+2*i+1])
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", idx, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\ngot:%s\ninput:%s", idx, exp, got, line)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
