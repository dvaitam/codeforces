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
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "1887B.go")
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcasesB.txt: %v\n", err)
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
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			fmt.Fprintf(os.Stderr, "bad test file at case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < m; i++ {
			var cnt int
			fmt.Fscan(reader, &cnt)
			sb.WriteString(fmt.Sprintf("%d\n", cnt))
			for j := 0; j < cnt; j++ {
				var u, v int
				fmt.Fscan(reader, &u, &v)
				sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			}
		}
		var k int
		fmt.Fscan(reader, &k)
		sb.WriteString(fmt.Sprintf("%d\n", k))
		for i := 0; i < k; i++ {
			var x int
			fmt.Fscan(reader, &x)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(x))
		}
		sb.WriteByte('\n')
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
