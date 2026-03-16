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
	oracle := filepath.Join(dir, "oracleA")
	cmd := exec.Command("go", "build", "-o", oracle, "590A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runProg(prog, input string) (string, error) {
	cmd := exec.Command(prog)
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

const testcasesARaw = `100
9 1 0 1 1 1 1 1 1 0
5 1 0 0 1 0
7 0 0 1 1 0 1 1
8 0 1 1 1 0 0 0 1
3 1 1 0
8 0 0 0 0 0 1 0 0
8 1 0 1 1 0 1 0 1
10 0 1 1 0 1 0 0 0 0 1
10 0 0 0 0 0 0 1 1 0 0
9 1 1 1 1 0 1 0 1 1
6 0 0 1 0 0 1
5 1 1 0 0 0
6 0 0 0 0 0 0
9 0 1 0 0 0 0 0 0 1
6 0 0 1 0 1 0
6 0 1 1 1 0 0
10 0 0 1 0 1 1 1 0 0 0
5 0 1 1 0 1
5 0 1 1 1 1
9 1 0 0 1 0 1 0 1 0
6 1 1 1 1 0 1
9 1 0 0 0 1 0 0 0 1
9 1 0 1 1 0 0 1 0 1
5 1 1 0 0 1
8 1 1 0 1 0 0 0 0
9 1 1 0 0 0 0 0 0 0
6 1 1 0 0 1 1
4 0 1 1 0
7 0 1 0 0 1 0 0
3 0 1 1
8 0 1 1 1 1 0 0 1
7 0 0 0 1 1 1 1
4 1 0 0 1
5 0 1 1 1 0
7 0 1 0 0 1 0 0
7 1 1 1 1 0 0 1
10 1 1 0 1 0 1 1 0 1 1
5 0 1 0 0 0
6 0 0 1 0 0 1
7 1 1 0 1 0 1 0
7 0 1 0 1 0 0 0
10 0 0 1 1 1 0 0 0 0 0
6 1 1 1 0 0 1
5 1 1 1 1 1
10 0 0 0 1 1 1 0 0 1 0
9 1 1 0 0 0 0 0 0 0
7 0 1 1 0 0 1 1
9 0 0 0 1 0 1 1 1 1
6 0 0 0 1 1 1
4 1 1 1 1
9 1 0 0 0 0 1 1 0 0
10 1 0 1 0 0 1 1 0 1 1
10 0 0 1 0 0 0 0 1 0 1
6 1 1 0 0 1 1
4 1 0 1 1
4 0 0 1 1
8 0 1 1 0 0 0 1 1
3 0 0 0
9 1 0 0 1 1 0 0 1 0
4 1 0 0 1
9 1 0 1 1 1 0 1 0 0
8 0 1 1 0 0 0 1 0
5 0 0 0 0 0
7 1 0 0 0 0 1 0
10 1 1 0 0 0 1 1 1 1 1
8 0 0 0 1 0 1 0 0
6 1 0 0 0 1 1
9 0 0 0 1 0 0 0 1 1
9 0 1 1 1 0 0 1 1 0
9 1 0 1 1 1 0 1 1 1
4 0 0 1 0
5 1 1 0 1 1
5 0 1 1 0 1
10 0 1 0 1 0 1 1 0 0 1
4 1 0 1 1
4 1 0 1 1
3 0 1 0
3 1 1 1
10 0 1 1 1 0 0 0 0 1 1
8 0 0 0 0 1 1 1 0
6 0 1 1 1 1 1
8 0 0 0 0 1 1 1 0
7 1 1 1 0 1 1 1
10 0 1 1 0 1 0 0 0 1 1
9 0 0 0 1 0 1 0 0 0
7 1 0 0 1 0 0 1
10 1 1 0 1 1 1 0 1 0 1
5 0 0 1 1 1
9 1 1 0 0 1 0 0 1 0
6 0 0 1 0 0 0
7 0 1 1 1 1 0 1
10 1 1 1 0 1 1 0 0 0 0
8 0 1 0 0 0 1 0 1
6 1 0 1 0 0 1
8 1 1 1 0 0 0 1 0
6 1 0 0 1 1 1
9 1 1 1 0 0 1 0 0 1
10 1 0 1 1 1 1 0 0 0 0
4 0 1 0 1
4 1 1 0 0
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	reader := bufio.NewReader(strings.NewReader(testcasesARaw))
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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
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
