package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleA")
	src := filepath.Join("1000-1999", "1500-1599", "1570-1579", "1576", "1576A.go")
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tcPath := filepath.Join("1000-1999", "1500-1599", "1570-1579", "1576", "testcasesA.txt")
	f, err := os.Open(tcPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseIdx := 1; caseIdx <= T; caseIdx++ {
		var nodeCount, edgeCount, constrCount, flowCount int
		if _, err := fmt.Fscan(reader, &nodeCount, &edgeCount, &constrCount, &flowCount); err != nil {
			fmt.Fprintf(os.Stderr, "case %d header read error: %v\n", caseIdx, err)
			os.Exit(1)
		}
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d %d %d\n", nodeCount, edgeCount, constrCount, flowCount)
		for i := 0; i < edgeCount; i++ {
			var id, group, u, v, dist, cap int
			fmt.Fscan(reader, &id, &group, &u, &v, &dist, &cap)
			fmt.Fprintf(&b, "%d %d %d %d %d %d\n", id, group, u, v, dist, cap)
		}
		for i := 0; i < constrCount; i++ {
			var a, b1, c int
			fmt.Fscan(reader, &a, &b1, &c)
			fmt.Fprintf(&b, "%d %d %d\n", a, b1, c)
		}
		for i := 0; i < flowCount; i++ {
			var id, s, t, rate int
			fmt.Fscan(reader, &id, &s, &t, &rate)
			fmt.Fprintf(&b, "%d %d %d %d\n", id, s, t, rate)
		}
		input := b.String()
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", caseIdx, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", caseIdx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\n got: %s\n", caseIdx, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", T)
}
