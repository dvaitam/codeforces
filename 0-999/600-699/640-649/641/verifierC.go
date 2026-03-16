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
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "641C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	const testcasesCRaw = `8 1 2
8 3 2 1 -5 2
4 2 1 -3 2
4 1 2
6 4 1 -1 2 2 1 5
2 3 1 -1 2 1 0
6 2 2 1 3
8 1 2
2 5 1 2 1 1 1 2 1 0 1 2
4 2 1 2 1 3
6 2 2 2
4 2 2 2
4 3 1 -3 2 2
6 1 1 3
2 3 1 -1 2 1 2
6 4 1 0 2 1 -5 2
6 2 2 1 4
8 6 2 1 -8 1 0 1 2 2 2
2 6 1 2 1 0 1 1 2 2 1 0
6 5 2 2 1 0 2 1 2
2 2 2 2
2 2 1 0 1 -1
4 4 1 1 2 1 4 1 3
6 3 2 2 2
2 2 2 1 -1
4 2 2 1 -2
6 5 1 -4 1 5 2 1 -3 1 3
4 6 1 4 2 1 3 1 -4 2 1 2
8 2 2 2
2 2 1 2 1 -2
2 6 1 1 2 2 2 2 2
6 6 1 6 1 3 1 5 1 -6 1 2 1 5
8 2 2 1 8
6 1 1 -1
6 3 2 2 1 -6
8 6 1 4 2 2 2 2 1 7
6 6 2 2 2 2 2 2
4 2 2 1 0
4 5 1 4 2 2 2 2
6 4 1 6 1 -2 2 1 -6
4 5 1 0 1 -1 1 -2 2 2
2 6 1 -2 1 0 1 1 2 2 2
2 5 2 1 -2 1 -2 2 2
4 4 1 1 1 -2 2 1 4
8 5 1 -8 1 5 1 -4 1 5 2
6 6 1 5 1 -1 1 -1 2 2 2
6 3 1 1 1 -5 2
8 5 2 2 1 5 1 -5 1 -3
8 3 1 5 2 2
2 6 2 1 -2 1 2 1 -1 1 1 2
8 1 2
4 1 2
4 1 1 2
2 1 2
4 5 2 2 2 1 3 1 4
4 6 1 -2 1 3 2 2 2 1 2
2 6 1 2 1 -2 2 1 -1 1 -2 2
2 6 2 1 0 1 -2 2 1 0 2
6 6 2 2 1 -5 2 2 2
8 5 2 2 1 -1 1 4 2
4 1 1 -2
2 4 1 -2 2 2 2
8 4 2 2 2 1 8
2 6 2 2 2 1 -1 1 1 1 0
8 5 1 -8 1 -7 2 1 3 1 2
6 4 2 2 1 6 1 4
8 6 2 1 -7 2 1 4 2 2
4 1 2
4 5 2 2 2 2 1 4
8 4 1 3 2 2 1 -2
6 2 2 2
4 5 2 1 1 1 -2 2 1 4
2 4 2 1 -1 2 1 -2
8 2 2 2
6 5 1 -1 2 2 1 -6 1 -3
4 5 2 1 0 2 2 2
6 6 2 1 -5 1 -1 1 6 2 2
6 5 2 2 2 1 4 1 5
2 6 1 2 1 2 2 1 2 1 -2 1 -2
6 1 2
4 5 1 0 2 1 -2 1 -2 1 -1
6 4 2 1 5 2 1 -5
4 6 2 1 -1 1 1 2 1 0 1 0
4 1 1 -1
2 4 1 -1 2 1 2 2
4 6 1 3 2 2 1 0 1 -3 2
8 3 2 2 1 1
8 1 1 4
2 5 1 2 1 1 1 1 2 1 -1
2 3 2 2 2
2 3 2 2 1 2
4 6 1 0 2 2 1 3 1 2 1 -3
8 1 2
4 1 1 -1
4 1 2
8 1 1 8
4 2 2 1 -1
2 5 1 2 1 2 2 2 2
8 6 1 3 2 1 -7 1 6 1 8 1 -8
4 2 2 2`

	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
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
		q := atoi(fields[1])
		ops := make([]string, 0, q)
		i := 2
		for op := 0; op < q; op++ {
			if i >= len(fields) {
				fmt.Printf("bad case %d\n", idx)
				os.Exit(1)
			}
			t := atoi(fields[i])
			i++
			if t == 1 {
				if i >= len(fields) {
					fmt.Printf("bad case %d\n", idx)
					os.Exit(1)
				}
				x := fields[i]
				i++
				ops = append(ops, fmt.Sprintf("1 %s", x))
			} else {
				ops = append(ops, "2")
			}
		}
		if i != len(fields) {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for _, opStr := range ops {
			input.WriteString(opStr)
			input.WriteByte('\n')
		}
		inputStr := input.String()

		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(inputStr)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		expected := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inputStr)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
