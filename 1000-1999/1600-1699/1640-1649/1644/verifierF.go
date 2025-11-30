package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesF.txt.
const encodedTestcases = `
H4sIAA8IK2kC/z1O2w0AMQj6d5pefND9JzultDFEqmD51rLqcksDWRKwILLhzWayNfFmg1GCztI+dANUzW501dj0gMqgd97nj5DPX4rTU56bBuT17tx8oYw/0ydTX8wAAAA=
`

type testCase struct {
	n int
	k int
}

// solve mirrors 1644F.go (stub output).
func solve(tc testCase) string {
	return "0"
}

func decodeTestcases() ([]testCase, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTestcases)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(r); err != nil {
		return nil, err
	}
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return nil, fmt.Errorf("no data")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t && pos+1 < len(fields); i++ {
		n, _ := strconv.Atoi(fields[pos])
		k, _ := strconv.Atoi(fields[pos+1])
		pos += 2
		cases = append(cases, testCase{n: n, k: k})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
	cmd := exec.Command(bin)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := solve(tc)
		got, err := runCandidate(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
