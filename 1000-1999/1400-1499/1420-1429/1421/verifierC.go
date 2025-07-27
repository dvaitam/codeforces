package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	s string
}

func buildCase(s string) testCase { return testCase{s: s} }

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 3 // length 3..22
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return buildCase(string(b))
}

func reverseBytes(b []byte) []byte {
	res := make([]byte, len(b))
	for i := range b {
		res[i] = b[len(b)-1-i]
	}
	return res
}

func applyOp(s []byte, opType string, idx int) ([]byte, error) {
	if idx < 2 || idx > len(s)-1 {
		return nil, fmt.Errorf("index out of range")
	}
	if opType == "L" {
		prefix := reverseBytes(s[1:idx])
		s = append(prefix, s...)
	} else if opType == "R" {
		suffix := reverseBytes(s[idx-1 : len(s)-1])
		s = append(s, suffix...)
	} else {
		return nil, fmt.Errorf("bad op type")
	}
	if len(s) > 1_000_000 {
		return nil, fmt.Errorf("length exceeds limit")
	}
	return s, nil
}

func isPalindrome(b []byte) bool {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		if b[i] != b[j] {
			return false
		}
	}
	return true
}

func runCase(bin string, tc testCase) error {
	input := tc.s + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	if k < 0 || k > 30 {
		return fmt.Errorf("k out of range")
	}
	type operation struct {
		typ string
		idx int
	}
	ops := make([]operation, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("incomplete op type")
		}
		typ := scanner.Text()
		if !scanner.Scan() {
			return fmt.Errorf("incomplete op index")
		}
		idxStr := scanner.Text()
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		ops[i] = operation{typ: typ, idx: idx}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	b := []byte(tc.s)
	for _, op := range ops {
		typ := op.typ
		idx := op.idx
		var err error
		b, err = applyOp(b, typ, idx)
		if err != nil {
			return err
		}
	}
	if !isPalindrome(b) {
		return fmt.Errorf("result is not palindrome")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, buildCase("abac"))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", idx+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
