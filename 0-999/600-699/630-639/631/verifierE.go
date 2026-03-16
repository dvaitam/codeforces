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
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "631E.go")
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
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesERaw = `100
5
-1 -7 2 5 -6
3
-8 -10 2
10
-1 -9 -3 6 7 1 -2 -5 -7 -2
5
-10 10 -2 -2 -4
4
-1 -1 10 1
3
9 0 2
10
-3 -5 -3 5 -2 -8 7 -1 -10 -1
6
6 -4 3 3 9 -1
8
4 -5 -3 -1 -2 -9 -8 -9
9
10 -2 6 7 10 5 0 -6 -4
3
3 -4 10
9
-2 -5 1 3 8 0 10 7 -4
7
-7 -9 -3 -2 8 9 -3
3
0 -5 -1
9
-10 -9 1 -8 -1 0 -10 0 -1
7
-6 10 3 9 -8 -1 9
5
4 -1 -6 -2 2
4
0 8 -10 1
2
4 -5
7
1 -1 8 -7 4 -4 3
5
-7 -9 -9 -9 -5
4
9 -9 7 5
5
0 -9 -7 6 -1
8
10 -4 5 -4 -3 4 3 5
2
-3 3
9
-3 10 3 -4 5 -4 -9 -9 -2
6
-3 6 -4 -3 3 -2
4
0 -9 0 8
3
8 2 10
2
5 2
3
3 -4 8
4
0 -1 5 10
7
3 6 -4 10 -2 0 2
9
-8 -2 10 -4 -9 2 9 -6 -2
2
-5 10
9
8 5 2 2 -4 -10 -4 -5 -10
6
-7 2 2 -3 7 -9
5
-5 9 0 7 5
10
4 -10 -8 -9 9 -7 5 7 -2 9
4
-9 1 -8 6
2
-1 1
3
-8 7 4
8
-4 -1 2 -3 5 2 -7 -8
3
9 1 6
8
3 4 -8 10 -4 10 -1 5
8
-7 7 -5 1 -5 -5 -6 0
9
0 -2 7 -10 -5 -10 10 -1 -7
10
-7 5 9 5 6 -8 6 -3 3 -1
7
-3 -5 10 -10 -9 9 0
10
4 8 -1 6 4 9 9 4 2 -6
6
9 1 0 -6 3 -8
4
9 -5 -1 1
5
8 1 9 -8 -8
8
10 -5 0 10 1 0 -5 -1
2
9 -10
10
-8 1 -7 -5 -5 8 5 8 -8 -7
4
10 5 -3 9
6
2 9 -3 5 -3 -1
7
-3 0 7 10 6 4 2
10
2 0 -1 4 3 8 -10 -2 -5 7
9
7 9 1 2 2 9 -10 -6 6
3
4 10 1
6
-8 -2 5 -3 10 5
3
-6 -3 -8
6
-6 -9 -5 2 8 9
10
7 -2 -10 6 -3 -5 -7 -4 -9 0
3
-7 -2 -9
6
9 10 -5 -6 10 3
4
-8 7 1 -10
10
-6 6 3 -6 -4 -1 5 6 -8 2
4
-5 -2 6 2
10
-1 2 0 -5 2 -9 3 -10 -2 -10
6
-6 -8 -5 -7 9 -10
5
-3 7 -10 5 7
4
4 2 0 -5
10
8 -4 -7 5 9 1 -1 10 3 9
2
10 -4
6
8 8 -1 5 10 0
3
-3 0 -7
2
9 0
10
5 1 -8 -5 -9 5 6 7 9 -3
2
-4 -8
7
9 -6 -1 -7 6 6 -10
3
-3 9 -2
2
-9 -10
9
-6 -4 1 -3 1 -1 2 10 9
8
-9 -5 3 5 -8 6 6 -1
10
-1 -4 3 0 -5 -4 -10 5 3 -2
7
3 10 3 0 9 -4 -2
6
6 -8 -10 2 -2 9
6
10 -6 9 0 -8 -5
3
-4 -2 -4
9
-7 10 0 -10 4 8 -4 -5 -8
5
6 -7 -7 0 -1
4
-8 -5 -10 -5
5
-2 -9 8 0 10
7
9 1 -5 -5 -4 -9 -9
4
-9 -2 5 -7
5
-6 2 5 9 7
9
7 0 1 3 2 -9 -9 -5 -6`

	reader := bufio.NewReader(strings.NewReader(testcasesERaw))
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
