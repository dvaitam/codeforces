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
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "1500D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read testcases: %v\n", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	idx := 0
	for scan.Scan() {
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("bad case")
			os.Exit(1)
		}
		q, _ := strconv.Atoi(scan.Text())
		arr := make([]string, n*n)
		for i := 0; i < n*n; i++ {
			scan.Scan()
			arr[i] = scan.Text()
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, q)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				input.WriteString(arr[i*n+j])
			}
			input.WriteByte('\n')
		}
		idx++
		expect, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
