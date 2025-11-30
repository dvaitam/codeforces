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

// Embedded gzipped+base64 testcases from testcasesD.txt.
const encodedTestcases = `
H4sIAEsEK2kC/zVTWXbkIAz89yl0hPSSTnIcDNhsBsxmw+lHlt/89HNrqSqVxOMLHtD7Evca9JKnH/hA5NF5a6cnpqyffuEb2tbTHLidvuEFom87m/7gF/y87q70WLHoAatO0oQjUmMt0xueYBJLCPoAJYKoMiLANxzax4jpNyS5rdPjC/7glFz34e3YseQB5jQ5I+oLQorOpn6hvMGuh8qbIkR7xmNs2/RCmlnuF8wHkGUP7JS+YvwBcejpczFulduGmn/A1qPHNBeOAl6gTu5Jb1iId7ZjNtjyAi/V7jOhLLwgI7qgYpYqYOEbdGMj0DQLDy5T7NDBrFj5A1GoLVP0Baxtml3qHlDOaPpQZTgyp+jWkQAdrXV6YH6Qj9mhGcaeiPSEVY7k9suUJ+g4zIniUC2XXTWiwjQfGyMpIrh20tdWnWhkMSvJXuxvqLyGw1QvJJEZJHuCULgDbCNv8nk5J+I17gvmw60hKCJsLlU3U3gVgbe13hPNwyo58so7Arwhr8p2OQojvTLpsZC7kiHKE7Jli55pZL8LzDwhdFpQFmFXhXQF+o207ybiETinEsmlLvW+lsNuZTdhL/fy0L6L3Qflbdj8Scj+vlbZIgJcVV6QiDTnlCyVuNvCczbFOHFDr9ighEv4/zJSDZ86zfoFrSZW2qHZQsjWsUVeqR/g8qyMn/u63We6OHPa/218lksvc+MaHcfbKMHrte+3yJsl7s3MtK4ys5XM67JrQ96zeHokui7wUHtqeD1X6znTM8BXmX3eaDViY3YZBWk+wL1QGL9p6GyeYPvC90pmdIeRD1SphNT3Ncre79T1nH+h555Dy/V+WNXKwGOoZpDg2e3cTP8A/vqMQj4EAAA=
`

type testCase struct {
	n int
	k int
	s string
}

func solve(tc testCase) string {
	freq := make([]int, 26)
	for i := 0; i < tc.n; i++ {
		freq[tc.s[i]-'a']++
	}
	pairs := 0
	singles := 0
	for _, v := range freq {
		pairs += v / 2
		singles += v % 2
	}
	m := pairs / tc.k
	ans := 2 * m
	if singles+2*(pairs%tc.k) >= tc.k {
		ans++
	}
	return fmt.Sprint(ans)
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
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return testCase{}, err
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return testCase{}, err
	}
	s := fields[2]
	if len(s) != n {
		return testCase{}, fmt.Errorf("length mismatch: n=%d len(s)=%d", n, len(s))
	}
	return testCase{n: n, k: k, s: s}, nil
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		input := fmt.Sprintf("1 %d %d %s\n", tc.n, tc.k, tc.s)
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
