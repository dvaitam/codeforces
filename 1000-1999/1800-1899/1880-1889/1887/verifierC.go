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
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1887C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcasesC.txt: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= t; caseIdx++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "bad test file at case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		var q int
		fmt.Fscan(reader, &q)
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			var l, r int
			var x int
			fmt.Fscan(reader, &l, &r, &x)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, x))
		}
		input := sb.String()
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", caseIdx, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\nexpected: %s\n got: %s\n", caseIdx, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
