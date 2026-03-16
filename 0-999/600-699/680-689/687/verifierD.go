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
	oracle := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", oracle, "687D.go")
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

	const testcasesDRaw = `4 6 3 1 3 2 1 2 10 1 3 10 1 4 2 2 4 2 1 4 8 4 6 4 6 1 1
5 3 2 4 5 7 1 3 6 3 5 0 2 3 2 2
5 1 2 2 5 9 1 1 1 1
3 3 1 1 3 7 2 3 3 2 3 6 2 3
2 1 2 1 2 2 1 1 1 1
3 1 3 2 3 0 1 1 1 1 1 1
5 6 2 1 2 1 4 5 8 2 4 8 2 3 4 4 5 2 1 5 8 3 4 6 6
5 7 1 3 5 7 2 3 9 3 5 3 1 5 7 4 5 9 2 5 3 2 3 10 3 3
3 2 1 2 3 9 2 3 4 1 2
5 4 3 2 5 2 4 5 2 1 5 3 1 3 3 2 4 3 4 1 1
2 1 1 1 2 0 1 1
2 1 1 1 2 10 1 1
6 14 3 4 6 9 2 6 10 2 4 8 2 3 2 4 6 2 1 4 6 3 6 9 2 4 10 2 6 2 2 4 0 2 6 5 2 5 0 3 6 5 2 3 0 13 14 3 10 11 13
4 2 3 1 2 6 1 4 1 2 2 2 2 1 2
3 3 1 1 2 7 1 2 9 2 3 0 3 3
5 6 2 3 4 3 1 2 6 1 5 7 1 3 9 1 3 10 3 5 4 3 6 2 4
5 1 2 1 2 6 1 1 1 1
3 1 2 2 3 5 1 1 1 1
4 2 3 3 4 10 1 3 5 2 2 1 1 1 1
3 2 2 1 2 2 2 3 0 2 2 1 1
2 1 3 1 2 8 1 1 1 1 1 1
6 10 3 4 5 8 1 2 10 3 5 6 1 6 6 1 6 2 2 4 4 5 6 10 3 4 1 2 3 10 2 4 2 6 6 8 9 7 9
4 1 3 1 3 8 1 1 1 1 1 1
6 9 1 1 2 7 2 5 5 2 6 10 4 5 1 2 3 5 2 4 3 4 5 2 1 2 8 1 3 9 5 5
2 1 2 1 2 2 1 1 1 1
6 12 2 1 6 10 3 4 3 1 2 10 4 6 2 1 3 4 2 6 8 1 5 0 2 4 4 1 5 3 5 6 3 3 4 4 1 3 9 6 6 12 12
5 6 3 2 4 5 3 5 4 3 4 0 1 3 9 2 4 6 2 3 10 4 5 6 6 6 6
5 2 3 1 4 6 2 5 1 2 2 2 2 1 2
5 3 2 1 2 6 1 5 7 2 3 1 3 3 2 3
6 3 1 1 2 6 5 6 3 1 3 5 3 3
4 2 2 3 4 3 2 3 8 2 2 2 2
6 4 3 2 6 7 2 3 8 1 5 3 2 5 8 2 4 2 4 3 4
6 1 3 1 5 4 1 1 1 1 1 1
4 4 2 2 3 3 1 4 5 1 4 6 2 4 8 3 4 4 4
3 1 1 2 3 10 1 1
2 1 1 1 2 3 1 1
2 1 2 1 2 1 1 1 1 1
6 11 2 1 2 1 3 6 10 5 6 9 3 5 3 4 5 4 3 4 9 1 3 1 2 4 7 3 5 8 3 4 1 2 6 4 6 11 3 8
3 1 1 1 3 5 1 1
2 1 3 1 2 0 1 1 1 1 1 1
3 3 1 1 2 9 1 2 4 2 3 0 2 2
5 4 2 2 4 9 1 3 0 1 4 3 2 4 3 3 3 1 4
5 3 2 3 4 5 4 5 6 1 2 5 2 2 1 2
6 8 2 1 4 4 2 6 0 2 6 9 3 6 9 4 5 3 3 5 5 3 6 5 3 4 4 8 8 2 5
5 7 1 1 4 7 1 3 5 1 2 6 3 5 6 1 2 3 3 5 9 1 2 9 6 7
5 9 2 3 5 7 3 5 0 1 2 6 1 5 7 1 5 6 1 4 3 3 4 5 2 4 0 1 3 6 4 8 1 6
2 1 1 1 2 4 1 1
5 2 1 2 4 1 3 5 0 1 2
6 5 1 1 2 1 3 5 3 1 5 7 1 5 10 3 6 9 1 4
6 3 1 1 6 5 2 6 8 4 6 9 2 3
2 1 1 1 2 8 1 1
4 5 2 1 3 4 3 4 8 2 4 0 1 3 9 2 4 2 2 2 4 4
3 1 2 1 2 9 1 1 1 1
6 1 3 5 6 8 1 1 1 1 1 1
5 7 2 1 2 0 4 5 7 2 3 2 2 4 10 4 5 3 2 3 4 2 4 2 7 7 5 6
5 10 2 4 5 4 3 5 0 4 5 7 1 5 7 3 5 6 3 4 7 1 2 3 1 3 4 1 5 6 1 5 5 6 8 9 9
4 2 3 1 2 2 2 4 6 2 2 1 2 2 2
5 4 1 3 4 5 3 5 8 3 4 7 2 5 1 1 2
5 1 1 3 5 2 1 1
5 1 1 2 3 9 1 1
3 1 1 1 2 2 1 1
5 7 1 3 4 3 3 5 6 4 5 0 1 4 9 1 2 2 2 5 9 1 2 4 7 7
3 2 1 2 3 8 1 2 5 2 2
4 5 1 1 3 2 1 2 9 1 3 6 3 4 4 1 2 4 3 4
5 3 2 4 5 7 3 5 3 2 4 2 3 3 2 2
5 2 1 2 3 4 3 4 1 1 2
6 13 3 2 5 2 3 6 10 2 3 3 5 6 5 1 3 5 2 4 7 1 5 0 3 4 8 2 5 1 4 5 9 2 5 3 1 5 8 2 4 2 11 12 1 12 1 11
3 1 3 1 3 3 1 1 1 1 1 1
5 6 2 1 4 7 4 5 1 3 5 2 2 4 3 2 5 9 2 3 10 4 4 5 6
2 1 1 1 2 1 1 1
4 6 3 1 2 1 2 4 5 1 4 4 2 4 8 3 4 2 3 4 7 1 1 1 2 4 6
2 1 2 1 2 4 1 1 1 1
3 3 1 2 3 7 1 3 1 1 2 3 1 3
3 3 1 1 2 6 2 3 7 1 2 4 1 1
4 5 3 1 2 7 2 4 5 1 2 6 1 3 7 1 3 6 3 5 5 5 4 4
2 1 2 1 2 2 1 1 1 1
5 3 1 2 3 9 1 5 9 1 4 8 3 3
3 3 1 1 2 7 1 2 10 2 3 6 3 3
4 4 1 2 3 7 2 4 2 2 3 4 1 4 1 3 3
4 3 1 1 3 7 2 4 8 2 4 3 2 2
2 1 1 1 2 10 1 1
3 2 2 1 2 2 1 2 10 1 1 2 2
6 15 1 1 2 1 1 4 4 1 3 2 3 5 3 2 5 5 2 4 8 2 6 10 3 4 9 1 4 0 1 5 10 1 5 3 2 3 4 4 6 7 2 6 8 2 4 6 6 12
6 2 3 3 5 1 3 4 10 2 2 1 1 2 2
2 1 1 1 2 3 1 1
5 6 1 1 2 10 1 2 0 3 5 9 2 3 5 3 5 2 1 3 10 1 4
6 2 1 1 2 1 1 5 3 2 2
2 1 3 1 2 6 1 1 1 1 1 1
3 3 3 1 2 8 1 2 9 1 3 5 2 3 3 3 1 1
2 1 2 1 2 7 1 1 1 1
6 10 3 2 4 5 2 6 10 1 2 0 2 3 7 2 3 2 3 5 9 1 5 6 1 6 8 1 6 4 4 5 2 2 4 2 8 10 10
5 9 1 2 5 5 1 2 10 1 3 5 2 3 6 1 3 7 2 4 0 3 5 1 1 2 8 2 5 6 1 9
4 4 1 1 3 7 1 4 5 3 4 6 3 4 2 1 3
2 1 3 1 2 10 1 1 1 1 1 1
3 2 2 1 3 9 1 3 4 2 2 1 1
4 1 3 1 3 6 1 1 1 1 1 1
5 10 3 1 2 1 2 3 9 3 4 6 3 4 9 1 2 9 2 5 7 1 2 2 1 3 7 1 3 6 1 4 1 6 7 5 6 7 8
6 11 3 5 6 9 1 3 8 3 6 7 2 6 10 1 3 5 1 6 3 1 4 0 1 5 0 1 3 10 2 3 1 3 6 1 4 6 10 11 4 5
4 3 2 2 4 9 1 4 9 1 4 0 3 3 2 3
3 1 3 1 2 7 1 1 1 1 1 1`

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("bad test case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		q, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+3*m+2*q {
			fmt.Printf("bad test case %d\n", idx)
			os.Exit(1)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", n, m, q)
		off := 3
		for i := 0; i < m; i++ {
			u := fields[off+3*i]
			v := fields[off+3*i+1]
			w := fields[off+3*i+2]
			fmt.Fprintf(&input, "%s %s %s\n", u, v, w)
		}
		off += 3 * m
		for i := 0; i < q; i++ {
			l := fields[off+2*i]
			r := fields[off+2*i+1]
			fmt.Fprintf(&input, "%s %s\n", l, r)
		}
		exp, err := run(oracle, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
