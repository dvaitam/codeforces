package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	n, k int
	a, b string
}

func parseCasesC(path string) ([]testCaseC, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, err
	}
	cases := make([]testCaseC, t)
	for i := 0; i < t; i++ {
		var n, k int
		if _, err := fmt.Fscan(in, &n, &k); err != nil {
			return nil, err
		}
		var s1, s2 string
		if _, err := fmt.Fscan(in, &s1); err != nil {
			return nil, err
		}
		if _, err := fmt.Fscan(in, &s2); err != nil {
			return nil, err
		}
		cases[i] = testCaseC{n: n, k: k, a: s1, b: s2}
	}
	return cases, nil
}

func solveC(tc testCaseC) int64 {
	letterIndex := make(map[byte]int)
	letters := make([]byte, 0, 10)
	for i := 0; i < tc.n; i++ {
		c := tc.a[i]
		if _, ok := letterIndex[c]; !ok {
			letterIndex[c] = len(letters)
			letters = append(letters, c)
		}
	}
	m := len(letters)
	if tc.k >= m {
		return int64(tc.n) * int64(tc.n+1) / 2
	}
	ans := int64(0)
	maxMask := 1 << m
	for mask := 0; mask < maxMask; mask++ {
		if bits.OnesCount(uint(mask)) > tc.k {
			continue
		}
		cur := int64(0)
		curLen := int64(0)
		for i := 0; i < tc.n; i++ {
			if tc.a[i] == tc.b[i] {
				curLen++
			} else {
				idx := letterIndex[tc.a[i]]
				if (mask>>idx)&1 == 1 {
					curLen++
				} else {
					cur += curLen * (curLen + 1) / 2
					curLen = 0
				}
			}
		}
		cur += curLen * (curLen + 1) / 2
		if cur > ans {
			ans = cur
		}
	}
	return ans
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCasesC("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "1\n%d %d\n%s\n%s\n", tc.n, tc.k, tc.a, tc.b)
		input := sb.String()
		expected := solveC(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(gotStr), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
