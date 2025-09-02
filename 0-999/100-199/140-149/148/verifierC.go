package main

import (
	"bufio"
	"bytes"
	"errors"
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
	cmd := exec.Command("go", "build", "-o", oracle, "148C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesC.txt")
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
		input := line + "\n"
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		oracleOut := strings.TrimSpace(outO.String())

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())

		// Parse test case parameters n a b
		n, a, b, perr := parseParams(line)
		if perr != nil {
			fmt.Printf("test %d: invalid testcase line: %q err: %v\n", idx, line, perr)
			os.Exit(1)
		}

		// Determine feasibility via oracle
		solvable := strings.TrimSpace(oracleOut) != "-1"

		// Verify candidate output
		if strings.TrimSpace(got) == "-1" {
			if solvable {
				fmt.Printf("test %d failed\nexpected any valid sequence for n=%d a=%d b=%d\n got: -1 (but a solution exists)\n", idx, n, a, b)
				os.Exit(1)
			}
			// both report impossible -> accept
			continue
		}

		seq, serr := parseSeq(got)
		if serr != nil {
			fmt.Printf("test %d failed: cannot parse output as sequence: %v\n got: %s\n", idx, serr, got)
			os.Exit(1)
		}
		if !solvable {
			fmt.Printf("test %d failed\nexpected: -1 (no solution)\n got: %s\n", idx, strings.Join(intsToStrs(seq), " "))
			os.Exit(1)
		}
		if verr := verifySeq(seq, n, a, b); verr != nil {
			fmt.Printf("test %d failed: %v\n got: %s\n", idx, verr, strings.Join(intsToStrs(seq), " "))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func parseParams(line string) (int, int, int, error) {
	fs := strings.Fields(line)
	if len(fs) < 3 {
		return 0, 0, 0, errors.New("expected three integers n a b")
	}
	n, err := strconv.Atoi(fs[0])
	if err != nil {
		return 0, 0, 0, err
	}
	a, err := strconv.Atoi(fs[1])
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.Atoi(fs[2])
	if err != nil {
		return 0, 0, 0, err
	}
	return n, a, b, nil
}

func parseSeq(out string) ([]int, error) {
	fs := strings.Fields(out)
	seq := make([]int, 0, len(fs))
	for _, tok := range fs {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		seq = append(seq, v)
	}
	return seq, nil
}

func verifySeq(seq []int, n, a, b int) error {
	if len(seq) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(seq))
	}
	// Check bounds and count oh/wow according to rules
	oh, wow := 0, 0
	sum := 0
	mx := 0
	for i, v := range seq {
		if v < 1 || v > 50000 {
			return fmt.Errorf("value out of bounds at position %d: %d (must be 1..50000)", i+1, v)
		}
		if i == 0 {
			sum = v
			mx = v
			continue
		}
		if v > sum {
			wow++
		} else if v > mx {
			oh++
		}
		sum += v
		if v > mx {
			mx = v
		}
	}
	if oh != a || wow != b {
		return fmt.Errorf("counts mismatch: expected oh=%d wow=%d, got oh=%d wow=%d", a, b, oh, wow)
	}
	return nil
}

func intsToStrs(xs []int) []string {
	ss := make([]string, len(xs))
	for i, v := range xs {
		ss[i] = strconv.Itoa(v)
	}
	return ss
}
