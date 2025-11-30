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
H4sIAMoDK2kC/zWO2Q3AMAhD/z0NRw7Yf7HatFWiBNCzTRgKjo0Lbxa+WG/4gRcCzrGTYLtV8S004QUPUT1/jdhkYtMlUpoU13quTMd5SVNED5n+fVNIjFnOLbHjuzU3hA377rkUNUsraX2Zk8PxoeJy0eJEspQg31QSOq70B1pDDBP9AAAA
`

const (
	mod = int64(1000000007)
	m   = int64(998244353)
)

func nextNum(curr, n int64) int64 {
	if curr*10 <= n {
		return curr * 10
	}
	for curr%10 == 9 || curr+1 > n {
		curr /= 10
	}
	return curr + 1
}

func solve(n int64) string {
	curr := int64(1)
	var ans int64
	for i := int64(1); i <= n; i++ {
		diff := (i - curr) % m
		if diff < 0 {
			diff += m
		}
		ans += diff
		ans %= mod
		if i == n {
			break
		}
		curr = nextNum(curr, n)
	}
	return fmt.Sprint(ans % mod)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines, err := decodeTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, line := range lines {
		val, err := strconv.ParseInt(strings.TrimSpace(line), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expect := solve(val)
		got, err := runCandidate(bin, line+"\n")
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
