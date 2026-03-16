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
	cmd := exec.Command("go", "build", "-o", oracle, "679D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errb.String(), err
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesDRaw = `5 6 2 5 2 3 1 2 3 4 1 5 1 4
3 3 1 2 1 3 2 3
5 4 2 3 1 3 2 5 1 4
3 2 1 2 1 3
2 1 1 2
4 5 1 4 2 4 2 3 3 4 1 2
2 1 1 2
2 1 1 2
3 2 2 3 1 3
3 3 1 3 2 3 1 2
5 6 2 4 4 5 2 3 1 4 1 3 3 4
2 1 1 2
2 1 1 2
4 5 1 4 3 4 2 3 1 2 1 3
3 2 2 3 1 2
5 5 2 3 2 4 1 2 2 5 3 4
3 3 1 2 2 3 1 3
2 1 1 2
2 1 1 2
5 5 3 5 2 3 1 2 2 4 1 4
3 2 2 3 1 2
3 3 2 3 1 2 1 3
3 2 1 2 2 3
5 4 2 4 2 5 2 3 1 3
3 2 1 2 1 3
3 2 1 2 2 3
2 1 1 2
2 1 1 2
3 2 2 3 1 2
5 6 1 2 4 5 1 5 1 3 3 4 2 3
5 5 3 5 2 4 2 3 1 2 1 3
2 1 1 2
5 6 3 4 1 4 2 5 4 5 1 5 2 3
2 1 1 2
5 6 1 5 3 4 1 4 2 5 2 4 1 2
2 1 1 2
4 3 1 2 3 4 1 4
2 1 1 2
3 2 1 2 2 3
5 4 1 3 2 3 1 4 1 5
4 4 2 4 3 4 1 3 1 4
2 1 1 2
2 1 1 2
5 6 3 4 2 5 3 5 1 5 4 5 2 3
5 5 3 4 2 5 2 4 1 3 4 5
4 4 2 4 2 3 1 2 1 3
2 1 1 2
4 4 1 3 2 4 2 3 1 2
4 4 2 3 1 2 3 4 1 3
2 1 1 2
5 4 2 5 3 4 1 2 3 5
4 5 1 3 2 3 3 4 2 4 1 2
4 5 3 4 1 4 2 3 1 3 1 2
4 4 1 3 1 2 1 4 2 3
2 1 1 2
5 4 2 4 3 4 1 3 4 5
5 4 3 5 2 4 1 4 1 3
3 2 1 2 1 3
3 2 2 3 1 3
3 2 1 3 1 2
3 2 1 2 1 3
5 6 1 4 3 4 3 5 2 4 1 2 2 5
2 1 1 2
2 1 1 2
5 4 1 3 1 2 3 5 4 5
2 1 1 2
3 3 2 3 1 2 1 3
2 1 1 2
4 3 1 2 2 3 2 4
3 3 1 2 1 3 2 3
5 5 1 2 1 3 2 5 4 5 3 4
3 3 2 3 1 2 1 3
3 3 1 3 1 2 2 3
4 5 2 4 1 2 1 3 3 4 2 3
5 6 4 5 2 5 3 4 1 4 3 5 2 3
2 1 1 2
4 4 1 2 3 4 1 4 2 4
2 1 1 2
3 3 2 3 1 3 1 2
3 3 2 3 1 2 1 3
3 2 2 3 1 3
5 6 1 3 2 4 1 4 2 5 2 3 1 2
2 1 1 2
3 3 1 2 2 3 1 3
5 6 1 4 3 4 2 5 1 5 4 5 2 3
5 6 2 4 1 2 2 3 4 5 3 5 2 5
2 1 1 2
2 1 1 2
3 3 2 3 1 3 1 2
5 5 2 3 3 4 2 5 1 3 1 5
5 6 4 5 3 4 1 4 2 3 1 5 3 5
4 5 2 4 1 2 2 3 1 3 1 4
2 1 1 2
3 2 1 3 2 3
4 3 2 4 3 4 1 3
4 3 2 4 1 2 1 3
5 6 1 3 1 4 2 4 4 5 3 4 3 5
3 3 1 3 2 3 1 2
2 1 1 2
4 4 1 2 1 3 3 4 2 4`

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		n := atoi(fields[0])
		m := atoi(fields[1])
		if len(fields) != 2+2*m {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < m; i++ {
			a := fields[2+2*i]
			b := fields[2+2*i+1]
			fmt.Fprintf(&sb, "%s %s\n", a, b)
		}
		inputStr := sb.String()

		exp, errStr, err := runBinary(oracle, inputStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n%s", idx, err, errStr)
			os.Exit(1)
		}
		got, errStr2, err := runBinary(bin, inputStr)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errStr2)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
