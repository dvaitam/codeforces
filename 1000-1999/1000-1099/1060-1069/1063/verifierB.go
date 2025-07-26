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
	cmd := exec.Command("go", "build", "-o", oracle, "1063B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scanner.Scan() {
			fmt.Printf("unexpected EOF on case %d\n", caseNum)
			os.Exit(1)
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) != 6 {
			fmt.Printf("bad header on case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		sr, _ := strconv.Atoi(parts[2])
		sc, _ := strconv.Atoi(parts[3])
		x, _ := strconv.Atoi(parts[4])
		y, _ := strconv.Atoi(parts[5])
		rows := make([]string, n)
		for i := 0; i < n; i++ {
			if !scanner.Scan() {
				fmt.Printf("unexpected EOF on case %d\n", caseNum)
				os.Exit(1)
			}
			rows[i] = scanner.Text()
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		sb.WriteString(fmt.Sprintf("%d %d\n", sr, sc))
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		for _, r := range rows {
			sb.WriteString(r)
			sb.WriteByte('\n')
		}
		input := sb.String()
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		out, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
