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
	cmd := exec.Command("go", "build", "-o", oracle, "679E.go")
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
		if len(fields) < 2+n {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		arr := fields[2 : 2+n]
		pos := 2 + n
		queries := make([]string, 0, q)
		for j := 0; j < q; j++ {
			if pos >= len(fields) {
				fmt.Printf("bad case %d\n", idx)
				os.Exit(1)
			}
			typ := fields[pos]
			pos++
			switch typ {
			case "1":
				if pos >= len(fields) {
					fmt.Printf("bad case %d\n", idx)
					os.Exit(1)
				}
				idxVal := fields[pos]
				pos++
				queries = append(queries, fmt.Sprintf("1 %s", idxVal))
			case "2":
				if pos+2 >= len(fields) {
					fmt.Printf("bad case %d\n", idx)
					os.Exit(1)
				}
				l := fields[pos]
				r := fields[pos+1]
				x := fields[pos+2]
				pos += 3
				queries = append(queries, fmt.Sprintf("2 %s %s %s", l, r, x))
			case "3":
				if pos+2 >= len(fields) {
					fmt.Printf("bad case %d\n", idx)
					os.Exit(1)
				}
				l := fields[pos]
				r := fields[pos+1]
				x := fields[pos+2]
				pos += 3
				queries = append(queries, fmt.Sprintf("3 %s %s %s", l, r, x))
			default:
				fmt.Printf("bad case %d\n", idx)
				os.Exit(1)
			}
		}
		if pos != len(fields) {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, q)
		sb.WriteString(strings.Join(arr, " "))
		sb.WriteByte('\n')
		for _, qs := range queries {
			sb.WriteString(qs)
			sb.WriteByte('\n')
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
