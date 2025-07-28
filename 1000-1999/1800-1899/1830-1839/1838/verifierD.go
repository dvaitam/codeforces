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

func buildOracle() (string, error) {
	exe := "oracleD"
	cmd := exec.Command("go", "build", "-o", exe, "1838D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, string(out))
	}
	return exe, nil
}

func runProg(exe, input string) (string, error) {
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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		s := scan.Text()
		queries := make([]int, q)
		for i := 0; i < q; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
				os.Exit(1)
			}
			queries[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n%s\n", n, q, s)
		for _, v := range queries {
			fmt.Fprintf(&sb, "%d\n", v)
		}
		input := sb.String()
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
