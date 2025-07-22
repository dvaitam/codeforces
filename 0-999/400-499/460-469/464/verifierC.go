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

const MOD = 1000000007
const PHI = MOD - 1

func modPow(a, e int64) int64 {
	res := int64(1)
	base := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		e >>= 1
	}
	return res
}

func solveC(s string, ds []int, ts []string) int64 {
	val := make([]int64, 10)
	length := make([]int64, 10)
	for d := 0; d < 10; d++ {
		val[d] = int64(d)
		length[d] = 1
	}
	for i := len(ds) - 1; i >= 0; i-- {
		d := ds[i]
		t := ts[i]
		var newVal int64
		var newLen int64
		for _, ch := range t {
			cd := int(ch - '0')
			pow := modPow(10, length[cd])
			newVal = (newVal*pow + val[cd]) % MOD
			newLen = (newLen + length[cd]) % PHI
		}
		val[d] = newVal
		length[d] = newLen
	}
	var result int64
	for _, ch := range s {
		d := int(ch - '0')
		pow := modPow(10, length[d])
		result = (result*pow + val[d]) % MOD
	}
	return result
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

type caseC struct {
	s  string
	ds []int
	ts []string
}

func parseCases(path string) ([]caseC, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	tests := []caseC{}
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		s := line
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing q")
		}
		qLine := strings.TrimSpace(scanner.Text())
		q, _ := strconv.Atoi(qLine)
		ds := make([]int, q)
		ts := make([]string, q)
		for i := 0; i < q; i++ {
			if !scanner.Scan() {
				return nil, fmt.Errorf("incomplete case")
			}
			row := strings.TrimSpace(scanner.Text())
			parts := strings.Split(row, "->")
			d, _ := strconv.Atoi(parts[0])
			ds[i] = d
			ts[i] = parts[1]
		}
		tests = append(tests, caseC{s: s, ds: ds, ts: ts})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected := solveC(tc.s, tc.ds, tc.ts)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s\n%d\n", tc.s, len(tc.ds))
		for i := 0; i < len(tc.ds); i++ {
			fmt.Fprintf(&sb, "%d->%s\n", tc.ds[i], tc.ts[i])
		}
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: no output\n", idx+1)
			os.Exit(1)
		}
		got, err := strconv.Atoi(scan.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx+1)
			os.Exit(1)
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		if int64(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
