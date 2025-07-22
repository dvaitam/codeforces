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
	oracle := "oracleD"
	cmd := exec.Command("go", "build", "-o", oracle, "313D.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read testcases: %v\n", err)
		os.Exit(1)
	}
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	scan.Scan()
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < m; i++ {
			scan.Scan()
			l, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			c, _ := strconv.Atoi(scan.Text())
			input.WriteString(fmt.Sprintf("%d %d %d\n", l, r, c))
		}
		expected, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
