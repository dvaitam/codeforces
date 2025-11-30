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

type point struct{ x, y int }

type testCase struct {
	n int
	k int
	p []point
}

// Embedded gzipped+base64 testcases from testcasesD.txt.
const encodedTestcases = `
H4sIAGQLK2kC/01WWXYjMQj89yl0BO3L/S821EJPnhPHaQGCoig8Systfm5ps8TbKQufWv+NOKkl3mvZ5cVJPOnxGjBqD49+YSj7pT+Xf2gI1xFG4VJ/O+wieC8HVr3Cv20GYoABI3jF84PIjcdw7HBsjLzkHglFhEqj9VuZ5oTLilwbw4RVpIwEGy/A+YQPymk3TyYuVSZ088FA/EPXKo8XN3Xfgqoafhu8UM5giY8OHZ5Vj/H+dE9AimqA9LMly3/M7sIyHjXBwGRxzQEKk3iqKGDdkMpgpkd+um4z1KB96wqBS+Sy2NL48PiCQdbEcpBp9ucJPPS70UvVxoEBb0wHZT0nFCkIuxNPRZU4G2WroiHqiF2blmwHwm+2gRUwx8HkzAFksoVEgKQbdfzxJBKiHwi4EWey2C7gb7wmydpF5k6+bLe/EoR8SgjY1kHGslxkVosowoQP+Xg9ES2bXU2f68QYkugIhMczcvcxneWZ4OBt3qQ0RfSbROaMDpXkbrLxHq1NO9EcQac/qZebNSyEQ2NX5ORJ7eq/qJ7NYd289uq4e1qr2wzkNu4cOi/VmPQs7n/ozTo0SDdzQDi3Y7CxBOkp3WxT/DPhdj2E5oyAvZ4yKclggWTR1+PmPJaJ0asRxRvtyZ4UhQavQQJkp66gRhjyr5M6SPkm1mwfwd6WuO6I1fycIg5VcYrEi+PRTrZ3imIwXGLO8ACAdYTRpH2a9WWVXgIEaHBEe049c5mco2laexYpXJeQPFFFHbofzbuV1hPB6WK8lvqiqjS8ZuKRshCJgefHPdmFqwLF1uwjlY7CnaL2WPi2DscNvDBD5vP2PCdNmlTV9Wp1NT8npFDLyiGAW9TCB+uPZH+avj3XT5NIolTJaJNgU1EwEMX96i7mahuY41VdI1vIwWmZHanam0RPjdjaWBqMGEvhPDihEjZp5iWxubP6X34sZZXL7HmLc70M8tLkX9RAEEwr6FCfG/Ev0thRVMBITVvWlEBFeTUxtEv9hJEZKUi62MOqq6PPMnKg2ehpkakC35v2JQrd68DTwnWnLeljrZhtJdrFkKdUiS7uyNC0UTq2+nw1lsuEG6aDW/qkxzWFmer03VizzY96pQVD2h6tVVGsPXWZ30QkvVKY861w79cAlghcFm0JsAC8r91e7c+AjqJtCIhXNvxInioJM9XjV7wj2O0qgWhwXDnPWlrScm5RjYCk7efh4Ze7hjpzN2sloYb5+wewoS9gPwoAAA==
`

// Precomputed expected outputs for those testcases.
const encodedExpected = `
H4sIAM0LK2kC/z1QyxFEMQi6W4Ut+O2/tAXDvsnEQUE0CSsrD8REnMPrY1exRp5A631MI0PVS119rLCxk3yAb93H5bn/a8QL9QAPOng/PTH2SO3CPK7jcZqHE+fA2JypOXV+rK62C/nwdamdG5HqkXrUG8/1/YW0Q48feNXcricBAAA=
`

func decodeExpected() ([]string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedExpected)
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	return lines, nil
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
	var cases []testCase
	for pos+2 <= len(fields) {
		n, err := strconv.Atoi(fields[pos])
		if err != nil {
			return nil, err
		}
		k, err := strconv.Atoi(fields[pos+1])
		if err != nil {
			return nil, err
		}
		pos += 2
		if pos+2*n > len(fields) {
			break
		}
		p := make([]point, n+1)
		for j := 1; j <= n; j++ {
			x, _ := strconv.Atoi(fields[pos])
			y, _ := strconv.Atoi(fields[pos+1])
			p[j] = point{x: x, y: y}
			pos += 2
		}
		cases = append(cases, testCase{n: n, k: k, p: p})
	}
	return cases, nil
}

func runCandidate(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 1; i <= tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.p[i].x, tc.p[i].y)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	expected, err := decodeExpected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(expected) != len(cases) {
		fmt.Fprintf(os.Stderr, "expected %d outputs, got %d\n", len(cases), len(expected))
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := expected[idx]
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
