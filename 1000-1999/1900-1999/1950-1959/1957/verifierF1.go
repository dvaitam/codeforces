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
	src := filepath.Join(dir, "1957F1.go")
	bin := filepath.Join(os.TempDir(), "oracle1957F1.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func lineToInput(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return "", fmt.Errorf("empty line")
	}
	idx := 0
	n, err := strconv.Atoi(fields[idx])
	if err != nil {
		return "", err
	}
	idx++
	if len(fields) < idx+n {
		return "", fmt.Errorf("not enough values")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fields[idx+i])
	}
	sb.WriteByte('\n')
	idx += n
	if len(fields) < idx+2*(n-1) {
		return "", fmt.Errorf("not enough edges")
	}
	for i := 0; i < n-1; i++ {
		sb.WriteString(fmt.Sprintf("%s %s\n", fields[idx], fields[idx+1]))
		idx += 2
	}
	if len(fields) <= idx {
		return "", fmt.Errorf("missing q")
	}
	q, err := strconv.Atoi(fields[idx])
	if err != nil {
		return "", err
	}
	idx++
	sb.WriteString(fmt.Sprintf("%d\n", q))
	if len(fields) != idx+5*q {
		return "", fmt.Errorf("expected %d query numbers got %d", 5*q, len(fields)-idx)
	}
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%s %s %s %s %s\n", fields[idx], fields[idx+1], fields[idx+2], fields[idx+3], fields[idx+4]))
		idx += 5
	}
	return sb.String(), nil
}

func loadTests(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tests []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tests = append(tests, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests, err := loadTests("testcasesF1.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, line := range tests {
		input, err := lineToInput(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
