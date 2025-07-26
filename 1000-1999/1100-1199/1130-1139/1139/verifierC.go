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
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, filepath.Join(dir, "1139C.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesC.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+3*(n-1) {
			fmt.Printf("bad testcase %d\n", idx)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d\n", n, k)
		for i := 0; i < n-1; i++ {
			u := fields[2+3*i]
			v := fields[2+3*i+1]
			w := fields[2+3*i+2]
			input += fmt.Sprintf("%s %s %s\n", u, v, w)
		}
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
