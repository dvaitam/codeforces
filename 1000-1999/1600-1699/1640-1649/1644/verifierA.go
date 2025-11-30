package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded gzipped+base64 testcases from testcasesA.txt.
const encodedTestcases = `
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
`

func solveCase(s string) string {
	keys := make(map[rune]bool)
	for _, ch := range s {
		if ch >= 'a' && ch <= 'z' {
			keys[ch] = true
		} else {
			if !keys[ch+32] {
				return "NO"
			}
		}
	}
	return "YES"
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
	if len(lines) == 0 {
		return nil, fmt.Errorf("no data")
	}
	var cases []string
	for i, line := range lines {
		if i == 0 {
			continue // skip t
		}
		line = strings.TrimSpace(line)
		if line != "" {
			cases = append(cases, line)
		}
	}
	return cases, nil
}

func runCandidate(bin string, s string) (string, error) {
	input := fmt.Sprintf("1\n%s\n", s)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, s := range tests {
		expect := solveCase(s)
		got, err := runCandidate(bin, s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
