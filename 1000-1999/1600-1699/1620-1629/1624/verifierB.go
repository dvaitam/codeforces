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

// Embedded gzipped+base64 testcases from testcasesB.txt.
const encodedTestcases = `
H4sIAAYEK2kC/y1Tt4HEMAzrMQVHMJNI7b/YA7qvLMsMSA7bY+PIa9vmBV/LMm942X2f4rNsy8GsxdgNVOmwiw0Vn8ZJy7SBh/GuG3nsmKPS7jWf/8pw3LYejb3OOptAcJv5h73GVfOhrdamMW5+rQ98OKw+1HsevOWWAWInTN74x04uTuwn2PEhnessEp1WenIz4ajkGM9sJy2W3E9ci5tw2EAUB32tWivz002B4gSRhxiR2wVp8VmNTTtHX5boWJqQfGSLoCLvIGFtBq0Gu4Uid2q7YCdr+8OU5meCs/1xtqN5INFLEVvwCZzYqEILAd8IZj9Nk3cHLKSPtKM+1ZJqUZMQw6SiLiPEovgVJEzE9zyY1yaVARbTGtcYGh5vk5xYW0oUOJ/c6xHZWBKXZEQWIYXc1UVOPG/ICgL0wB2jGzvC9dOXvjI3yDD2z1PTOfpC2HhTwnkf8v7slPRlMdu2FBQunddErT/W9JtLj78nG/VUHGixUwaWrZBTkrPg/YyVC5WSeCXbaYWWCvW+kB7h6gUzK/Dx8kvR70uaVjBGLlbJsIzA87/IxzNbIaEpjDIV4Q8VL0ckshd/F/lgqXEDAAA=
`

type testCase struct {
	a int64
	b int64
	c int64
}

func canMakeAP(a, b, c int64) bool {
	target := 2*b - c
	if target > 0 && target%a == 0 {
		return true
	}
	sum := a + c
	if sum%2 == 0 {
		mid := sum / 2
		if mid > 0 && mid%b == 0 {
			return true
		}
	}
	target = 2*b - a
	if target > 0 && target%c == 0 {
		return true
	}
	return false
}

func solve(tc testCase) string {
	if canMakeAP(tc.a, tc.b, tc.c) {
		return "YES"
	}
	return "NO"
}

func decodeTestcases() ([]string, error) {
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
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	var res []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			res = append(res, l)
		}
	}
	return res, nil
}

func parseCase(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return testCase{}, fmt.Errorf("invalid case: %q", line)
	}
	toInt := func(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }
	a, err := toInt(fields[0])
	if err != nil {
		return testCase{}, err
	}
	b, err := toInt(fields[1])
	if err != nil {
		return testCase{}, err
	}
	c, err := toInt(fields[2])
	if err != nil {
		return testCase{}, err
	}
	return testCase{a: a, b: b, c: c}, nil
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, line := range lines {
		tc, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := solve(tc)
		input := fmt.Sprintf("1 %d %d %d\n", tc.a, tc.b, tc.c)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
